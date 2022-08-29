package providers

import (
	"context"
	"errors"
	"fmt"
	"github.com/darklam/gofig/interfaces"
	vault "github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
	"reflect"
)

// VaultProvider is the build-in provider for Hashicorp's Vault
// It currently only supports approle auth and the KV engine
type VaultProvider struct {
	cache  map[string]map[string]interface{}
	client interfaces.Valuer
}

// Name returns this provider's name
func (vp VaultProvider) Name() string {
	return "vault"
}

// GetValue returns the value retrieved from Vault for this field
func (vp *VaultProvider) GetValue(field reflect.StructField, parentField *reflect.StructField) (string, error) {
	mountPath := ""
	secretPath := ""
	key := field.Name

	if parentField != nil {
		mountPath = parentField.Tag.Get("mountPath")
		secretPath = parentField.Tag.Get("secretPath")
	}

	fieldMountPath := field.Tag.Get("mountPath")
	fieldSecretPath := field.Tag.Get("secretPath")
	fieldKey := field.Tag.Get("key")

	if fieldMountPath != "" {
		mountPath = fieldMountPath
	}

	if fieldSecretPath != "" {
		secretPath = fieldSecretPath
	}

	if fieldKey != "" {
		key = fieldKey
	}

	cacheKey := fmt.Sprintf("%s/%s", mountPath, secretPath)

	if _, found := vp.cache[cacheKey]; !found {
		result, err := vp.client.Value(mountPath, secretPath)
		if err != nil {
			return "", err
		}
		vp.cache[cacheKey] = result
	}

	value, found := vp.cache[cacheKey][key]
	if !found {
		return "", nil
	}

	result, ok := value.(string)
	if !ok {
		return "", errors.New(fmt.Sprintf("value type assertion failed for %s/%s - %s", mountPath, secretPath, key))
	}

	return result, nil
}

// VaultWrapper wraps the Vault API to aid in unit testing
// It implements the interfaces.Valuer interface
type VaultWrapper struct {
	client *vault.Client
}

// Value returns the value present in mountPath/secretPath of the configured Vault instance
func (vault *VaultWrapper) Value(mountPath string, secretPath string) (map[string]interface{}, error) {
	result, err := vault.client.KVv2(mountPath).Get(context.Background(), secretPath)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}

func newVaultWrapperAppRole(address string, roleId string, secretId string) (*VaultWrapper, error) {
	vaultConfig := vault.DefaultConfig()
	vaultConfig.Address = address

	client, err := vault.NewClient(vaultConfig)
	if err != nil {
		return nil, err
	}

	roleSecretId := &approle.SecretID{
		FromString: secretId,
	}

	appRoleAuth, err := approle.NewAppRoleAuth(roleId, roleSecretId)
	if err != nil {
		return nil, err
	}

	if appRoleAuth == nil {
		return nil, errors.New("could not get appRoleAuth")
	}

	authInfo, err := client.Auth().Login(context.Background(), appRoleAuth)
	if err != nil {
		return nil, err
	}

	if authInfo == nil {
		return nil, errors.New("could not get authInfo")
	}

	return &VaultWrapper{
		client: client,
	}, nil
}

// NewVaultProviderAppRole returns a new Vault provider using approle authentication
// It's named this way to maybe allow more authentication types in the future
func NewVaultProviderAppRole(address string, roleId string, secretId string) (*VaultProvider, error) {
	wrapper, err := newVaultWrapperAppRole(address, roleId, secretId)
	if err != nil {
		return nil, err
	}

	return &VaultProvider{
		client: wrapper,
		cache:  map[string]map[string]interface{}{},
	}, nil
}
