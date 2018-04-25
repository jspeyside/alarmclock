package rest

import (
	"fmt"
	"github.com/jspeyside/alarmclock/domain"
	"github.com/jspeyside/alarmclock/interfaces"
	"github.com/julienschmidt/httprouter"
	check "gopkg.in/check.v1"
	"net/http/httptest"
	"testing"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type RestSuite struct {
	// mockNetwork  Network
	// mockStatsite *mockStatsite
	config domain.Config
}

var _ = check.Suite(&RestSuite{})

func (s *RestSuite) SetUpSuite(c *check.C) {
	host := domain.Host{
		MacAddress: "12:34:56:78:90:AB",
		Username:   "Username",
		Password:   "Password",
	}
	invalidMac := domain.Host{
		MacAddress: "12:34:56:78:90:AB:00:00:00",
	}
	emptyHost := domain.Host{}
	s.config = domain.Config{
		Broadcast: "127.255.255.255",
		Hosts: map[string]domain.Host{
			"host1":      host,
			"invalidmac": invalidMac,
			"emptyHost":  emptyHost,
		},
	}
	Config = s.config
}

func (s *RestSuite) TestWake(c *check.C) {
	request := httptest.NewRequest("GET", "/v1/wake/", nil)
	recorder := httptest.NewRecorder()
	param := httprouter.Param{
		Key:   "host",
		Value: "host1",
	}
	Wake(recorder, request, []httprouter.Param{param})
	resp := recorder.Result()
	c.Assert(resp.StatusCode, check.Equals, 200)
	c.Assert(recorder.Body.String(), check.Equals, "")
}

func (s *RestSuite) TestWakeUndefiendHost(c *check.C) {
	request := httptest.NewRequest("GET", "/v1/wake/", nil)
	recorder := httptest.NewRecorder()
	param := httprouter.Param{
		Key:   "host",
		Value: "undefined",
	}
	Wake(recorder, request, []httprouter.Param{param})
	resp := recorder.Result()
	c.Assert(resp.StatusCode, check.Equals, 500)
	c.Assert(recorder.Body.String(), check.Equals, "No mac address found in config for undefined")
}

func (s *RestSuite) TestWakeNoHost(c *check.C) {
	request := httptest.NewRequest("GET", "/v1/wake/", nil)
	recorder := httptest.NewRecorder()
	Wake(recorder, request, nil)
	resp := recorder.Result()
	c.Assert(resp.StatusCode, check.Equals, 500)
	c.Assert(recorder.Body.String(), check.Equals, "No mac address found in config for ")
}

func (s *RestSuite) TestWakeInvalidMac(c *check.C) {
	request := httptest.NewRequest("GET", "/v1/wake/", nil)
	recorder := httptest.NewRecorder()
	param := httprouter.Param{
		Key:   "host",
		Value: "invalidmac",
	}
	Wake(recorder, request, []httprouter.Param{param})
	resp := recorder.Result()
	c.Assert(resp.StatusCode, check.Equals, 500)
	errorString := fmt.Sprintf("Error waking host: Invalid MAC Address String: %s", s.config.Hosts["invalidmac"].MacAddress)
	c.Assert(recorder.Body.String(), check.Matches, errorString)
}

func (s *RestSuite) TestSleep(c *check.C) {
	request := httptest.NewRequest("GET", "/v1/sleep/", nil)
	recorder := httptest.NewRecorder()
	param := httprouter.Param{
		Key:   "host",
		Value: "host1",
	}
	Sleep(recorder, request, []httprouter.Param{param})
	resp := recorder.Result()
	c.Assert(resp.StatusCode, check.Equals, 500)
	c.Assert(recorder.Body.String(), check.Matches, "Error putting host to sleep.*")
}

func (s *RestSuite) TestSleepNoCredentials(c *check.C) {
	request := httptest.NewRequest("GET", "/v1/sleep/", nil)
	recorder := httptest.NewRecorder()
	param := httprouter.Param{
		Key:   "host",
		Value: "emptyHost",
	}
	Sleep(recorder, request, []httprouter.Param{param})
	resp := recorder.Result()
	c.Assert(resp.StatusCode, check.Equals, 500)
	c.Assert(recorder.Body.String(), check.Equals, "Missing username/password for emptyHost")
}

func (s *RestSuite) TestSleepVault(c *check.C) {
	hostname := "host1"
	domain.Vault = interfaces.NewMockVaultClient([]string{})
	data := make(map[string]interface{})
	userKey := fmt.Sprintf("%s_username", hostname)
	passKey := fmt.Sprintf("%s_password", hostname)
	data[userKey] = "username"
	data[passKey] = "password"
	domain.Vault.Write(domain.VaultPath, data)
	request := httptest.NewRequest("GET", "/v1/sleep/", nil)
	recorder := httptest.NewRecorder()
	param := httprouter.Param{
		Key:   "host",
		Value: hostname,
	}
	Sleep(recorder, request, []httprouter.Param{param})
	resp := recorder.Result()
	c.Assert(resp.StatusCode, check.Equals, 500)
	c.Assert(recorder.Body.String(), check.Matches, "Error putting host to sleep.*")
}

func (s *RestSuite) TestSleepVaultNoSecret(c *check.C) {
	domain.Vault = interfaces.NewMockVaultClient([]string{})
	request := httptest.NewRequest("GET", "/v1/sleep/", nil)
	recorder := httptest.NewRecorder()
	param := httprouter.Param{
		Key:   "host",
		Value: "host1",
	}
	Sleep(recorder, request, []httprouter.Param{param})
	resp := recorder.Result()
	c.Assert(resp.StatusCode, check.Equals, 500)
	c.Assert(recorder.Body.String(), check.Matches, "Error putting host to sleep.*")
}

func (s *RestSuite) TestSleepVaultNoData(c *check.C) {
	domain.Vault = interfaces.NewMockVaultClient([]string{})
	data := make(map[string]interface{})
	domain.Vault.Write(domain.VaultPath, data)
	request := httptest.NewRequest("GET", "/v1/sleep/", nil)
	recorder := httptest.NewRecorder()
	param := httprouter.Param{
		Key:   "host",
		Value: "host1",
	}
	Sleep(recorder, request, []httprouter.Param{param})
	resp := recorder.Result()
	c.Assert(resp.StatusCode, check.Equals, 500)
	c.Assert(recorder.Body.String(), check.Matches, "Error putting host to sleep.*")
}

func (s *RestSuite) TestSleepVaultClientErr(c *check.C) {
	clientErrs := []string{"bad_client"}
	domain.Vault = interfaces.NewMockVaultClient(clientErrs)
	request := httptest.NewRequest("GET", "/v1/sleep/", nil)
	recorder := httptest.NewRecorder()
	param := httprouter.Param{
		Key:   "host",
		Value: "host1",
	}
	Sleep(recorder, request, []httprouter.Param{param})
	resp := recorder.Result()
	c.Assert(resp.StatusCode, check.Equals, 500)
	c.Assert(recorder.Body.String(), check.Matches, "Error putting host to sleep.*")
}

func (s *RestSuite) TestPing(c *check.C) {
	request := httptest.NewRequest("GET", "/ping", nil)
	recorder := httptest.NewRecorder()

	Ping(recorder, request, nil)
	resp := recorder.Result()
	c.Assert(resp.StatusCode, check.Equals, 200)
	c.Assert(recorder.Body.String(), check.Equals, "pong")
}

func (s *RestSuite) TestVersion(c *check.C) {
	request := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	Version(recorder, request, nil)
	resp := recorder.Result()
	c.Assert(resp.StatusCode, check.Equals, 200)
	versionStr := fmt.Sprintf("Version: %s", domain.Version)
	c.Assert(recorder.Body.String(), check.Equals, versionStr)
}
