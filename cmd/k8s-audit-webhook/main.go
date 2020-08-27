package main

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"github.com/gorilla/mux"
	"os"
	"strconv"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	EnvFilename = "LOG_FILENAME"
	EnvMaxSize = "LOG_MAX_SIZE"
	EnvMaxBackups = "LOG_MAX_BACKUPS"
	EnvMaxAge = "LOG_MAX_AGE"
	EnvCompress = "LOG_COMPRESS"
	EnvListenAddr = "LOG_LISTEN_ADDR"
)

var (
	logger *log.Logger
)

func getEnvDefaultString(key, defaultValue string) string {
	if v, exists := os.LookupEnv(key); exists != true {
		log.Printf("looking for %s, defaulting to %s", key, defaultValue)
		return defaultValue
	} else {
		log.Printf("looking for %s, found %s", key, v)
		return v
	}
}

func getEnvDefaultInt(key string, defaultValue int) int {
	if v, exists := os.LookupEnv(key); exists != true {
		log.Printf("looking for %s, defaulting to %d", key, defaultValue)
		return defaultValue
	} else {
		if vv, err := strconv.Atoi(v); err != nil {
			log.Printf("looking for %s, found %s, defaulting to %d, error: %s", key, v, defaultValue, err)
			return defaultValue
		} else {
			log.Printf("looking for %s, found %d", key, vv)
			return vv
		}
	}
}

func getEnvDefaultBool(key string, defaultValue bool) bool {
	if v, exists := os.LookupEnv(key); exists != true {
		log.Printf("looking for %s, defaulting to %t", key, defaultValue)
		return defaultValue
	} else {
		if vv, err := strconv.ParseBool(v); err != nil {
			log.Printf("looking for %s, found %s, defaulting to %t, error: %s", key, v, defaultValue, err)
			return defaultValue
		} else {
			log.Printf("looking for %s, found %t", key, vv)
			return vv
		}
	}
}

func defaultHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("unhandled: %s", req.URL)
}

func auditHandler(w http.ResponseWriter, req *http.Request) {
	if body, err := ioutil.ReadAll(req.Body); err != nil {
		log.Printf("error reading body: %s", err)
	} else {
		logger.Print(string(body))
	}
}

func main() {
	logger = log.New(&lumberjack.Logger{
		Filename: getEnvDefaultString(EnvFilename, "k8s-audit-webhook.log"),
		MaxSize:  getEnvDefaultInt(EnvMaxSize, 1024),
		MaxBackups: getEnvDefaultInt(EnvMaxBackups, 30),
		MaxAge: getEnvDefaultInt(EnvMaxAge, 30),
		Compress: getEnvDefaultBool(EnvCompress, true),
	}, "", 0)

	r := mux.NewRouter()
	r.HandleFunc("/audit", auditHandler)
	r.PathPrefix("/").HandlerFunc(defaultHandler)

	log.Fatal(http.ListenAndServe(getEnvDefaultString(EnvListenAddr, ":8080"), r))
}
