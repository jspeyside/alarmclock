package rest

import (
	"fmt"
	"github.com/jspeyside/alarmclock/domain"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// Handler handles the main api endpoint for alarm clock and wakes hosts up
func Handler(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	_, err := ioutil.ReadAll(req.Body)
	host := params.ByName("host")
	if err != nil {
		writeError(w, errors.Wrap(err, "Error reading request body"))
		return
	}

	// Wake up the host
	err = wake(host)
	if err != nil {
		writeError(w, errors.Wrap(err, "Error reading request body"))
		return
	}
}

// Ping is a healthcheck endpoint and returns a simple pong
func Ping(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Write(toBytes("pong"))
}

// Version is a version endpoint for the api
func Version(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	v := fmt.Sprintf("Version: %s", domain.Version)
	w.Write(toBytes(v))
}

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	w.Write(toBytes(err.Error()))
	log.Error(err)
}

func toBytes(s string) []byte {
	return []byte(s)
}

func wake(host string) error {
	// Placeholder
	return nil
}
