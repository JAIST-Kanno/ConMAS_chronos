all:
	make init
	make build
	make containerize

init:
	go get

build:
	go build -a -tags netgo -ldflags '-extldflags "-w -static"'

containerize: ./ConMAS_chronos
	docker build -t docker.pkg.github.com/jaist-kanno/conmas_chronos/chronos:1.0 .

clean: ./ConMAS_chronos
	rm ConMAS_chronos
