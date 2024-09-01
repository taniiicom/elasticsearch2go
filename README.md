# elasticsearch2go

[English](https://github.com/taniiicom/elasticsearch2go/blob/main/README.md)
| [日本語](https://github.com/taniiicom/elasticsearch2go/blob/main/README.ja.md)

`elasticsearch2go` is a tool that automatically generates Go structs based on Elasticsearch mapping definitions. It helps automate the process of generating Go structs from Elasticsearch JSON schemas, ensuring consistency and reducing manual coding efforts.

https://github.com/user-attachments/assets/33f1d144-b4b5-44cb-a399-9de3d7e1f522

## Special Thanks!

This project was originally developed during my time with the minimo division at MIXI, Inc. ([@mixigroup](https://github.com/mixigroup)) as part of our MLOps automation efforts. I am deeply grateful to the team for their generosity in allowing this package to be released as OSS. Thank you for your support!

web: [https://mixi.co.jp/](https://mixi.co.jp/)

## Installation

To install this package, use the following command:

```bash
go get github.com/taniiicom/elasticsearch2go
```

## Usage

Below is an example of how to use this package to generate Go structs from Elasticsearch mapping definitions.

### Command Line Exec

e.g.

```bash
go run github.com/taniiicom/elasticsearch2go/cmd \
    --in example/elasticsearch/cafe-mapping.json \
    --out example/infrastructure/datamodel/searchmodel/cafe.gen.go \
    --struct CafeDocJson \
    --package searchmodel \
```

### Go Exec

```go
package gen

import (
	"log"

	"github.com/taniiicom/elasticsearch2go" // import the package
)

func main() {
	// required arguments
	inputPath := "example/elasticsearch/cafe-mapping.json"
	outputPath := "example/infrastructure/datamodel/searchmodel/cafe.gen.go"
	packageName := "searchmodel"
	structName := "CafeDocJson"

	// optional arguments
	opts := &elasticsearch2go.GeneratorOptions{
		InitClassName:      nil, // optional
		TypeMappingPath:    nil, // optional
		ExceptionFieldPath: nil, // optional
		ExceptionTypePath:  nil, // optional
		SkipFieldPath:      nil, // optional
		FieldCommentPath:   nil, // optional
		TmplPath:           nil, // optional
	}

	// generate datamodel
	err := elasticsearch2go.GenerateDatamodel(inputPath, outputPath, packageName, structName, opts)
	if err != nil {
		log.Fatalf("Failed to generate data model: %v", err)
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

This package is currently maintained by [@taniiicom](https://github.com/taniiicom) (Taniii.com). If you encounter any issues, please report them using the GitHub Issues.
