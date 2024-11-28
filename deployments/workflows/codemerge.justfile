[no-cd]
[doc("Ensure codemerge is installed")]
check-installs:
    @command -v codemerge >/dev/null || @cargo install codemerge


[no-cd]
[doc("Merge files using codemerge")]
merge: check-installs
    @codemerge merge -i "./app.db" -i "./changescout/internal/api/gql/generated/*" -i "./changescout/internal/infrastructure/database/ent/ent/*" -i ".idea/" -i "./.idea/dataSources/b800681d-81cb-4b99-9746-3b9d034b8592.xml" -i "./changescout/internal/**/mocks/*.go" -i "./changescout/internal/**/*_test.go" -i "*.md" -i "*.sum"

[no-cd]
[doc("Generate tokens using codemerge")]
tokens: check-installs
    @codemerge tokens -i "./app.db" -i "./changescout/internal/api/gql/generated/*" -i "./changescout/internal/infrastructure/database/ent/ent/*" -i ".idea/" -i "./.idea/dataSources/b800681d-81cb-4b99-9746-3b9d034b8592.xml" -i "./changescout/internal/**/mocks/*.go" -i "./changescout/internal/**/*_test.go" -i "*.md" -i "*.sum"
