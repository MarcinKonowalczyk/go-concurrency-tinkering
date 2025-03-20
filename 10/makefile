
.PHONY: all clean

PROTO = ./outliers.proto

all: pb/outliers_pb2.py pb/outliers_pb2_grpc.py pb/outliers.pb.go pb/outliers_grpc.pb.go

pb:
	mkdir -p pb

pb/outliers_pb2.py: $(PROTO) pb
	protoc -I. --python_out=./pb $(PROTO)

pb/outliers_pb2_grpc.py: $(PROTO) pb
	python -m grpc_tools.protoc -I. --grpc_python_out=./pb $(PROTO)

pb/outliers.pb.go: $(PROTO) pb
	protoc --go_out=./pb --go_opt=paths=source_relative $(PROTO)

pb/outliers_grpc.pb.go: $(PROTO) pb
	protoc --go-grpc_out=./pb --go-grpc_opt=paths=source_relative $(PROTO)

clean:
	rm -rf pb
