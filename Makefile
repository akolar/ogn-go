

all: fmt vet test

test:
	@go test ./...

vet:
	@if ! $$(go vet ./...); then \
		echo "Go vet failed"; \
	fi
	

fmt:
	@failed=$$(gofmt -l .)
	@if [ -n "$$failed" ]; then \
		echo "Gofmt failed for:"; \
		echo $$failed; \
	fi
