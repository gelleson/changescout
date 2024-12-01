args := '-i "./app.db" -i "./changescout/internal/api/gql/generated/*" -i "./changescout/internal/infrastructure/database/ent/ent/*" -i ".idea/" -i "./.idea/dataSources/b800681d-81cb-4b99-9746-3b9d034b8592.xml" -i "./changescout/internal/**/mocks/*.go" -i "./changescout/internal/**/*_test.go" -i "*.md" -i "*.sum" -i "./web/node_modules/*" -i "./web/bun.lockb"'

[doc("Ensure codemerge is installed")]
[no-cd]
check-installs:
    @command -v codemerge >/dev/null || @cargo install codemerge

[doc("Merge files using codemerge")]
[no-cd]
merge: check-installs
    @codemerge merge {{ args }}

[doc("Generate tokens using codemerge")]
[no-cd]
tokens: check-installs
    @codemerge tokens {{ args }}
