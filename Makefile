.PHONY: run
run:
	go run .

.PHONY: build
build: albumin

albumin: rice-box.go
	go build -ldflags "-w -s" -trimpath -o $@

rice-box.go:
	go get github.com/GeertJohan/go.rice/rice
	rice embed-go

.PHONY: clean
clean:
	rm -rf `cat .gitignore`
