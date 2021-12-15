PROJECT_DIR := ${CURDIR}

TEST_DIR := ${PROJECT_DIR}/internal


## lint: run go liners
lint:
	golangci-lint run


## test-func: run go test func
test:
	cd ${TEST_DIR}
	go test -coverpkg=./... -coverprofile=cover ./... && cat cover | grep -v "mock" | grep -v  "easyjson" | grep -v "proto" > cover.out && go tool cover -func=cover.out
	rm -rf cover
	rm -rf cover.out