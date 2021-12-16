PROJECT_DIR := ${CURDIR}

TEST_DIR := ${PROJECT_DIR}/internal


## lint: run go liners
lint:
	golangci-lint run


## test-func: run go test func
test:
	cd ${TEST_DIR}
	go test -coverpkg=./... -coverprofile=cover ./... && cat cover | grep -v "logs" | grep -v "mock" | grep -v "database" | grep -v "cmd" | grep -v  "easyjson" | grep -v "proto" | grep -v "config" | grep -v "metrics" | grep -v "error" | grep -v "mocks" > cover.out && go tool cover -func=cover.out
	rm -f cover
	rm -f cover.out
