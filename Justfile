
# runs the golangci-lint linter
lint:
    golangci-lint --verbose run

# runs the golangci-lint linter and fixes the issues
lint-fix:
    golangci-lint --verbose run --fix

# runs tests in all packages
test:
	go test -v ./...
