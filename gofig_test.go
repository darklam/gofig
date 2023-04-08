package gofig

import (
	"github.com/darklam/gofig/interfaces/interfacesfakes"
	"github.com/darklam/gofig/providers"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGofig_PopulateConfigNoProvider(t *testing.T) {
	// GIVEN
	gofig := NewGofig()

	envProvider := providers.NewEnvProvider()

	gofig.RegisterProvider(envProvider)

	type child struct {
		URL string `prop:"url"`
	}

	type parent struct {
		C                   *child `prop:"smth"`
		SuperSecretPassword string `prop:"super.secret.password" default:"1234"`
	}

	err := os.Setenv("SUPER_SECRET_PASSWORD", "4321")
	if err != nil {
		t.Error("got error when setting environment")
	}

	err = os.Setenv("SMTH_URL", "some_url")
	if err != nil {
		t.Error("got error when setting environment")
	}

	cfg := new(parent)

	// WHEN
	err = gofig.PopulateConfig(cfg)

	// THEN
	assert.Nil(t, err)
	assert.Equal(t, "4321", cfg.SuperSecretPassword)
	assert.Equal(t, "some_url", cfg.C.URL)
}

func TestGofig_PopulateConfigNoProviderVeryDeep(t *testing.T) {
	// GIVEN
	gofig := NewGofig()

	envProvider := providers.NewEnvProvider()

	gofig.RegisterProvider(envProvider)

	type childChildChild struct {
		EvenMoar string `prop:"even.moar"`
	}

	type childChild struct {
		Moar string           `prop:"moar"`
		Em   *childChildChild `prop:"em"`
	}

	type child struct {
		URL string      `prop:"url"`
		M   *childChild `prop:"m"`
	}

	type parent struct {
		C                   *child `prop:"c"`
		SuperSecretPassword string `prop:"super.secret.password" default:"1234"`
	}

	err := os.Setenv("SUPER_SECRET_PASSWORD", "4321")
	if err != nil {
		t.Error("got nil when setting environment")
	}

	err = os.Setenv("C_URL", "some_url")
	if err != nil {
		t.Error("got nil when setting environment")
	}

	err = os.Setenv("C_M_MOAR", "moar")
	if err != nil {
		t.Error("got nil when setting environment")
	}

	err = os.Setenv("C_M_EM_EVEN_MOAR", "even_moar")
	if err != nil {
		t.Error("got nil when setting environment")
	}

	cfg := new(parent)

	// WHEN
	err = gofig.PopulateConfig(cfg)

	// THEN
	assert.Nil(t, err)
	assert.Equal(t, "4321", cfg.SuperSecretPassword)
	assert.Equal(t, "some_url", cfg.C.URL)
	assert.Equal(t, "moar", cfg.C.M.Moar)
	assert.Equal(t, "even_moar", cfg.C.M.Em.EvenMoar)
}

func TestGofig_PopulateConfigWithProviderSimple(t *testing.T) {
	// GIVEN
	provider := &interfacesfakes.FakeProvider{}

	provider.GetValueReturns("providerValue", nil)

	gofig := NewGofig()

	gofig.RegisterProvider(provider)

	err := os.Setenv("SOMETHING", "something")
	if err != nil {
		t.Error("got nil when setting environment")
	}

	type simple struct {
		SomeField string `prop:"smth"`
	}

	cfg := new(simple)

	// WHEN
	err = gofig.PopulateConfig(cfg)

	// THEN
	providerGetValueArgs := provider.GetValueArgsForCall(0)

	assert.Nil(t, err)
	assert.Equal(t, 1, provider.GetValueCallCount())
	assert.Equal(t, 1, len(providerGetValueArgs))
	assert.Equal(t, "smth", providerGetValueArgs[0])
	assert.Equal(t, "providerValue", cfg.SomeField)
}

func TestGofig_PopulateConfigProviderPrecedence(t *testing.T) {
	// GIVEN
	gofig := NewGofig()

	provider1 := &interfacesfakes.FakeProvider{}
	provider2 := &interfacesfakes.FakeProvider{}
	provider3 := &interfacesfakes.FakeProvider{}

	provider1.GetValueReturns("provider1", nil)

	provider2.GetValueStub = func(strings []string) (string, error) {
		property := strings[0]
		if property == "provider2" || property == "provider3" {
			return "provider2", nil
		}
		return "", nil
	}

	provider3.GetValueStub = func(strings []string) (string, error) {
		if strings[0] == "provider3" {
			return "provider3", nil
		}
		return "", nil
	}

	gofig.RegisterProvider(provider1)
	gofig.RegisterProvider(provider2)
	gofig.RegisterProvider(provider3)

	type config struct {
		Value1 string `prop:"provider1"`
		Value2 string `prop:"provider2"`
		Value3 string `prop:"provider3"`
	}

	cfg := new(config)

	// WHEN
	err := gofig.PopulateConfig(cfg)

	// THEN
	assert.Nil(t, err)
	assert.Equal(t, cfg.Value1, "provider1")
	assert.Equal(t, cfg.Value2, "provider2")
	assert.Equal(t, cfg.Value3, "provider3")
}
