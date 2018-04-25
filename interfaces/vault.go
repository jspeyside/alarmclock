package interfaces

import (
	"github.com/hashicorp/vault/api"
	"github.com/jspeyside/alarmclock/domain"
)

type VaultClient struct {
	client *api.Client
}

func NewVaultClient(config *api.Config) (domain.VaultClient, error) {
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &VaultClient{
		client: client,
	}, nil
}

func (t *VaultClient) Read(path string) (*api.Secret, error) {
	return t.client.Logical().Read(path)
}

func (t *VaultClient) Write(path string, data map[string]interface{}) (*api.Secret, error) {
	return t.client.Logical().Write(path, data)
}
