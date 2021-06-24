
test: ## run tests
	cd test && ./test.sh

lint: ## lint code using local version of golangci-lint
	golangci-lint run

help:
	@grep -E '^[\.a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test lint help
