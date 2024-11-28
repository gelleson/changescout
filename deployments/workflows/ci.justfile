#!/usr/bin/env just --justfile

[no-cd]
[doc("Run all tests with detailed output and generate a coverage profile (coverage.txt)")]
tests: deps
    @go test -v $(go list ./... | grep -v 'ent/ent' | grep -v 'gql/' | grep -v 'mocks/' ) -coverprofile coverage.txt

[no-cd]
[doc("Display the test coverage report in HTML format using the generated coverage profile")]
coverage: tests
    @go tool cover -html=coverage.txt

[no-cd]
[doc("Generate a coverage SVG visualization from the coverage profile")]
coverage-svg: tests
    @$HOME/go/bin/go-cover-treemap -coverprofile coverage.txt > out.svg
