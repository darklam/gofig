# Gofig - a simple (but powerful) config library
![Coverage](https://img.shields.io/badge/Coverage-3-red)

Gofig is a lightweight and extendable configuration library for Go projects. Originally created for personal use,
it has been shared to benefit the community. ❤️

## Usage

```shell
go get -u github.com/darklam/gofig@v1.0.0
```

This example demonstrates the main features of Gofig:

```go
package main

import (
	"github.com/darklam/gofig"
	"github.com/darklam/gofig/providers"
)

type PgConfig struct {
	Username string `prop:"username"`
	Password string `prop:"password"`
	Host     string `prop:"host"`
	Port     string `prop:"port"`
}

type Config struct {
	Port      string    `prop:"port" default:"3000"`
	Postgres  *PgConfig `prop:"postgres"`
	RedisHost string    `prop:"redis.host" default:"localhost"`
}

func main() {
	fig := gofig.NewGofig()
	
	json5Provider, err := providers.NewJSONProvider("local.json5")
	if err != nil {
		panic(err)
    }
	
	fig.RegisterProvider(json5Provider)
	
	envProvider := providers.NewEnvProvider()
	
	fig.RegisterProvider(envProvider)
	
	cfg := new(Config)
	
	err = fig.PopulateConfig(cfg)

	if err != nil {
		panic(err)
	}
	
    // The cfg variable now has all the fields populated
}
```

## General information

The v1 version of the library has been improved to support JSON (and JSON5) sources and offer better extensibility.

However, it has also become more opinionated about naming properties in the configuration. This decision helps keep the code simple without sacrificing extendability.


## Available tags

* **prop**: Specifies the name of the property which will be used to fetch its value from the different providers
* **default**: The default value of the field

**_You must use at least one of these tags in the struct fields (unless if the field is a struct pointer)_**

## Field types

Gofig supports two field types: string and a pointer to a struct.

If a field is a string, it will be treated as a field to populate.

If a field is a struct pointer, it will be replaced by an instance of the struct with its fields populated according to the tags of its fields.

Fields with no value found will have empty strings. Structs are always instantiated. All fields will be populated recursively.



## Provider precedence

The default value has the lowest precedence if set.

The order in which providers are registered determines precedence. For example, if you register providers in this order:

- env
- json


Values fetched from JSON will replace the values from the environment, which in turn will replace the defaults where set.

## Adding more providers

Built-in providers include:

- env
- json

To create custom providers, implement the interfaces/Provider interface in your code (see the interface documentation for more information).

If you think a new provider might be useful, please create a PR.

## Testing

You can run:
```shell
go generate ./...
go test $(go list ./... | grep -v tools)
```

Alternatively you can just run the run_tests.sh script.

The tools directory is excluded, as it only contains imports to manage the versions of tools and
it would just error out.

## Contributing

It's not a popular package to have specific guidelines. If you have issues or want to add extra functionality
(hopefully more providers) feel free to open an issue or even better create a PR.

Just be civil and respectful.