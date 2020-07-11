package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/greatfocus/gf-frame/config"
	"github.com/greatfocus/gf-frame/database"
	"github.com/greatfocus/gf-user/models"
	"github.com/greatfocus/gf-user/repositories"
)

// NotifyRequest struct
type NotifyRequest struct {
	Host     string `json:"host,omitempty"`
	Port     string `json:"port,omitempty"`
	Messages []models.Notify
	Status   []bool
	URI      []map[string]interface{} `json:"uri,omitempty"`
}

// Task struct
type Task struct {
	notifyRepository *repositories.NotifyRepository
	config           *config.Config
	db               *database.DB
}

// Init required parameters
func (t *Task) Init(db *database.DB, config *config.Config) {
	t.notifyRepository = &repositories.NotifyRepository{}
	t.notifyRepository.Init(db)
	t.config = config
	t.db = db
}

// SendNotification intiates the job to send queued messages
func (t *Task) SendNotification() {
	log.Println("Scheduler_SendQueuedEmails started")
	request := NotifyRequest{
		Host: t.config.Notify.Host,
		Port: t.config.Notify.Port,
		URI:  t.config.Notify.URI,
	}
	t.SendNotifications(t.notifyRepository, &request)
	log.Println("Scheduler_SendQueuedEmails ended")
}

// SendNotifications function to send messages
func (t *Task) SendNotifications(repo *repositories.NotifyRepository, request *NotifyRequest) {
	log.Println("Scheduler_SendNotifications Fetching Notification messages")
	msgs, err := repo.GetNotification()
	if err != nil {
		log.Println("Scheduler_SendNotifications Error fetching Notification messages")
		return
	}

	if len(msgs) > 0 {
		request.Messages = msgs
		request.Status = make([]bool, len(msgs))
		sendBulkNotification(repo, msgs, request)
	} else {
		log.Println("Scheduler_SendNotifications Notification queued is empty")
	}
}

// SendBulk initiates sending of the messages
func sendBulkNotification(repo *repositories.NotifyRepository, msgs []models.Notify, request *NotifyRequest) {
	log.Println("Scheduler_SendNotifications Sending bulk Email messages")
	var wg sync.WaitGroup

	for i := 0; i < len(request.Messages); i++ {
		uri := GetURI(i, request.Messages[i], request)
		if uri == "" {
			// is a sad place to be
		} else {
			wg.Add(1)
			go SendNotify(i, request, uri, &wg)
		}
	}

	wg.Wait()
	updateNotifications(repo, msgs, request)
}

// updateMessage change message status
func updateNotifications(repo *repositories.NotifyRepository, msgs []models.Notify, request *NotifyRequest) {
	for i := 0; i < len(request.Messages); i++ {
		// check status of email sent
		err := repo.Update(msgs[i])
		if err != nil {
			log.Println("Failed to update Notification with ID", msgs[i].ID)
			log.Println(err)
		}
	}
}

// GetURI returns the operations uri
func GetURI(i int, msg models.Notify, request *NotifyRequest) string {
	var uri = ""
	for _, result := range request.URI {
		v := fmt.Sprintf("%v", result["operation"])
		templateID := fmt.Sprintf("%d", result["templateId"])
		if msg.Operation == v {
			if n, err := strconv.Atoi(templateID); err == nil {
				request.Messages[i].TemplateID = int64(n)
				uri = v
			}
		}
	}
	return uri
}

// SendNotify creates the messages
func SendNotify(i int, request *NotifyRequest, uri string, wg *sync.WaitGroup) {
	reqBody, err := json.Marshal(request.Messages[i])
	if err != nil {
		print(err)
	}
	resp, err := http.Post(request.Host+":"+uri,
		"application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		print(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
		request.Messages[i].Sent = false
		request.Messages[i].Status = "queue"
	}
	created := models.Notify{}
	err = json.Unmarshal(body, &created)

	if created.ID > 0 {
		request.Messages[i].Sent = true
		request.Messages[i].Status = "done"
	} else {
		request.Messages[i].Sent = false
		request.Messages[i].Status = "queue"
	}
	wg.Done()
}
