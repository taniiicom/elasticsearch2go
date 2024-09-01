go run github.com/taniiicom/elasticsearch2go \
    --in example/elasticsearch/cafe-mapping.json \
    --out example/infrastructure/datamodel/searchmodel/cafe.gen.go \
    --struct CafeDocJson \
    --package searchmodel \
    --type-mapping example/script/custom/type-mapping.json \
    --tmpl example/script/custom/custom-template.tmpl
