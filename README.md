# elasticsearch2go

[English](https://github.com/taniiicom/elasticsearch2go/blob/main/README.md)
| [日本語](https://github.com/taniiicom/elasticsearch2go/blob/main/README.ja.md)

`elasticsearch2go` is a tool that automatically generates Go structs based on Elasticsearch mapping definitions. It helps automate the process of generating Go structs from Elasticsearch JSON schemas, ensuring consistency and reducing manual coding efforts.

## Installation

To install this package, use the following command:

```bash
go get github.com/taniiicom/elasticsearch2go
```

## Usage

Below is an example of how to use this package to generate Go structs from Elasticsearch mapping definitions.

### Command Line Usage

```bash
elasticsearch2go --in=mapping.json --out=model.go --package=mypackage --struct=MyStruct
```

### Code Example

```go
package main

import (
    "log"
    "github.com/taniiicom/elasticsearch2go"
)

func main() {
    err := elasticsearch2go.GenerateStructs(
        "mapping.json",
        "model.go",
        "mypackage",
        "MyStruct",
        "MyWrapperStruct",
        "custom_mapping.json",
        "custom_field_exceptions.json",
        "custom_type_exceptions.json",
        "skip_fields.json",
        "field_comments.json",
        "custom_template.tmpl",
    )
    if err != nil {
        log.Fatalf("Failed to generate structs: %v", err)
    }
}
```

## Command-line Options

This package offers the following command-line options:

- `--in`: **(required)** Specifies the path to the input JSON schema file.
- `--out`: **(required)** Specifies the path to the output Go file.
- `--package`: **(required)** Specifies the package name for the generated Go file.
- `--struct`: **(required)** Specifies the name of the generated struct.
- `--init`: Specifies the name of the initial wrapper struct (optional).
- `--type-mapping`: Specifies the path to a JSON file that maps Elasticsearch types to Go types (optional).
- `--exception-field`: Specifies the path to a JSON file that defines exceptions for field names (optional).
- `--exception-type`: Specifies the path to a JSON file that defines exceptions for field types (optional).
- `--skip-field`: Specifies the path to a JSON file that defines fields to skip during struct generation (optional).
- `--field-comment`: Specifies the path to a JSON file that adds comments to fields (optional).
- `--tmpl`: Specifies the path to a custom Go template file (optional).

## Customization

This package supports various customization options. Below are some examples of how you can customize the output.

### Custom Type Mapping

To map Elasticsearch types to specific Go types, use a `custom_mapping.json` file. Example:

```json
{
  "text": "*string",
  "integer": "int"
}
```

This configuration maps the `text` type to `*string` and the `integer` type to `int`.

### Field Name Exceptions

To apply specific transformations to field names, use a `custom_field_exceptions.json` file. Example:

```json
{
  "my_field": "MyCustomField"
}
```

This will map `my_field` to the Go field name `MyCustomField`.

### Skipping Fields

To skip specific fields during struct generation, use a `skip_fields.json` file. Example:

```json
{
  "unnecessary_field": true
}
```

This configuration will exclude the `unnecessary_field` from the generated struct.

## Features

### Support for Nested Struct Generation

Fields defined as `"type": "nested"` in Elasticsearch are treated as sub-structures when generating Go structs. These sub-structures are generated under the main struct within the same output file. This allows the nested structure to be directly reflected in the Go code.

### Flexible Customization

You can specify exceptions for Elasticsearch to Go mappings, including field names and types, using a JSON file. Exceptions can be applied to both field names and types.

Additionally, you can specify a JSON file to skip certain properties. This allows you to exclude unnecessary fields from being generated.

### Adding Comments to Fields

You can provide a JSON file to add comments to Elasticsearch fields. This feature allows you to include comments in the generated Go structs.

### Custom Template Support

To customize the format of the generated files, you can specify a custom template file. This feature gives you control over the formatting of the output Go code.

### [Todo] Specifying Field Order in Structs

Currently, the fields in the generated structs are ordered alphabetically. In future versions, it will be possible to maintain the original field order as defined in the Elasticsearch schema.

## Contributing

Contributions to this package are welcome. To contribute, please follow these steps:

1. Fork this repository.
2. Create a new branch (`git checkout -b feature/your-feature-name`).
3. Commit your changes (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin feature/your-feature-name`).
5. Open a Pull Request on GitHub.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

## Maintenance and Support

This package is currently maintained by [Your Name or Organization]. If you encounter any issues, please report them using the GitHub Issues.
