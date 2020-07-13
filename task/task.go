package task

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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
		wg.Add(1)
		go UpdateNotify(repo, request.Messages[i])
		go SendNotify(i, request, &wg)
	}

	wg.Wait()
	updateNotifications(repo, msgs, request)
}

// UpdateNotify change message status
func UpdateNotify(repo *repositories.NotifyRepository, msgs models.Notify) {
	// check status of email sent
	msgs.Status = "processing"
	err := repo.Update(msgs)
	if err != nil {
		log.Println("Failed to update Notification with ID", msgs.ID)
		log.Println(err)
	}
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

// SendNotify creates the messages
func SendNotify(i int, request *NotifyRequest, wg *sync.WaitGroup) {
	reqBody, err := json.Marshal(request.Messages[i])
	if err != nil {
		print(err)
	}
	resp, err := http.Post(request.Host+":"+request.Port+request.Messages[i].URI,
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
