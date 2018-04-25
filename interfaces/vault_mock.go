package interfaces

import (
	"errors"
	"github.com/hashicorp/vault/api"
	"github.com/jspeyside/alarmclock/domain"
)

type MockVaultClient struct {
	data   map[string](map[string]interface{})
	errors []string
}

func NewMockVaultClient(errors []string) domain.VaultClient {
	return &MockVaultClient{
		errors: errors,
		data:   make(map[string](map[string]interface{})),
	}
}

func (t *MockVaultClient) Read(path string) (*api.Secret, error) {
	for _, clientErr := range t.errors {
		if clientErr == "bad_client" {
			return nil, errors.New(clientErr)
		}
	}
	if t.data[path] == nil {
		return nil, nil
	}
	return &api.Secret{
		Data: t.data[path],
	}, nil
}

func (t *MockVaultClient) Write(path string, data map[string]interface{}) (*api.Secret, error) {
	t.data[path] = data
	return &api.Secret{Data: data}, nil
}
