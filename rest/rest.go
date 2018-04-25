package rest

import (
	"fmt"
	"github.com/ghthor/gowol"
	"github.com/jspeyside/alarmclock/domain"
	"github.com/julienschmidt/httprouter"
	"github.com/masterzen/winrm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	// "os"
)

var (
	// Config defines a configuration for waking/sleeping hosts
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
		writeError(w, errors.Errorf("No mac address found in config for %s", hostname))
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

// Sleep shuts down a configured PC
func Sleep(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	hostname := params.ByName("host")
	host := Config.Hosts[hostname]
	// Attempt to load username/password from vault
	user, pass, err := loadVaultCredentials(hostname)
	if err == nil {
		host.Username = user
		host.Password = pass
	} else {
		log.Error(err)
	}
	if host.Username == "" || host.Password == "" {
		writeError(w, errors.Errorf("Missing username/password for %s", hostname))
		return
	}
	log.Debugf("Putting %s to sleep", hostname)
	endpoint := winrm.NewEndpoint(hostname, 5985, false, false, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, host.Username, host.Password)
	if err != nil {
		writeError(w, err)
	}
	code, err := client.Run("shutdown /s /t 0", w, w)
	if err != nil {
		writeError(w, errors.Wrap(err, "Error putting host to sleep"))
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
}

func loadVaultCredentials(hostname string) (username string, password string, err error) {
	if domain.Vault == nil {
		return "", "", errors.New("No Vault Client")
	}
	usernameKey := fmt.Sprintf("%s_username", hostname)
	passwordKey := fmt.Sprintf("%s_password", hostname)
	secret, err := domain.Vault.Read(domain.VaultPath)
	if err != nil {
		return "", "", errors.New("Error loading vault client")
	}
	if secret == nil {
		return "", "", fmt.Errorf("No secret at %s", domain.VaultPath)
	}
	data := secret.Data
	if data == nil || data[usernameKey] == nil || data[passwordKey] == nil {
		return "", "", errors.New("No data loaded for vault secret")
	}
	username = fmt.Sprintf("%s", data[usernameKey])
	password = fmt.Sprintf("%s", data[passwordKey])
	return username, password, nil
}
