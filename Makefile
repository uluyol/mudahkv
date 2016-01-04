
BINS = mudahd mudahc
REPO = uluyol/mudah
VERSION = 0.1

.PHONY: all docker-build

all: pb/mudah.pb.go

pb/mudah.pb.go: pb/mudah.proto
	protoc --go_out=plugins=grpc:. pb/mudah.proto

$(BINS):
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build github.com/uluyol/mudahkv/cmd/$@

docker-build: pb/mudah.pb.go $(BINS)
	docker build -t $(REPO):$(VERSION) .

docker-push: docker-build
	docker push $(REPO):$(VERSION)

clean:
	rm -f $(BINS)