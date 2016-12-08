get-deps:
	go get github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api
	go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go get golang.org/x/net/context
	go get google.golang.org/grpc
	go get github.com/pkg/errors
	go get github.com/satori/go.uuid

gen-rpc:
	protoc -I/usr/local/include -I. \
	-I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:. \
	./api/feature-toggle.proto

gen-gw:
	protoc -I/usr/local/include -I. \
	-I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--grpc-gateway_out=logtostderr=true:. \
	./api/feature-toggle.proto

gen-swagger:
	protoc -I/usr/local/include -I. \
	 -I${GOPATH}/src \
	 -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	 --swagger_out=logtostderr=true:. \
	 ./api/feature-toggle.proto

build-grpc-server: gen-rpc
	go build -o server/grpc-server \
	./grpc-server/grpc-server.go

build-gw-server: gen-gw
	go build -o server/gw-server \
	./gw-server/gw-server.go

clean:
	rm -rf server/ 
	rm api/*.go 
	rm api/*.swagger.json

build-all: get-deps gen-rpc gen-gw gen-swagger build-grpc-server build-gw-server

