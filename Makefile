SHELL=/bin/bash
PROJECT_NAME=piktoctl
######################
# Go
.PHYONY: run build build_copy test test_cover get docs
run:
	go run ./main.go

run_sonar:
	go run ./main.go sonar

build: clean
	mkdir -p ./bin/ ./bin/m1
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/${PROJECT_NAME} ./
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./bin/m1/${PROJECT_NAME} ./

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

######################
# common
.PHYONY: status scan
status:
	ls -ltra
	go run ./main.go --scan status

scan:
	go run ./main.go sonar --scan -o "soyuntest" -p "test1,test2"

######################
# Vagrant
.PHYONY: vagrant_up vagrant_rm vagrant_ssh vagrant_ssh vagrant_reload vagrant_remote_deploy
vagrant_up:
	VAGRANT_VAGRANTFILE=./infra/Vagrantfile vagrant up

vagrant_rm:
	VAGRANT_VAGRANTFILE=./infra/Vagrantfile vagrant destroy --force && rm -fr .vagrant/

vagrant_ssh:
	VAGRANT_VAGRANTFILE=./infra/Vagrantfile vagrant ssh

vagrant_reload:
	VAGRANT_VAGRANTFILE=./infra/Vagrantfile vagrant reload

vagrant_remote_deploy: build_copy
	scp -P 2222 ./bin/piktoctl vagrant@127.0.0.1:.

######################
# Docker
.PHYONY: vagrant_up vagrant_rm vagrant_ssh vagrant_ssh vagrant_reload vagrant_remote_deploy
docker_build: build
	docker build -f ./infra/Dockerfile . -t ${PROJECT_NAME}:latest

docker_start: docker_build
	docker run -t --rm -v ${PWD}/bin/:/app ${PROJECT_NAME}:latest ls /app

docker_status: docker_build
	docker run -t --rm -v ${PWD}/bin/:/app ${PROJECT_NAME}:latest ls -ltra /app/

docker_piktoctl: docker_build
	docker run -t --rm -v ${PWD}/bin/:/app ${PROJECT_NAME}:latest bash -c /app/piktoctl

docker_piktoctl_install_sudo: docker_build
	docker run -t --rm -v ${PWD}/bin/:/app ${PROJECT_NAME}:latest /bin/bash -c "apt update && apt install -y sudo"

docker_piktoctl_sonar: docker_build
	docker run -t --rm -v ${PWD}/bin/:/app ${PROJECT_NAME}:latest /bin/bash -c "/app/piktoctl sonar"

docker_piktoctl_sonar_install: docker_build
	docker run -t --rm -v ${PWD}/bin/:/app ${PROJECT_NAME}:latest /bin/bash -c "/app/piktoctl sonar -i"

docker_piktoctl_sonar_install_debug: docker_build
	docker run -t --rm -v ${PWD}/bin/:/app ${PROJECT_NAME}:latest /bin/bash -c "/app/piktoctl sonar -i --debug"