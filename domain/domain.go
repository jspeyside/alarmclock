package domain

import (
	"github.com/hashicorp/vault/api"
)

const (
	VaultPath = "secret/alarmclock"
)

var (
	// Vault is the client connection to the secure vault server
	Vault VaultClient
	// Version is the current version of the app. It is generated from VERSION at build time
	Version string
)

type VaultClient interface {
	Read(path string) (*api.Secret, error)
	Write(path string, data map[string]interface{}) (*api.Secret, error)
}
