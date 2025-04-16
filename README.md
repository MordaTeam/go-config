# 🛠️ go-config: Universal Configuration Library for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/hikitani/go-config.svg)](https://pkg.go.dev/github.com/hikitani/go-config)
[![Go Report Card](https://goreportcard.com/badge/github.com/hikitani/go-config)](https://goreportcard.com/report/github.com/hikitani/go-config)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`go-config` is a universal library for Go that simplifies working with configuration files and data from various sources. The library provides a flexible interface for loading, decoding, and merging configurations.

## 🌟 Key Features

- **Support for multiple data sources**:

  - Command-line arguments
  - Consul KV
  - Environment variables
  - Files
  - Any object implementing the `io.Reader` interface

- **Flexibility in decoding**:

  - Built-in JSON decoder
  - Ability to use third-party decoders (e.g., XML)

- **Configuration merging**:

  - Loading configuration from multiple sources
  - Merging strategies: `AllOf` (all sources must succeed) and `OneOf` (at least one source succeeds)

- **Additional features**:
  - Partial filling of existing structures
  - Simple integration into existing projects

---

## 📦 Installation

```bash
go get github.com/hikitani/go-config
```

---

## 🚀 Usage

### Simple example: Reading configuration from `io.Reader`

```go
type Config struct {
    Foo string `json:"foo"`
    Bar int    `json:"bar"`
}

source := strings.NewReader(`{"foo": "hello", "bar": 42}`)
cfg, err := config.New[Config](config.FromReader(source))
if err != nil {
    panic(err)
}

fmt.Printf("%+v\n", cfg) // Output: {Foo:hello Bar:42}
```

> By default, the JSON decoder is used.

---

### Using a different format (XML)

```go
type Config struct {
    Sizes []string `xml:"size"`
}

source := strings.NewReader(`
<sizes>
    <size>small</size>
    <size>regular</size>
    <size>large</size>
</sizes>`)

cfg, err := config.New[Config](
    config.FromReader(source),
    config.WithDecoder(xml.NewDecoder),
)
if err != nil {
    panic(err)
}

fmt.Printf("%+v\n", cfg) // Output: {Sizes:[small regular large]}
```

---

### Partially filling an existing configuration

If you already have a partially filled structure, you can fill it further:

```go
type Config struct {
    Foo string `json:"foo"`
    Bar string `json:"bar"`
}

source := strings.NewReader(`{"bar": "hello"}`)
cfg := Config{
    Foo: "initialized",
}

if err := config.Fill(&cfg, config.FromReader(source)); err != nil {
    panic(err)
}

fmt.Printf("%+v\n", cfg) // Output: {Foo:initialized Bar:hello}
```

---

### Merging configurations from multiple sources

```go
type Config struct {
    Foo   string   `json:"foo"`
    Bar   string   `json:"bar"`
    Sizes []string `xml:"size"`
}

source1 := strings.NewReader(`{"foo": "hello", "bar": "world"}`)
source2 := strings.NewReader(`
<sizes>
    <size>small</size>
    <size>regular</size>
    <size>large</size>
</sizes>`)

cfg, err := config.Multi[Config]().
    Add(config.FromReader(source1)).
    Add(config.FromReader(source2), config.WithDecoder(xml.NewDecoder)).
    AllOf()
if err != nil {
    panic(err)
}

fmt.Printf("%+v\n", cfg) // Output: {Foo:hello Bar:world Sizes:[small regular large]}
```

> The `AllOf` method requires successful reading from all sources. If at least one successful source is needed, use the `OneOf` method.

---

## 🛠️ API

### Interfaces

1. **`ConfigProvider`**:

   - Responsible for retrieving configuration from a specific source.
   - Implementations:
     - `FromCmdline`
     - `FromConsul`
     - `FromEnv`
     - `FromFile`
     - `FromReader`
     - `NoProvider`

2. **`Decoder`**:
   - Responsible for converting configuration into the desired structure.
   - Supports custom decoders.
