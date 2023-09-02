package providers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/darklam/gofig/mocks/providers"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
)

func TestValidateOptions(t *testing.T) {
	testCases := []struct {
		name    string
		options VaultOptions
		wantErr error
	}{
		{
			name: "Valid AppRoleAuth Options",
			options: VaultOptions{
				AppRoleAuth: &VaultAppRoleAuthOptions{
					RoleId:   "test-role-id",
					SecretId: "test-secret-id",
				},
			},
			wantErr: nil,
		},
		{
			name: "Valid KubernetesAuth Options",
			options: VaultOptions{
				KubernetesAuth: &VaultKubernetesAuthOptions{
					Jwt:  "test-jwt",
					Role: "test-role",
				},
			},
			wantErr: nil,
		},
		{
			name: "Both Auth Options Present (Invalid)",
			options: VaultOptions{
				AppRoleAuth: &VaultAppRoleAuthOptions{
					RoleId:   "test-role-id",
					SecretId: "test-secret-id",
				},
				KubernetesAuth: &VaultKubernetesAuthOptions{
					Jwt:  "test-jwt",
					Role: "test-role",
				},
			},
			wantErr: ErrInvalidVaultAuthConfig,
		},
		{
			name:    "No Auth Options (Invalid)",
			options: VaultOptions{},
			wantErr: ErrInvalidVaultAuthConfig,
		},
		{
			name: "AppRoleAuth Set to nil and KubernetesAuth with Valid Options",
			options: VaultOptions{
				AppRoleAuth: nil,
				KubernetesAuth: &VaultKubernetesAuthOptions{
					Jwt:  "test-jwt",
					Role: "test-role",
				},
			},
			wantErr: nil,
		},
		{
			name: "KubernetesAuth Set to nil and AppRoleAuth with Valid Options",
			options: VaultOptions{
				KubernetesAuth: nil,
				AppRoleAuth: &VaultAppRoleAuthOptions{
					RoleId:   "test-role-id",
					SecretId: "test-secret-id",
				},
			},
			wantErr: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := validateOptions(testCase.options)
			if !errors.Is(err, testCase.wantErr) {
				t.Errorf("Expected error: %v, but got: %v", testCase.wantErr, err)
			}
		})
	}
}

func TestSetupVaultClient(t *testing.T) {
	options := VaultOptions{
		KubernetesAuth: &VaultKubernetesAuthOptions{
			Jwt:  "test-jwt",
			Role: "test-role",
		},
		MountPath:      "test-mount-path",
		Path:           "test-path",
		RequestTimeout: 15,
		Url:            "test-url",
	}

	optionsAppRoleAuth := VaultOptions{
		AppRoleAuth: &VaultAppRoleAuthOptions{
			RoleId:   "test-role-id",
			SecretId: "test-secret-id",
		},
		MountPath:      "test-mount-path",
		Path:           "test-path",
		RequestTimeout: 15,
		Url:            "test-url",
	}

	t.Run("Initialize called with correct arguments for kubernetes auth", func(t *testing.T) {
		// GIVEN
		client := providers.NewMockVaultClienter(t)
		ctx := context.Background()

		client.
			EXPECT().
			Initialize(options.Url, time.Second*time.Duration(options.RequestTimeout)).
			Return(nil)

		client.
			EXPECT().
			KubernetesLogin(ctx, options.KubernetesAuth.Jwt, options.KubernetesAuth.Role).
			Return(nil)

		// WHEN
		err := setupVaultClient(ctx, client, options)

		// THEN
		assert.Nil(t, err)
	})

	t.Run("Initialize called with correct arguments for approle auth", func(t *testing.T) {
		// GIVEN
		client := providers.NewMockVaultClienter(t)
		ctx := context.Background()

		client.
			EXPECT().
			Initialize(options.Url, time.Second*time.Duration(options.RequestTimeout)).
			Return(nil)

		client.
			EXPECT().
			AppRoleLogin(ctx, optionsAppRoleAuth.AppRoleAuth.RoleId, optionsAppRoleAuth.AppRoleAuth.SecretId).
			Return(nil)

		// WHEN
		err := setupVaultClient(ctx, client, optionsAppRoleAuth)

		// THEN
		assert.Nil(t, err)
	})

	t.Run("Returns correct error for initialization", func(t *testing.T) {
		// GIVEN
		client := providers.NewMockVaultClienter(t)
		ctx := context.Background()

		client.
			EXPECT().
			Initialize(options.Url, time.Second*time.Duration(options.RequestTimeout)).
			Return(errors.New("something went wrong"))

		// WHEN
		err := setupVaultClient(ctx, client, options)

		// THEN
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrVaultConnection)
	})

	t.Run("Returns correct error for auth", func(t *testing.T) {
		// GIVEN
		client := providers.NewMockVaultClienter(t)
		ctx := context.Background()

		client.
			EXPECT().
			Initialize(options.Url, time.Second*time.Duration(options.RequestTimeout)).
			Return(nil)

		client.
			EXPECT().
			KubernetesLogin(ctx, options.KubernetesAuth.Jwt, options.KubernetesAuth.Role).
			Return(errors.New("something went wrong"))

		// WHEN
		err := setupVaultClient(ctx, client, options)

		// THEN
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrVaultAuth)
	})
}

func TestGetSecretData(t *testing.T) {
	options := VaultOptions{
		Path:      "test-path",
		MountPath: "test-mount-path",
	}

	t.Run("Returns correct error for GetValues", func(t *testing.T) {
		// GIVEN
		client := providers.NewMockVaultClienter(t)
		ctx := context.Background()

		client.
			EXPECT().
			GetValues(ctx, options.Path, options.MountPath).
			Return(nil, errors.New("something went wrong"))

		// WHEN
		_, err := getSecretData(ctx, client, options)

		// THEN
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrVaultSecretFetch)
	})

	t.Run("Mapping the values from Vault", func(t *testing.T) {
		// GIVEN
		client := providers.NewMockVaultClienter(t)
		ctx := context.Background()

		resultMap := map[string]interface{}{
			"test":      "value",
			"another":   "value",
			"otherTest": "data",
		}

		client.
			EXPECT().
			GetValues(ctx, options.Path, options.MountPath).
			Return(resultMap, nil)

		// WHEN
		result, err := getSecretData(ctx, client, options)

		// THEN
		assert.Nil(t, err)
		for _, key := range maps.Keys(resultMap) {
			assert.Equal(t, resultMap[key], result[key])
		}
	})

	t.Run("Returns correct error when the secret value is not a string", func(t *testing.T) {
		// GIVEN
		client := providers.NewMockVaultClienter(t)
		ctx := context.Background()

		resultMap := map[string]interface{}{
			"test":    "value",
			"invalid": 100,
		}

		client.
			EXPECT().
			GetValues(ctx, options.Path, options.MountPath).
			Return(resultMap, nil)

		// WHEN
		_, err := getSecretData(ctx, client, options)

		// THEN
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, ErrVaultSecretValueType)
	})
}

func TestVaultProvider_GetValue(t *testing.T) {
	provider := &VaultProvider{
		data: map[string]string{
			"SOME_KEY":                   "some_value",
			"ANOTHER_KEY_MULTIPLE_PARTS": "foo",
			"MIXED_CASING":               "something",
		},
	}

	testCases := []struct {
		name      string
		keys      []string
		wantValue string
	}{
		{
			name:      "Resolves simple key",
			keys:      []string{"some", "key"},
			wantValue: "some_value",
		},
		{
			name:      "Resolves key with multiple parts",
			keys:      []string{"another", "key", "multiple", "parts"},
			wantValue: "foo",
		},
		{
			name:      "Resolves key with mixed casing",
			keys:      []string{"mIxEd", "CASING"},
			wantValue: "something",
		},
		{
			name:      "Returns empty string if key is not found",
			keys:      []string{"not", "found"},
			wantValue: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value, err := provider.GetValue(testCase.keys)
			assert.Nil(t, err)
			assert.Equal(t, testCase.wantValue, value)
		})
	}
}
