.PHONY: all

all: pb/mudah.pb.go

pb/mudah.pb.go: pb/mudah.proto
	protoc --go_out=plugins=grpc:. pb/mudah.proto