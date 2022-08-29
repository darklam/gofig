package gofig

import (
	"errors"
	"github.com/darklam/gofig/mocks"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
)

func TestGofig_PopulateConfigNoProvider(t *testing.T) {
	gofig := NewGofig()

	type child struct {
		URL string `env:"URL"`
	}

	type parent struct {
		C                   *child `provider:"vault" mountPath:"somewhere" secretPath:"secure"`
		SuperSecretPassword string `env:"SUPER_SECRET_PASSWORD" default:"1234"`
	}

	err := os.Setenv("SUPER_SECRET_PASSWORD", "4321")
	if err != nil {
		t.Error("got nil when setting environment")
	}

	err = os.Setenv("URL", "some_url")
	if err != nil {
		t.Error("got nil when setting environment")
	}

	cfg := new(parent)

	err = gofig.PopulateConfig(cfg)
	assert.Nil(t, err)

	assert.Equal(t, "4321", cfg.SuperSecretPassword)
	assert.Equal(t, "some_url", cfg.C.URL)
}

func TestGofig_PopulateConfigNoProviderVeryDeep(t *testing.T) {
	gofig := NewGofig()

	type childChildChild struct {
		EvenMoar string `env:"EVEN_MOAR"`
	}

	type childChild struct {
		Moar string `env:"MOAR"`
		Em   *childChildChild
	}

	type child struct {
		URL string `env:"URL"`
		M   *childChild
	}

	type parent struct {
		C                   *child
		SuperSecretPassword string `env:"SUPER_SECRET_PASSWORD" default:"1234"`
	}

	err := os.Setenv("SUPER_SECRET_PASSWORD", "4321")
	if err != nil {
		t.Error("got nil when setting environment")
	}

	err = os.Setenv("URL", "some_url")
	if err != nil {
		t.Error("got nil when setting environment")
	}

	err = os.Setenv("MOAR", "moar")
	if err != nil {
		t.Error("got nil when setting environment")
	}

	err = os.Setenv("EVEN_MOAR", "even_moar")
	if err != nil {
		t.Error("got nil when setting environment")
	}

	cfg := new(parent)

	err = gofig.PopulateConfig(cfg)
	assert.Nil(t, err)

	assert.Equal(t, "4321", cfg.SuperSecretPassword)
	assert.Equal(t, "some_url", cfg.C.URL)
	assert.Equal(t, "moar", cfg.C.M.Moar)
	assert.Equal(t, "even_moar", cfg.C.M.Em.EvenMoar)
}

func TestGofig_PopulateConfigWithProviderSimple(t *testing.T) {
	provider := &mocks.ProviderMock{
		NameFunc: func() string {
			return "test"
		},
		GetValueFunc: func(field reflect.StructField, parentField *reflect.StructField) (string, error) {
			return "something_else", nil
		},
	}

	gofig := NewGofig()

	gofig.RegisterProvider(provider)

	err := os.Setenv("SOMETHING", "something")
	if err != nil {
		t.Error("got nil when setting environment")
	}

	type simple struct {
		SomeField string `env:"SOMETHING" provider:"test"`
	}

	cfg := new(simple)

	err = gofig.PopulateConfig(cfg)

	calls := provider.GetValueCalls()
	call := calls[0]

	assert.Nil(t, err)
	assert.Equal(t, "something_else", cfg.SomeField)
	assert.Equal(t, 1, len(calls), "GetValue should have been called only once")
	assert.Equal(t, "SomeField", call.Field.Name, "GetValue should have been called on SomeField")
	assert.Nil(t, call.ParentField, "The parent should have been nil")
}

func TestGofig_PopulateConfigWithProviderDeep(t *testing.T) {
	provider := &mocks.ProviderMock{
		NameFunc: func() string {
			return "test"
		},
		GetValueFunc: func(field reflect.StructField, parentField *reflect.StructField) (string, error) {
			if parentField == nil {
				return "", nil
			}
			value := field.Name + "_value"
			if parentField != nil {
				value = parentField.Name + "_" + value
			}
			return value, nil
		},
	}

	gofig := NewGofig()

	gofig.RegisterProvider(provider)

	err := os.Setenv("VALUE", "different")
	assert.Nil(t, err)

	type childChild struct {
		Value string
	}

	type child struct {
		Value string
		CC    *childChild `provider:"test"`
	}

	type parent struct {
		Value string `provider:"test" env:"VALUE"`
		C     *child `provider:"test"`
	}

	cfg := new(parent)

	err = gofig.PopulateConfig(cfg)
	assert.Nil(t, err)

	assert.Equal(t, "different", cfg.Value)
	assert.Equal(t, "C_Value_value", cfg.C.Value)
	assert.Equal(t, "CC_Value_value", cfg.C.CC.Value)
}

func TestGofig_PopulateConfigProviderError(t *testing.T) {
	provider := &mocks.ProviderMock{
		NameFunc: func() string {
			return "test"
		},
		GetValueFunc: func(field reflect.StructField, parentField *reflect.StructField) (string, error) {
			return "", errors.New("some error")
		},
	}

	gofig := NewGofig()

	gofig.RegisterProvider(provider)

	type config struct {
		Value string `provider:"test"`
	}

	cfg := new(config)

	err := gofig.PopulateConfig(cfg)

	assert.Equal(t, err.Error(), "some error")
}

func TestGofig_PopulateConfigDefaultValues(t *testing.T) {
	gofig := NewGofig()

	err := os.Setenv("SOMETHING", "value")
	assert.Nil(t, err)

	type config struct {
		Value   string `env:"SOMETHING"`
		Another string `default:"one"`
	}

	cfg := new(config)

	err = gofig.PopulateConfig(cfg)
	assert.Nil(t, err)

	assert.Equal(t, "value", cfg.Value)
	assert.Equal(t, "one", cfg.Another)
}

func TestGofig_PopulateConfigErrorOnInvalidValue(t *testing.T) {
	gofig := NewGofig()

	err := os.Setenv("SOMETHING", "10")
	assert.Nil(t, err)

	type config struct {
		Value int `env:"SOMETHING"`
	}

	cfg := new(config)

	err = gofig.PopulateConfig(cfg)

	assert.Equal(t, err.Error(), "only struct pointers and strings are allowed")
}

func TestGofig_PopulateConfigErrorOnInvalidValuePointer(t *testing.T) {
	gofig := NewGofig()

	err := os.Setenv("SOMETHING", "10")
	assert.Nil(t, err)

	type config struct {
		Value *string `env:"SOMETHING"`
	}

	cfg := new(config)

	err = gofig.PopulateConfig(cfg)

	assert.Equal(t, err.Error(), "only struct pointers and strings are allowed")
}

func TestGofig_PopulateConfigInvalidTags(t *testing.T) {
	gofig := NewGofig()

	err := os.Setenv("SOMETHING", "10")
	assert.Nil(t, err)

	type config struct {
		Value string
	}

	cfg := new(config)

	err = gofig.PopulateConfig(cfg)

	assert.Equal(t, err.Error(), "no provider, env or default for field Value")
}
