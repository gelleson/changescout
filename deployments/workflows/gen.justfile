[no-cd]
[doc("Generate GraphQL code")]
gql:
    @echo "Generating code from schema.graphqls"
    @go run -mod=mod github.com/99designs/gqlgen gen
    @echo "Generated code from schema.graphqls"