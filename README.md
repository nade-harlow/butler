# Butler

Butler is a configuration management library for Go that enables you to load your configuration variables from different file formats, such as `.env`, `.yaml`, or `.yml`, into a configuration struct.

## Installation

Install Butler as a library using `go get`:

```shell
go get github.com/nade-harlow/butler
```

## Usage

### To read From .env file follow the steps below: 
- Create a .env file in your project directory and populate it with your configuration variables:
- Load your application configuration from your `.env` file located at any path of your project:

```shell
PORT=8000
S3_BUCKET=YOURS3BUCKET
SECRET_KEY=YOURSECRETKEYGOESHERE
```

In your Go application, import Butler:

```go
package main

import (
    "log"

    "github.com/nade-harlow/butler"
)
- Define a struct that represents your configuration:

type Config struct {
    Port      int32  `env:"port"`
    SecretKey string `env:"secret_key"`
    S3Bucket  string `env:"s3_bucket"`
}

- Load the configuration from the .env file:

func main() {
  config := &Config{}

  err := butler.LoadConfig(&config, "./.env")
  if err != nil {
    log.Fatalf("Error loading configuration: %v", err)
  }

  // now do something with config
}
```

### Reading From .yaml or .yml file
-  Create a .yaml or .yml file in your project directory and define your configuration variables:
- Load your application configuration from your `.yaml` or `yml` file located at any path of your project:

```yaml
PORT:8000
S3_BUCKET:YOURS3BUCKET
SECRET_KEY:YOURSECRETKEYGOESHERE
```

Follow the same steps as for reading from a .env file, but replace the file path with the path to your YAML file:

```go
err := butler.LoadConfig(config, "./config.yaml")

func main() {
  config := &Config{}

  err := butler.LoadConfig(&config, "./config.yaml")
  if err != nil {
    log.Fatalf("Error loading configuration: %v", err)
  }

  // now do something with config
}
```
## Contributing
We welcome contributions from the community! To contribute to Butler, please follow these steps:

1. Fork the repository and clone it to your local machine.
2. Create a new branch for your feature or bug fix.
3. Make your changes and ensure that the tests pass.
4. Commit your changes and push them to your forked repository.
5. Submit a pull request with a detailed description of your changes.

## License

This project is licensed under the MIT License. See the LICENSE file for details.