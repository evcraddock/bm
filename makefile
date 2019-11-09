.PHONY: build
build: 
	GOOS=linux GOARCH=amd64 go build -o .build/bm cmd/bm/bm.go

.PHONY: install
install: build
	cp .build/bm ${GOPATH}/bin/bm
