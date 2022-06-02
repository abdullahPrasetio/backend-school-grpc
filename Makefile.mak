create:
	protoc --proto_path=proto proto/*.proto  --go_out=common/models
	protoc --proto_path=proto proto/*.proto  --go-grpc_out=common/models
	protoc -I . --grpc-gateway_out ./common/models/ \
    --grpc-gateway_opt logtostderr=true \
    --grpc-gateway_opt paths=source_relative \
    --proto_path=proto proto/*.proto

clean:
	rm common/models/proto/*.go