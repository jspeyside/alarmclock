package rest

import (
	"fmt"
	"github.com/ghthor/gowol"
	"github.com/jspeyside/alarmclock/domain"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	// wol "github.com/sabhiram/go-wol"
	"github.com/masterzen/winrm"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	// "os"
)

var (
	Config domain.Config
)

// Wake handles the main api endpoint for alarm clock and wakes hosts up
func Wake(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	_, err := ioutil.ReadAll(req.Body)
	if err != nil {
		writeError(w, errors.Wrap(err, "Error reading request body"))
		return
	}

	hostname := params.ByName("host")
	host := Config.Hosts[hostname]
	mac := host.MacAddress
	if mac == "" {
		log.Errorf("No mac address found in config for %s", hostname)
		return
	}
	log.Debugf("Waking %s[%s]", hostname, mac)

	// Wake up the host
	err = wakeOnLAN(mac)
	if err != nil {
		writeError(w, errors.Wrap(err, "Error waking host"))
		return
	}
}

func Sleep(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	hostname := params.ByName("host")
	host := Config.Hosts[hostname]
	mac := host.MacAddress
	if mac == "" {
		log.Errorf("No mac address found in config for %s", hostname)
		return
	}
	log.Debugf("Putting %s to sleep", hostname)
	endpoint := winrm.NewEndpoint(hostname, 5985, false, false, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, host.Username, host.Password)
	if err != nil {
		panic(err)
	}
	code, err := client.Run("shutdown /s /t 0", w, w)
	if err != nil {
		writeError(w, errors.Wrap(err, "Error putting host to slep"))
		return
	}
	log.Debug("Sleep cmd exited with ", code)
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

func wakeOnLAN(mac string) error {
	return wol.MagicWake(mac, Config.Broadcast)
	// return wol.SendMagicPacket(mac, Config.Broadcast+":9", "")
}
