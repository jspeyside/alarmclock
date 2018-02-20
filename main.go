package main

import (
	"fmt"
	"github.com/jspeyside/alarmclock/domain"
	"github.com/jspeyside/alarmclock/rest"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"

	"net/http"
	"os"
)

const (
	port = 5050
)

var (
	app        = kingpin.New("alamrmclock", "")
	configFile = app.Flag("file", "Config file for hosts and mac addresses").Short('f').Required().String()
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

}

func main() {
	log.Info("Starting alarmclock server...")

	log.Debug("Parsing CLI")
	kingpin.Version(domain.Version)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	// Load Config
	loadConfig()

	// Start the server
	start()
}

func loadConfig() {
	raw, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Panic(err)
	}
	if err = yaml.Unmarshal(raw, &config); err != nil {
		log.Panic(err)
	}
}

func start() {
	router := httprouter.New()
	router.GET("/ping", rest.Ping)
	router.GET("/v1/:host", rest.Handler)
	router.GET("/", rest.Version)

	{
		port := fmt.Sprintf(":%d", port)
		log.Panic(http.ListenAndServe(port, router))
	}

}
