.PHONY: all client server debuild

all: client server

client:
	cd ./mosiuterclient && \
	go build -o mosiuterclient
	@echo "client built in mosiuterclient folder"

server:
	cd ./mosiuterserver && \
	go build
	@echo "server built in mosiuterserver folder"

debuild:
	rm -rf mosiuterclient/mosiuterclient mosiuterserver/mosiuterserver
	@echo "Deleted build files successfully"
