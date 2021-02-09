# Kafka Message Bus

## How to Use

### Use at Other Projects

Since this library does not reside in a Go Proxy repository, you need to use it with Go module's `replace` directive. You need to set your `go.mod` file as following:

```
require (
    ...
    kata.ai/messagebus-golang-kafka v1.1.0
)

replace kata.ai/messagebus-golang-kafka => ../shared-libs/golang/messagebus-kafka/1.1.0
```

The replace directive is filled with this library's module name and its relative path to your project.

### Use Avro Schema

This library only accept [Gogen-avro](https://github.com/actgardner/gogen-avro) schema. If you have an Avro schema in a `avsc` file, you need to convert it to Gogen-avro schema with these steps:

1. Install Gogen-avro with its [installation guide](https://github.com/actgardner/gogen-avro#installation)
2. Convert your `avsc` file to Go file with Gogen-avro command:
    ```bash
    gogen-avro --package=<generated-file-package-name> schemas <schema-directory>
    ```
    For example, if you have an Avro schema in `schemas/prediction_log.avsc`, you have to execute this command:

    ```bash
    gogen-avro --package=schemas schemas schemas/prediction_log.avsc
    ```
    It will create a Go file in `schemas` directory and you can use it with this library.

### Examples

#### [Consumer Example](./1.1.0/example/consumer_example/consume_example.go)

#### [Producer Example](./1.1.0/example/publisher_example/publish_example.go)