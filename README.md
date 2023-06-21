# Butler

Butler enables you load your configuration file onto a configuration struct. Successfully reading your configuration variables of any Go type.

## Installation

As a library

```shell
go get github.com/nade-harlow/butler
```

## Usage

### Reading From .env file

Load your application configuration from your `.env` file located at any path of your project:

```shell
PORT=8000
S3_BUCKET=YOURS3BUCKET
SECRET_KEY=YOURSECRETKEYGOESHERE
```

Then in your Go app you can do something like

```go
package main

import (
    "log"

    "github.com/nade-harlow/butler"
)

type Config struct {
    Port      int32  `env:"port"`
    SecretKey string `env:"secret_key"`
    S3Bucket  string `env:"s3_bucket"`
}

func main() {
  config := &Config{}

  err := butler.LoadConfig(&config, "./env")
  if err != nil {
    log.Fatalf("Error loading configuration: %v", err)
  }

  // now do something with config
}
```

### Reading From .yaml or .yml file

Load your application configuration from your `.yaml` or `yml` file located at any path of your project:

```yaml
PORT:8000
S3_BUCKET:YOURS3BUCKET
SECRET_KEY:YOURSECRETKEYGOESHERE
```

Then in your Go app you can do something like

```go
package main

import (
    "log"

    "github.com/nade-harlow/butler"
)

type Config struct {
    Port      int32  `env:"port"`
    SecretKey string `env:"secret_key"`
    S3Bucket  string `env:"s3_bucket"`
}

func main() {
  config := &Config{}

  err := butler.LoadConfig(&config, "./config.yaml")
  if err != nil {
    log.Fatalf("Error loading configuration: %v", err)
  }

  // now do something with config
}
```