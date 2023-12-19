agent-build:
	go build -buildvcs=false  -o agent  -ldflags "-X main.buildCommit=$$(git rev-parse --short HEAD)\
 	-X main.buildDate=$$(date +'%Y-%m-%d_%H:%M')\
 	-X main.buildVersion="1.0"" cmd/agent/*.go
server-build:
	go build -buildvcs=false  -o server  -ldflags "-X main.buildCommit=$$(git rev-parse --short HEAD)\
	 -X main.buildDate=$$(date +'%Y-%m-%d_%H:%M')\
	  -X main.buildVersion="1.0"" cmd/server/*.go