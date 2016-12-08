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
	perl -i -0pe \
	's/.*filter_FeatureToggleService_GetFeaturesForProperties_0.*\n.*\n.*}.*/protoReq.Properties = make(map[string]string)\nfor k, v := range req.URL.Query() {\nprotoReq.Properties[k] = v[0]\n}/' \
	api/feature-toggle.pb.gw.go

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


build-all: gen-rpc gen-gw gen-swagger build-grpc-server build-gw-server

