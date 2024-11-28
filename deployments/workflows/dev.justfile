[no-cd]
[doc("Start the GraphQL server")]
start secret="secret":
    @go run main.go start -s {{ secret }} --log-level debug -d "app.db?_pragma=foreign_keys(1)"
