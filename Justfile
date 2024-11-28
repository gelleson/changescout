gen:
    @echo "Generating code from schema.graphqls"
    @go run -mod=mod github.com/99designs/gqlgen gen
    @echo "Generated code from schema.graphqls"

dev:
    @go run main.go start -s secret --log-level debug -d "app.db?_pragma=foreign_keys(1)"

test:
    @go test -v ./...

codemerge command="merge":
    @codemerge {{ command }} -i "./app.db" -i "./changescout/internal/api/gql/generated/*" -i "./changescout/internal/infrastructure/database/ent/ent/*" -i ".idea/" -i "./.idea/dataSources/b800681d-81cb-4b99-9746-3b9d034b8592.xml" -i "./changescout/internal/**/mocks/*.go" -i "./changescout/internal/**/*_test.go" -i "*.md" -i "*.sum"
