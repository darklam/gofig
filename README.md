# Gofig - a simple (but powerful) config library

This was created mainly to use for my own projects, but felt like a good thing to share ❤️

## Usage

```shell
go get -u github.com/darklam/gofig@latest
```

This example should explain pretty much everything

```go
package main

import (
	"github.com/darklam/gofig"
	"github.com/darklam/gofig/providers"
	"os"
)

type PgConfig struct {
	Username string `env:"PG_USERNAME" key:"USERNAME"`
	Password string `env:"PG_PASSWORD" key:"PASSWORD"`
	Host     string `env:"PG_HOST" key:"HOST"`
	Port     string `env:"PG_PORT" key:"PORT"`
}

type Config struct {
	Port      string    `env:"PORT" default:"3000"`
	Postgres  *PgConfig `provider:"vault" mountPath:"kv" secretPath:"postgres"`
	RedisHost string    `env:"REDIS_HOST" default:"localhost" provider:"vault" mountPath:"kv" secretPath:"redis" key:"HOST"`
}

func main() {
	fig := gofig.NewGofig()
	
	vaultAddress := os.Getenv("VAULT_ADDR")
	roleId := os.Getenv("ROLE_ID")
	secretId := os.Getenv("SECRET_ID")
	vault, err := providers.NewVaultProviderAppRole(vaultAddress, roleId, secretId)
	if err != nil {
		panic(err)
        }
	
	fig.RegisterProvider(vault)
	
	cfg := new(Config)
	
	err = fig.PopulateConfig(cfg)

	if err != nil {
		panic(err)
        }
	
    // The cfg variable now has all the fields populated
}
```

## Available tags

* **provider**: The name of the provider responsible for populating this field
* **env**: The name of the environment variable to populate the field with
* **default**: The default value of the field

**_You must use at least one of these tags in the struct fields (unless if the field is a struct pointer)_**

All other tags are provider specific.

## Field types

The only supported types are string and a pointer to a struct. If a string then it will be treated
as a field to populate. If a struct pointer it will be replaced by an instance of the struct with its 
fields populated according to the tags of its fields.

Fields for which the value was not found by anything will be empty strings.

Structs are always instantiated.

All fields will be populated recursively. GO nuts!

## Tag precedence

* default
* env
* provider

This means that if an env variable is found it will replace the default value and if a value from the
provider is retrieved it will replace the env or default values.

## Adding more providers

Currently, there's only the Vault provider with AppRole auth which is what I needed, but I might add more
providers in the future.

You can also create custom ones in your code by implementing the interfaces/Provider
interface (more documentation there). Please do create a PR to add a new provider if you think it might be useful.

## Testing

This module uses moq for mocks, so the first thing you need is to install moq
```shell
go install github.com/matryer/moq@v0.2.7
```
Currently using v0.2.7 so keep it the same in case you dug this up 5 years later and there are breaking changes

Then just run 
```shell
go test -v ./...
```

## Contributing

It's not a popular package to have specific guidelines. If you have issues or want to add extra functionality
(hopefully more providers) feel free to open an issue or even better create a PR.

Just be civil and respectful.