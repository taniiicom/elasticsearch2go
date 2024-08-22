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
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties,omitempty"`
}

type Mappings struct {
	Properties map[string]Property `json:"properties"`
}

type ElasticsearchMapping struct {
	Mappings Mappings `json:"mappings"`
}

const structTemplateWithWrapper = `package {{.PackageName}}

type {{.InitClassName}} struct {
	{{.StructName}}
}

{{.StructDefinitions}}
`

const structTemplateWithoutWrapper = `package {{.PackageName}}

{{.StructDefinitions}}
`

type Field struct {
	FieldName    string
	FieldType    string
	JSONName     string
	FieldComment string
}

type StructData struct {
	PackageName       string
	InitClassName     string
	StructName        string
	StructDefinitions string
}

// GoTypeMap holds the mapping from Elasticsearch types to Go types.
var GoTypeMap map[string]string
var FieldExceptions map[string]string
var TypeExceptions map[string]string
var SkipFields map[string]bool
var FieldComments map[string]string

// StructNameTracker to avoid generating duplicate struct names
var StructNameTracker map[string]bool

func main() {
	inputPath := flag.String("in", "", "Input JSON schema file (including file name)")
	outputPath := flag.String("out", "", "Output Go file (including file name)")
	packageName := flag.String("package", "searchmodel", "Name of the Go package")
	structName := flag.String("struct", "GeneratedStruct", "Name of the generated Go struct")
	initClassName := flag.String("init", "", "Name of the initial wrapper struct (optional)")
	typeMappingPath := flag.String("type-mapping", "", "Path to JSON file specifying Elasticsearch to Go type mapping")
	exceptionFieldPath := flag.String("exception-field", "", "Path to JSON file specifying exceptions for field names")
	exceptionTypePath := flag.String("exception-type", "", "Path to JSON file specifying exceptions for field types")
	skipFieldPath := flag.String("skip-field", "", "Path to JSON file specifying fields to skip")
	fieldCommentPath := flag.String("field-comment", "", "Path to JSON file specifying comments for fields")
	tmplPath := flag.String("tmpl", "", "Path to custom Go template file")
	flag.Parse()

	if *inputPath == "" || *outputPath == "" || *structName == "" || *packageName == "" {
		log.Fatalf("All --in, --out, --struct, and --package must be specified")
	}

	// initialize StructNameTracker
	StructNameTracker = make(map[string]bool)

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

	// load skip fields if provided
	if *skipFieldPath != "" {
		loadSkipFields(*skipFieldPath)
	} else {
		SkipFields = make(map[string]bool)
	}

	// load field comments if provided
	if *fieldCommentPath != "" {
		loadFieldComments(*fieldCommentPath)
	} else {
		FieldComments = make(map[string]string)
	}

	// load custom template if provided
	var tmpl *template.Template
	var err error
	if *tmplPath != "" {
		tmpl, err = template.ParseFiles(*tmplPath)
		if err != nil {
			log.Fatalf("Failed to load template file %s: %v", *tmplPath, err)
		}
	} else {
		// choose default template based on the presence of initClassName
		if *initClassName != "" {
			tmpl, err = template.New("structWithWrapper").Parse(structTemplateWithWrapper)
		} else {
			tmpl, err = template.New("structWithoutWrapper").Parse(structTemplateWithoutWrapper)
		}
		if err != nil {
			log.Fatalf("Error parsing template: %v", err)
		}
	}

	processFile(*inputPath, *outputPath, *packageName, *structName, *initClassName, tmpl)
}

func processFile(inputPath, outputPath, packageName, structName, initClassName string, tmpl *template.Template) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatalf("Failed to read file %s: %v", inputPath, err)
	}

	var esMapping ElasticsearchMapping
	err = json.Unmarshal(data, &esMapping)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON from file %s: %v", inputPath, err)
	}

	structDefinitions := generateStructDefinitions(structName, esMapping.Mappings.Properties)

	structData := StructData{
		PackageName:       packageName,
		InitClassName:     initClassName,
		StructName:        structName,
		StructDefinitions: structDefinitions,
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

func generateStructDefinitions(structName string, properties map[string]Property) string {
	var structDefs strings.Builder

	generateStruct(&structDefs, structName, properties)

	return structDefs.String()
}

func generateStruct(structDefs *strings.Builder, structName string, properties map[string]Property) {
	// check if the struct has already been generated
	if _, exists := StructNameTracker[structName]; exists {
		return
	}

	// mark this struct as generated
	StructNameTracker[structName] = true

	fields := []Field{}
	nestedStructs := []string{}

	for name, prop := range properties {
		// skip fields that are in the SkipFields map
		if _, skip := SkipFields[name]; skip {
			continue
		}

		fieldName := mapElasticsearchFieldToGoField(name)
		var fieldType string

		if prop.Type == "object" || prop.Type == "nested" {
			// check if the type has a custom exception
			if customType, exists := TypeExceptions[name]; exists {
				var nestedStructName string
				fieldType = customType
				if strings.HasPrefix(fieldType, "*") {
					nestedStructName = fieldType[1:] // "*" を取り除く
				} else if strings.HasPrefix(fieldType, "[]") {
					nestedStructName = fieldType[2:] // "[]" を取り除く
				} else {
					nestedStructName = fieldType
				}
				nestedStructs = append(nestedStructs, generateStructDefinitions(nestedStructName, prop.Properties))
			} else {
				nestedStructName := toPascalCase(name)
				fieldType = "*" + nestedStructName
				nestedStructs = append(nestedStructs, generateStructDefinitions(nestedStructName, prop.Properties))
			}
		} else {
			fieldType = mapElasticsearchTypeToGoType(name, prop.Type)
		}

		fieldComment := mapElasticsearchFieldToComment(name)

		fields = append(fields, Field{
			FieldName:    fieldName,
			FieldType:    fieldType,
			JSONName:     name,
			FieldComment: fieldComment,
		})
	}

	// sort fields alphabetically
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].FieldName < fields[j].FieldName
	})

	// generate struct definition
	structDefs.WriteString(fmt.Sprintf("type %s struct {\n", structName))
	for _, field := range fields {
		if field.FieldComment != "" {
			structDefs.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\"` // %s\n", field.FieldName, field.FieldType, field.JSONName, field.FieldComment))
		} else {
			structDefs.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\"`\n", field.FieldName, field.FieldType, field.JSONName))
		}
	}
	structDefs.WriteString("}\n\n")

	// append nested structs
	for _, nestedStruct := range nestedStructs {
		structDefs.WriteString(nestedStruct)
	}
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

	return toPascalCase(esFieldName)
}

func mapElasticsearchFieldToComment(esFieldName string) string {
	// check if the field has a custom comment
	if comment, exists := FieldComments[esFieldName]; exists {
		return comment
	}

	return ""
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

func loadSkipFields(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read skip fields file %s: %v", filePath, err)
	}

	err = json.Unmarshal(data, &SkipFields)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON from skip fields file %s: %v", filePath, err)
	}
}

func loadFieldComments(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read field comments file %s: %v", filePath, err)
	}

	err = json.Unmarshal(data, &FieldComments)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON from field comments file %s: %v", filePath, err)
	}
}

func toCamelCase(s string) string {
	caser := cases.Title(language.Und) // or: `language.English`
	parts := strings.Split(s, "_")
	for i, part := range parts {
		parts[i] = caser.String(part)
	}
	parts[0] = strings.ToLower(parts[0])
	return strings.Join(parts, "")
}

func toPascalCase(s string) string {
	caser := cases.Title(language.Und)
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
