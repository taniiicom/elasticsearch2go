package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Property struct {
	Type string `json:"type"`
}

type Mappings struct {
	Properties map[string]Property `json:"properties"`
}

type ElasticsearchMapping struct {
	Mappings Mappings `json:"mappings"`
}

const structTemplate = `package {{.PackageName}}

type {{.InitClassName}} struct {
	{{.StructName}}
}

type {{.StructName}} struct {
{{- range .Fields}}
	{{.FieldName}} {{.FieldType}} ` + "`json:\"{{.JSONName}}\"`" + `
{{- end}}
}
`

type Field struct {
	FieldName string
	FieldType string
	JSONName  string
}

type StructData struct {
	PackageName   string
	InitClassName string
	StructName    string
	Fields        []Field
}

// GoTypeMap holds the mapping from Elasticsearch types to Go types.
var GoTypeMap map[string]string
var FieldExceptions map[string]string
var TypeExceptions map[string]string

func main() {
	inputPath := flag.String("in", "", "Input JSON schema file (including file name)")
	outputPath := flag.String("out", "", "Output Go file (including file name)")
	packageName := flag.String("package", "searchmodel", "Name of the Go package")
	structName := flag.String("struct", "GeneratedStruct", "Name of the generated Go struct")
	initClassName := flag.String("init", "", "Name of the initial wrapper struct (optional)")
	typeMappingPath := flag.String("type-mapping", "", "Path to JSON file specifying Elasticsearch to Go type mapping")
	exceptionFieldPath := flag.String("exception-field", "", "Path to JSON file specifying exceptions for field names")
	exceptionTypePath := flag.String("exception-type", "", "Path to JSON file specifying exceptions for field types")
	flag.Parse()

	if *inputPath == "" || *outputPath == "" || *structName == "" || *packageName == "" {
		log.Fatalf("All --in, --out, --struct, and --package must be specified")
	}

	if *initClassName == "" {
		*initClassName = fmt.Sprintf("%sWrapper", *structName)
	}

	// load custom type mapping if provided
	if *typeMappingPath != "" {
		loadTypeMapping(*typeMappingPath)
	} else {
		// default mapping
		GoTypeMap = map[string]string{
			"integer":   "*uint64",
			"float":     "*float64",
			"boolean":   "bool",
			"text":      "*string",
			"keyword":   "*string",
			"date":      "*time.Time",
			"geo_point": "*GeoPoint",
			"object":    "*map[string]interface{}",
			"nested":    "[]interface{}",
		}
	}

	// load field exceptions if provided
	if *exceptionFieldPath != "" {
		loadFieldExceptions(*exceptionFieldPath)
	} else {
		FieldExceptions = make(map[string]string)
	}

	// load type exceptions if provided
	if *exceptionTypePath != "" {
		loadTypeExceptions(*exceptionTypePath)
	} else {
		TypeExceptions = make(map[string]string)
	}

	processFile(*inputPath, *outputPath, *packageName, *structName, *initClassName)
}

func processFile(inputPath, outputPath, packageName, structName, initClassName string) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatalf("Failed to read file %s: %v", inputPath, err)
	}

	var esMapping ElasticsearchMapping
	err = json.Unmarshal(data, &esMapping)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON from file %s: %v", inputPath, err)
	}

	fields := []Field{}
	for name, prop := range esMapping.Mappings.Properties {
		fieldName := mapElasticsearchFieldToGoField(name)
		fieldType := mapElasticsearchTypeToGoType(name, prop.Type)
		fields = append(fields, Field{
			FieldName: fieldName,
			FieldType: fieldType,
			JSONName:  name,
		})
	}

	// sort field names alphabetically
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].FieldName < fields[j].FieldName
	})

	structData := StructData{
		PackageName:   packageName,
		InitClassName: initClassName,
		StructName:    structName,
		Fields:        fields,
	}

	tmpl, err := template.New("struct").Parse(structTemplate)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, structData)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	err = os.WriteFile(outputPath, buf.Bytes(), 0644)
	if err != nil {
		log.Fatalf("Failed to write output file %s: %v", outputPath, err)
	}

	fmt.Printf("Generated Go struct for %s and saved to %s\n", inputPath, outputPath)
}

func mapElasticsearchTypeToGoType(name, esType string) string {
	// check if the type has a custom exception
	if customType, exists := TypeExceptions[name]; exists {
		return customType
	}

	goType, exists := GoTypeMap[esType]
	if !exists {
		goType = "interface{}"
	}

	return goType
}

func mapElasticsearchFieldToGoField(esFieldName string) string {
	// check if the field has a custom exception
	if customFieldName, exists := FieldExceptions[esFieldName]; exists {
		return customFieldName
	}

	return toCamelCase(esFieldName)
}

func loadTypeMapping(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read type mapping file %s: %v", filePath, err)
	}

	err = json.Unmarshal(data, &GoTypeMap)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON from type mapping file %s: %v", filePath, err)
	}
}

func loadFieldExceptions(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read field exception file %s: %v", filePath, err)
	}

	err = json.Unmarshal(data, &FieldExceptions)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON from field exception file %s: %v", filePath, err)
	}
}

func loadTypeExceptions(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read type exception file %s: %v", filePath, err)
	}

	err = json.Unmarshal(data, &TypeExceptions)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON from type exception file %s: %v", filePath, err)
	}
}

func toCamelCase(s string) string {
	caser := cases.Title(language.Und) // or: `language.English`
	parts := strings.Split(s, "_")
	for i, part := range parts {
		parts[i] = caser.String(part)
	}
	return strings.Join(parts, "")
}

// GeoPoint struct for handling geo_point type in Elasticsearch
type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
