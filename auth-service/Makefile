runp:
	@ go run ./playground/main.go
.PHONY: runp

run:
	@ go run ./main.go
.PHONY: run

proto-gen:
	@ protoc --proto_path=./proto ./proto/*.proto \
		--proto_path=./proto/lib \
		--plugin=$(go env GOPATH)/bin/protoc-gen-go.exe \
		--plugin=$(go env GOPATH)/bin/protoc-gen-grpc-gateway.exe \
		--plugin=$(go env GOPATH)/bin/protoc-gen-openapiv2.exe \
		--plugin=$(go env GOPATH)/bin/protoc-gen-go-grpc.exe \
		--go_out=./internal/handler/rpc/pb --go_opt paths=source_relative \
		--go-grpc_out=./internal/handler/rpc/pb --go-grpc_opt paths=source_relative \
		--grpc-gateway_out ./internal/handler/rpc/pb \
		--grpc-gateway_opt allow_delete_body=true,logtostderr=true,paths=source_relative,repeated_path_param_separator=ssv \
		--openapiv2_out ./docs \
		--openapiv2_opt logtostderr=true,repeated_path_param_separator=ssv
.PHONY: proto-gen

cert:
	@ mkdir ./etc/cert -p
	@ openssl genpkey -algorithm RSA -outform PEM -out ./etc/cert/id_rsa_pkcs8.pem
	@ openssl rsa -pubout -in ./etc/cert/id_rsa.pem -out ./etc/cert/id_rsa.pub
	@ openssl rsa -in id_rsa_pkcs8.pem -out id_rsa.pem
.PHONY: cert