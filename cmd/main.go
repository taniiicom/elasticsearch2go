package main

import (
	"flag"
	"log"

	. "github.com/taniiicom/elasticsearch2go"
)

func main() {
	// required arguments
	inputPath := flag.String("in", "", "Input JSON schema file (including file name)")
	outputPath := flag.String("out", "", "Output Go file (including file name)")
	packageName := flag.String("package", "searchmodel", "Name of the Go package")
	structName := flag.String("struct", "GeneratedStruct", "Name of the generated Go struct")

	// optional arguments
	initClassName := flag.String("init", "", "Name of the initial wrapper struct (optional)")
	typeMappingPath := flag.String("type-mapping", "", "Path to JSON file specifying Elasticsearch to Go type mapping")
	exceptionFieldPath := flag.String("exception-field", "", "Path to JSON file specifying exceptions for field names")
	exceptionTypePath := flag.String("exception-type", "", "Path to JSON file specifying exceptions for field types")
	skipFieldPath := flag.String("skip-field", "", "Path to JSON file specifying fields to skip")
	fieldCommentPath := flag.String("field-comment", "", "Path to JSON file specifying comments for fields")
	tmplPath := flag.String("tmpl", "", "Path to custom Go template file")

	flag.Parse()

	// validate required arguments
	if *inputPath == "" || *outputPath == "" || *structName == "" || *packageName == "" {
		log.Fatalf("All --in, --out, --struct, and --package must be specified")
	}

	// set up generator options
	opts := &GeneratorOptions{
		InitClassName:      nullableString(initClassName),
		TypeMappingPath:    nullableString(typeMappingPath),
		ExceptionFieldPath: nullableString(exceptionFieldPath),
		ExceptionTypePath:  nullableString(exceptionTypePath),
		SkipFieldPath:      nullableString(skipFieldPath),
		FieldCommentPath:   nullableString(fieldCommentPath),
		TmplPath:           nullableString(tmplPath),
	}

	// generate datamodel
	err := GenerateDatamodel(*inputPath, *outputPath, *packageName, *structName, opts)
	if err != nil {
		log.Fatalf("Failed to generate data model: %v", err)
	}
}

// nullableString is a helper function to treat flag.String values as nullable.
func nullableString(flagValue *string) *string {
	if *flagValue == "" {
		return nil
	}
	return flagValue
}
