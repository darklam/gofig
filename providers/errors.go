package providers

import "errors"

var (
	ErrInvalidVaultAuthConfig = errors.New("exactly one auth method options must be specified")
	ErrVaultConnection        = errors.New("error connecting to the Vault server")
	ErrVaultAuth              = errors.New("error authenticating with Vault")
	ErrVaultSecretFetch       = errors.New("error fetching secret from Vault")
	ErrVaultSecretValueType   = errors.New("error getting secret value as string")
)
