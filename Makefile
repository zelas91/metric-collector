agent-build:
	go build -buildvcs=false  -o agent  -ldflags "-X main.buildCommit=$$(git rev-parse --short HEAD)\
 	-X main.buildDate=$$(date +'%Y-%m-%d_%H:%M')\
 	-X main.buildVersion="1.0"" cmd/agent/*.go
server-build:
	go build -buildvcs=false  -o server  -ldflags "-X main.buildCommit=$$(git rev-parse --short HEAD)\
	 -X main.buildDate=$$(date +'%Y-%m-%d_%H:%M')\
	  -X main.buildVersion="1.0"" cmd/server/*.go


.PHONY: build
BINARY_NAME_SERVER = server-metrics
BINARY_NAME_AGENT = agent-metrics
CMD_SERVER=cmd/server/*.go
CMD_AGENT=cmd/agent/*.go
ARCH = amd64 arm64 arm 386
PLATFORMS = linux darwin windows
build-server:
	go build -o build/server/$(BINARY_NAME_SERVER)-$(GOOS)-$(GOARCH) -ldflags "-X main.buildCommit=$$(git rev-parse --short HEAD)\
                                                       	 -X main.buildDate=$$(date +'%Y-%m-%d_%H:%M')\
                                                       	  -X main.buildVersion="1.0"" $(CMD_SERVER)
build-agent:
	go build -o build/agent/$(BINARY_NAME_AGENT)-$(GOOS)-$(GOARCH) -ldflags "-X main.buildCommit=$$(git rev-parse --short HEAD)\
                                                        -X main.buildDate=$$(date +'%Y-%m-%d_%H:%M')\
                                                        -X main.buildVersion="1.0"" $(CMD_AGENT)

clean:
	 if [ -d "build" ]; then \
            rm -r build; \
        fi
build-all: clean
	$(foreach GOOS,$(PLATFORMS),\
		$(foreach GOARCH,$(ARCH),\
			GOOS=$(GOOS) GOARCH=$(GOARCH) make build-server build-agent;))
generate:
	protoc --go_out=./api/gen --go_opt=paths=source_relative --go-grpc_out=./api/gen --go-grpc_opt=paths=source_relative api/proto/api.proto