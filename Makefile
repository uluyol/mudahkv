
BINS = mudahd mudahc sendout
REPO = quay.io/uluyol/mudah
VERSION = 0.5

.PHONY: all docker-build

all: lib/pb/mudah.pb.go

lib/pb/mudah.pb.go: lib/pb/mudah.proto
	cd lib/pb && protoc --go_out=plugins=grpc:. mudah.proto

sendout mudahc mudahd: cmd
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build github.com/uluyol/mudahkv/cmd/$@

docker-build: lib/pb/mudah.pb.go $(BINS)
	docker build -t $(REPO):$(VERSION) .

docker-push: docker-build
	docker push $(REPO):$(VERSION)

clean:
	rm -rf $(BINS)
