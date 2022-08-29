package providers

import (
	"errors"
	"github.com/darklam/gofig/mocks"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestVaultProvider_Name(t *testing.T) {
	client := mocks.ValuerMock{
		ValueFunc: func(mountPath string, secretPath string) (map[string]interface{}, error) {
			m := map[string]interface{}{
				"some_value": mountPath + "_" + secretPath + "_",
			}
			return m, nil
		},
	}
	vault := &VaultProvider{client: &client, cache: map[string]map[string]interface{}{}}

	name := vault.Name()

	assert.Equal(t, "vault", name)
}

func TestVaultProvider_GetValueNoParent(t *testing.T) {
	client := mocks.ValuerMock{
		ValueFunc: func(mountPath string, secretPath string) (map[string]interface{}, error) {
			m := map[string]interface{}{
				"oof": mountPath + "_" + secretPath + "_" + "oof",
			}
			return m, nil
		},
	}
	vault := &VaultProvider{client: &client, cache: map[string]map[string]interface{}{}}

	type config struct {
		Field1 string `mountPath:"something" secretPath:"secret" key:"oof"`
	}

	cfg := new(config)

	field := reflect.TypeOf(cfg).Elem().Field(0)

	result, err := vault.GetValue(field, nil)
	assert.Nil(t, err)
	assert.Equal(t, "something_secret_oof", result)
}

func TestVaultProvider_GetValueWithParent(t *testing.T) {
	client := mocks.ValuerMock{
		ValueFunc: func(mountPath string, secretPath string) (map[string]interface{}, error) {
			m := map[string]interface{}{
				"oof": mountPath + "_" + secretPath + "_" + "oof",
			}
			return m, nil
		},
	}
	vault := &VaultProvider{client: &client, cache: map[string]map[string]interface{}{}}

	type config struct {
		Child  string `key:"oof"`
		Parent string `mountPath:"something" secretPath:"secret"`
	}

	cfg := new(config)

	typ := reflect.TypeOf(cfg).Elem()

	childField := typ.Field(0)
	parentField := typ.Field(1)

	result, err := vault.GetValue(childField, &parentField)
	assert.Nil(t, err)
	assert.Equal(t, "something_secret_oof", result)
}

func TestVaultProvider_GetValueVaultError(t *testing.T) {
	client := mocks.ValuerMock{
		ValueFunc: func(mountPath string, secretPath string) (map[string]interface{}, error) {
			return nil, errors.New("i have a snake in my boot")
		},
	}
	vault := &VaultProvider{client: &client, cache: map[string]map[string]interface{}{}}

	type config struct {
		Child  string `key:"oof"`
		Parent string `mountPath:"something" secretPath:"secret"`
	}

	cfg := new(config)

	typ := reflect.TypeOf(cfg).Elem()

	childField := typ.Field(0)
	parentField := typ.Field(1)

	result, err := vault.GetValue(childField, &parentField)
	assert.NotNil(t, err)
	assert.Equal(t, "", result)
	assert.Equal(t, "i have a snake in my boot", err.Error())
}

func TestVaultProvider_GetValueKeyNotFound(t *testing.T) {
	client := mocks.ValuerMock{
		ValueFunc: func(mountPath string, secretPath string) (map[string]interface{}, error) {
			m := map[string]interface{}{
				"other_oof": mountPath + "_" + secretPath + "_" + "oof",
			}
			return m, nil
		},
	}
	vault := &VaultProvider{client: &client, cache: map[string]map[string]interface{}{}}

	type config struct {
		Child  string `key:"oof"`
		Parent string `mountPath:"something" secretPath:"secret"`
	}

	cfg := new(config)

	typ := reflect.TypeOf(cfg).Elem()

	childField := typ.Field(0)
	parentField := typ.Field(1)

	result, err := vault.GetValue(childField, &parentField)
	assert.Nil(t, err)
	assert.Equal(t, "", result)
}

func TestVaultProvider_GetValueTypeAssertionError(t *testing.T) {
	client := mocks.ValuerMock{
		ValueFunc: func(mountPath string, secretPath string) (map[string]interface{}, error) {
			m := map[string]interface{}{
				"oof": 42,
			}
			return m, nil
		},
	}
	vault := &VaultProvider{client: &client, cache: map[string]map[string]interface{}{}}

	type config struct {
		Child  string `key:"oof"`
		Parent string `mountPath:"something" secretPath:"secret"`
	}

	cfg := new(config)

	typ := reflect.TypeOf(cfg).Elem()

	childField := typ.Field(0)
	parentField := typ.Field(1)

	result, err := vault.GetValue(childField, &parentField)
	assert.NotNil(t, err)
	assert.Equal(t, "", result)
	assert.Equal(t, "value type assertion failed for something/secret - oof", err.Error())
}
