all:
	make build
	make containerize
init:
	go get
build:
	go build -a -tags netgo -ldflags '-extldflags "-w -static"'
containerize: ./chronos
clean: ./chronos
	rm chronos
