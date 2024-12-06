#!/usr/bin/env just --justfile
args := "-v $(go list ./... | grep -v 'ent/ent' | grep -v 'gql/' | grep -v 'mocks/' | grep -v '**/cmd/*' | grep -v '**/internal/api/' )"
output := "-coverprofile coverage.txt"

[no-cd]
[doc("Run all tests with detailed output and generate a coverage profile (coverage.txt)")]
tests:
    @go test {{ args }} {{ output }}

[no-cd]
[doc("Display the test coverage report in HTML format using the generated coverage profile")]
coverage: tests
    @go tool cover -html=coverage.txt

[no-cd]
[doc("Show total test coverage")]
coverage-total:
    @go test {{ args }} {{ output }} > /dev/null 2>&1
    @go tool cover -func=coverage.txt | grep total

[no-cd]
[doc("Generate a coverage SVG visualization from the coverage profile")]
coverage-svg: tests
    @$HOME/go/bin/go-cover-treemap -coverprofile coverage.txt > out.svg


[no-cd]
[doc("Download all dependencies")]
deps:
    @go install github.com/vektra/mockery/v2@v2.49.1
    @go install github.com/nikolaydubina/go-cover-treemap@latest
