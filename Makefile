
all: vet test demo

test:
	go test ./...

vet:
	go vet ./...

demo:
	make -C demo

.PHONY: demo
