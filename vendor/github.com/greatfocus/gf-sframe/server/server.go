package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/greatfocus/gf-sframe/database"
	"github.com/greatfocus/gf-sframe/logger"
	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

// NewServer get new instance of server
func NewServer(serviceName, URI string) *Server {
	// Load environment variables
	env := os.Getenv("ENV")
	if env == "" || os.Getenv("ENV") == "dev" {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal(err)
		}
	}
	env = os.Getenv("ENV")

	// init creates instance of logger
	logger := logger.NewLogger(serviceName)

	timeout, err := strconv.ParseUint(os.Getenv("SERVER_TIMEOUT"), 0, 64)
	if err != nil {
		log.Fatal(fmt.Println(err))
	}

	srv := Server{
		Name:     serviceName,
		URI:      URI,
		Env:      env,
		Logger:   logger,
		Cache:    initCache(),
		JWT:      initJWT(),
		Database: initDatabase(logger),
		Timeout:  timeout,
	}
	return &srv
}

// Server struct
type Server struct {
	Name             string
	Env              string
	URI              string
	Mux              *http.ServeMux
	Cache            *cache.Cache
	Database         database.Database
	JWT              JWT
	Logger           *logrus.Logger
	clientPublicKey  *rsa.PublicKey
	ServerPublicKey  *rsa.PublicKey
	serverPrivateKey *rsa.PrivateKey
	Timeout          uint64
}

// Start the server
func (s *Server) Start() {
	pki(s)

	setUploadPath(s.Mux, s.URI)

	serverProbe(s.Mux, s.URI)

	start(s.Mux, s.Logger, int(s.Timeout))
}

func pki(s *Server) {
	privatekey, publicKey := GetServerPKI()
	s.ServerPublicKey = publicKey
	s.serverPrivateKey = privatekey
	clientPublicKey := os.Getenv("CLIENT_PUBLICKEY")
	if clientPublicKey != "" {
		clientPublicKeyString, err := base64.StdEncoding.DecodeString(clientPublicKey)
		if err == nil {
			publicBlock, _ := pem.Decode([]byte(clientPublicKeyString))
			pubKey, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
			if err == nil {
				s.clientPublicKey = pubKey.(*rsa.PublicKey)
			}

		}
	}
}

func serverProbe(mux *http.ServeMux, uri string) {
	probe := liveProbe{}
	probeLoc := "/" + uri + "/info"
	mux.Handle(probeLoc, probe)
}

func initCache() *cache.Cache {
	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	expireVal, err := strconv.ParseInt(os.Getenv("CACHE_EXPIRE"), 0, 64)
	if err != nil {
		log.Fatal(fmt.Println(err))
	}

	intervalVal, err := strconv.ParseInt(os.Getenv("CACHE_INTERVAL"), 0, 64)
	if err != nil {
		log.Fatal(fmt.Println(err))
	}

	return cache.New(time.Duration(expireVal), time.Duration(intervalVal))
}

func initJWT() JWT {
	minutes, err := strconv.ParseInt(os.Getenv("JWT_Minutes"), 0, 64)
	if err != nil {
		log.Fatal(fmt.Println(err))
	}

	authorized, err := strconv.ParseBool(os.Getenv("JWT_Authorized"))
	if err != nil {
		log.Fatal(fmt.Println(err))
	}

	secret := os.Getenv("JWT_Secret")
	return NewJWT(secret, minutes, authorized)
}

func initDatabase(logger *logrus.Logger) database.Database {
	logger.Info(fmt.Sprintln("Preparing Database configuration"))
	host := os.Getenv("DB_HOST")
	databaseName := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	sslmode := "require"
	sslkey := os.Getenv("DB_SSL_KEY")
	sslcert := os.Getenv("DB_SSL_CERT")
	var sslkeyPath, sslcertPath, sslcaPath string

	// prepare ssl connection files
	if sslkey != "" && sslcert != "" {
		sslkeyPath, sslcertPath, sslcaPath = getDatabaseCertificate()
	}

	port, err := strconv.ParseUint(os.Getenv("DB_PORT"), 0, 64)
	if err != nil {
		log.Fatal(fmt.Println(err))
	}

	maxLifetimeVal, err := strconv.ParseUint(os.Getenv("DB_MaxLifetime"), 0, 64)
	if err != nil {
		log.Fatal(fmt.Println(err))
	}
	maxLifetime := time.Duration(maxLifetimeVal) * time.Minute

	maxIdleConns, err := strconv.ParseInt(os.Getenv("DB_MaxIdleConns"), 0, 64)
	if err != nil {
		log.Fatal(fmt.Println(err))
	}

	maxOpenConns, err := strconv.ParseInt(os.Getenv("DB_MaxOpenConns"), 0, 64)
	if err != nil {
		log.Fatal(fmt.Println(err))
	}

	var psqlInfo string
	if sslcaPath != "" {
		sslmode = "verify-full"
		psqlInfo = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s sslrootcert=%s sslkey=%s sslcert=%s",
			host, port, user, password, databaseName, sslmode, sslcaPath, sslkeyPath, sslcertPath)
	} else {
		psqlInfo = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s sslkey=%s sslcert=%s",
			host, port, user, password, databaseName, sslmode, sslkeyPath, sslcertPath)
	}

	// create database connection
	params := database.DatabaseParam{
		ConnectionStr: psqlInfo,
		DatabaseName:  databaseName,
		MaxLifetime:   maxLifetime,
		MaxOpenConns:  int(maxOpenConns),
		MaxIdleConns:  int(maxIdleConns),
	}

	return database.NewConnection(params, logger)
}

// setUploadPath creates an upload path
func setUploadPath(mux *http.ServeMux, uri string) {
	path := os.Getenv("UPLOAD_PATH")
	if path == "" {
		path = "./data/upload"
	}
	fs := http.FileServer(http.Dir(path))
	fileLoc := "/" + uri + "/resource"
	mux.Handle(fileLoc+"/", http.StripPrefix(fileLoc, fs))
}

// start creates server instance
func start(mux *http.ServeMux, logger *logrus.Logger, timeout int) {
	addr := ":" + os.Getenv("SERVER_PORT")
	srv := &http.Server{
		Addr:           addr,
		ReadTimeout:    time.Duration(timeout) * time.Second,
		WriteTimeout:   time.Duration(timeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        mux,
	}

	// Get key certificate
	logger.Info("Listening to port HTTP" + addr)
	crt, key := GetServerCertificate()
	if crt != "" && key != "" {
		log.Fatal(srv.ListenAndServeTLS(crt, key))
	}
	log.Fatal(srv.ListenAndServe())
}

// Success returns object as json
func (s *Server) Success(w http.ResponseWriter, r *http.Request, data interface{}) {
	if data != nil {
		s.response(w, r, data, "success")
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	s.response(w, r, nil, "success")
}

// Error returns error as json
func (s *Server) Error(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		s.response(w, r, struct {
			Error string `json:"error"`
		}{Error: err.Error()}, "error")
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	s.response(w, r, nil, "error")
}

// request returns payload
func (s *Server) Request(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		req := Params{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			derr := errors.New("invalid payload request")
			w.WriteHeader(http.StatusBadRequest)
			s.Error(w, r, derr)
			return nil, err
		}

		err = s.checkRequestId(req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			s.Error(w, r, err)
			return nil, err
		}

		if s.clientPublicKey != nil {
			res, err := serverDecrypt(req.Cipher, s.serverPrivateKey)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				s.Error(w, r, err)
				return nil, err
			}
			return res.Params, nil
		}
		return req.Params, nil
	}
	return nil, nil
}

// CheckRequestId validates requestID
func (s *Server) checkRequestId(p Params) error {
	if p.ID == "" {
		return errors.New("invalid request")
	}

	_, found := s.Cache.Get(p.ID)
	if found {
		return errors.New("duplicate request")
	}

	s.Cache.Set(p.ID, p.Params, time.Duration(s.Timeout)*time.Second)
	return nil
}

// response returns payload
func (s *Server) response(w http.ResponseWriter, r *http.Request, data interface{}, message string) {
	out, _ := json.Marshal(data)
	if s.clientPublicKey != nil {
		result, err := serverEncrypt(string(out), s.clientPublicKey)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			s.Error(w, r, err)
		}
		res := Response{
			Cipher: result,
		}
		_ = json.NewEncoder(w).Encode(res)
	} else {
		res := Response{
			Result: string(out),
		}
		_ = json.NewEncoder(w).Encode(res)
	}

}

// Connect method make a database connection
func getDatabaseCertificate() (cert, key, ca string) {
	var root string
	key = CreateSSLCert("postgresql-client.key", os.Getenv("DB_SSL_KEY"))
	cert = CreateSSLCert("postgresql-client.crt", os.Getenv("DB_SSL_CERT"))
	sslrootcert := os.Getenv("DB_ROOT_CA")
	if sslrootcert != "" {
		root = CreateSSLCert("postgresql-ca.crt", os.Getenv("DB_ROOT_CA"))
	}
	return key, cert, root
}

// CreateSSLCert makes cert in image
func CreateSSLCert(filename string, content string) string {
	var path = os.Getenv("APP_PATH") + "/ssl/" + filename
	path = filepath.Clean(path)
	_, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			_ = file.Close()
		}()
		err = os.Chmod(path, 0600)
		if err != nil {
			log.Fatal(err)
		}

		cnt, err := base64.StdEncoding.DecodeString(content)
		if err != nil {
			file.WriteString(content)
		}
		file.WriteString(string(cnt))
	}

	return path
}

// Response result
type Response struct {
	Result string `json:"result,omitempty"`
	Cipher string `json:"cipher,omitempty"`
}

// Params
type Params struct {
	ID     string      `json:"id,omitempty"`
	Params interface{} `json:"params,omitempty"`
	Cipher string      `json:"cipher,omitempty"`
}

// decrypt payload
func serverDecrypt(cipherText string, privateKey *rsa.PrivateKey) (Params, error) {
	params := Params{}
	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	unencrypted, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ct)
	if err != nil {
		derr := errors.New("invalid payload request")
		return params, derr
	}
	err = json.Unmarshal([]byte(unencrypted), &params)
	if err != nil {
		derr := errors.New("invalid payload request")
		return params, derr
	}
	return params, nil
}

// encrypt payload
func serverEncrypt(payload string, publicKey *rsa.PublicKey) (string, error) {
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(payload))
	if err != nil {
		derr := errors.New("invalid payload request")
		return "", derr
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// GetServerCertificate returns private and public key
func GetServerCertificate() (string, string) {
	var sslcert = os.Getenv("API_SSL_CERT")
	var sslkey = os.Getenv("API_SSL_KEY")

	// prepare ssl connection files
	if sslkey != "" && sslcert != "" {
		crt := CreateSSLCert("api-server.crt", sslcert)
		key := CreateSSLCert("api-server.key", sslkey)
		return crt, key
	}

	return "", ""
}

// GetServerPKI returns public key infrustructure
func GetServerPKI() (*rsa.PrivateKey, *rsa.PublicKey) {
	var privateKey = os.Getenv("API_PRIVATE_KEY")
	privateKeyString, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, nil
	}
	privateBlock, _ := pem.Decode([]byte(privateKeyString))
	privKey, err := x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
	if err != nil {
		return nil, nil
	}

	var publicKey = os.Getenv("API_PUBLIC_KEY")
	publicKeyString, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, nil
	}
	publicBlock, _ := pem.Decode([]byte(publicKeyString))
	pubKey, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
	if err != nil {
		return nil, nil
	}

	return privKey, pubKey.(*rsa.PublicKey)
}
