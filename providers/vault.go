package providers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/exp/maps"
)

type VaultAppRoleAuthOptions struct {
	RoleId   string
	SecretId string
}

type VaultKubernetesAuthOptions struct {
	Jwt  string
	Role string
}

type VaultOptions struct {
	// The Vault Server url
	Url string

	// The request timeout for the vault client in seconds (default 1m)
	RequestTimeout int

	// Options for app role authentication
	AppRoleAuth *VaultAppRoleAuthOptions

	// Options for kubernetes authentication
	KubernetesAuth *VaultKubernetesAuthOptions

	// The KV mount path
	MountPath string

	// The path of the secret
	Path string
}

type VaultProvider struct {
	mountPath string
	data      map[string]string
}

func NewVaultProvider(options VaultOptions) (*VaultProvider, error) {
	vaultClient := NewVaultClient()
	err := validateOptions(options)
	if err != nil {
		return nil, fmt.Errorf("vault config invalid: %w", err)
	}

	ctx := context.Background()

	err = setupVaultClient(ctx, vaultClient, options)
	if err != nil {
		return nil, err
	}

	data, err := getSecretData(ctx, vaultClient, options)
	if err != nil {
		return nil, err
	}

	return &VaultProvider{
		data: data,
	}, nil
}

func (vp *VaultProvider) GetValue(fieldPath []string) (string, error) {
	key := strings.ToUpper(strings.Join(fieldPath, "_"))
	value, exists := vp.data[key]
	if !exists {
		return "", nil
	} else {
		return value, nil
	}
}

func validateOptions(options VaultOptions) error {
	if options.AppRoleAuth == nil && options.KubernetesAuth == nil {
		return ErrInvalidVaultAuthConfig
	}

	if options.AppRoleAuth != nil && options.KubernetesAuth != nil {
		return ErrInvalidVaultAuthConfig
	}

	return nil
}

func setupVaultClient(ctx context.Context, client VaultClienter, options VaultOptions) error {
	err := client.Initialize(options.Url, time.Duration(options.RequestTimeout)*time.Second)
	if err != nil {
		return errors.Join(ErrVaultConnection, err)
	}

	if options.AppRoleAuth != nil {
		err = client.AppRoleLogin(ctx, options.AppRoleAuth.RoleId, options.AppRoleAuth.SecretId)
	} else if options.KubernetesAuth != nil {
		err = client.KubernetesLogin(ctx, options.KubernetesAuth.Jwt, options.KubernetesAuth.Role)
	}

	if err != nil {
		return errors.Join(ErrVaultAuth, err)
	}

	return nil
}

func getSecretData(ctx context.Context, client VaultClienter, options VaultOptions) (map[string]string, error) {
	result, err := client.GetValues(ctx, options.Path, options.MountPath)
	if err != nil {
		return nil, errors.Join(ErrVaultSecretFetch, err)
	}

	keys := maps.Keys(result)
	data := make(map[string]string, len(keys))

	for _, key := range keys {
		value, ok := result[key].(string)
		if !ok {
			return nil, ErrVaultSecretValueType
		}

		data[key] = value
	}

	return data, nil
}
