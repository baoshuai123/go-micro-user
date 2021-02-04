
GOPATH:=$(shell go env GOPATH)
.PHONY: proto
proto:
	docker run --rm -v d:/GOLANG/src/user:/d/GOLANG/src/user -w /d/GOLANG/src/user  -e ICODE=2606C833CD172F4C cap1573/cap-protoc -I ./ --go_out=./ --micro_out=./ ./proto/user/user.proto
	
.PHONY: build
build:
	go build -o user

.PHONY: docker
docker:
	docker build . -t user:latest
