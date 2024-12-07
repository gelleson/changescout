target := "changescout/internal/infrastructure/database/ent/ent/schema/"

[no-cd]
[doc("Generate GraphQL code")]
gql:
    @echo "Generating code from schema.graphqls"
    @go run -mod=mod github.com/99designs/gqlgen gen
    @echo "Generated code from schema.graphqls"


[no-cd]
[doc("Generate Ent schema code")]
ent name:
    @echo "Generating new schema for {{ titlecase(name) }}"
    @go run -mod=mod entgo.io/ent/cmd/ent new {{ titlecase(name) }} --target {{ target }}
