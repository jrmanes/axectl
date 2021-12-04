PROJECT_NAME=piktoctl

# Go
.PHYONY: run build test test_cover get docs clean
run:
	go run ./main.go

build: clean
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/${PROJECT_NAME} ./

build_copy:
	rm ./bin/${PROJECT_NAME} || true
	rm ~/.local/bin/${PROJECT_NAME} || true
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/${PROJECT_NAME} ./
	mkdir -p ~/.local/bin/
	cp ./bin/${PROJECT_NAME} ~/.local/bin/

clean:
	rm ./bin/${PROJECT_NAME} || true
	rm ~/.local/bin/${PROJECT_NAME} || true

test:
	go test ./... -v -cover .

test_cover:
	go test ./... -v -coverprofile cover.out
	go tool cover -func ./cover.out | grep total | awk '{print $3}'

get:
	go get ./...

docs:
	godoc -http=:6060

.PHYONY: status scan
status:
	ls -ltra
	go run ./main.go --scan status

scan:
	go run ./main.go sonar --scan -o "soyuntest" -p "test1,test2"


