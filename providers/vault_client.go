package providers

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

// VaultClienter Serves as an abstraction layer to the actual vault client
// We're using this, so we can unit test the vault provider without worrying about the Vault client
type VaultClienter interface {
	Initialize(url string, requestTimeout time.Duration) error
	AppRoleLogin(ctx context.Context, roleId string, secretId string) error
	KubernetesLogin(ctx context.Context, jwt string, role string) error
	GetValues(ctx context.Context, path string, mountPath string) (map[string]interface{}, error)
}

type VaultClient struct {
	client *vault.Client
}

func NewVaultClient() *VaultClient {
	return &VaultClient{}
}

func (vc *VaultClient) Initialize(url string, requestTimeout time.Duration) error {
	client, err := vault.New(vault.WithAddress(url), vault.WithRequestTimeout(requestTimeout))
	if err != nil {
		return err
	}

	vc.client = client
	return nil
}

func (vc *VaultClient) AppRoleLogin(ctx context.Context, roleId string, secretId string) error {
	res, err := vc.client.Auth.AppRoleLogin(ctx, schema.AppRoleLoginRequest{
		RoleId:   roleId,
		SecretId: secretId,
	})

	if err != nil {
		return err
	}

	return vc.client.SetToken(res.Auth.ClientToken)
}

func (vc *VaultClient) KubernetesLogin(ctx context.Context, jwt string, role string) error {
	res, err := vc.client.Auth.KubernetesLogin(ctx, schema.KubernetesLoginRequest{
		Jwt:  jwt,
		Role: role,
	})
	if err != nil {
		return nil
	}

	return vc.client.SetToken(res.Auth.ClientToken)
}

func (vc *VaultClient) GetValues(ctx context.Context, path string, mountPath string) (map[string]interface{}, error) {
	result, err := vc.client.Secrets.KvV2Read(ctx, path, vault.WithMountPath(mountPath))
	if err != nil {
		return nil, err
	}
	return result.Data.Data, nil
}
