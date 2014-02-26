spinner-node: *.go ../sensor/*.go ../*.go
	go build -v -ldflags "-s" -o spinner-node

clean:
	rm -f spinner-node

install:
	mkdir -p /var/lib/spinner-node/bin
	mkdir -p /var/lib/spinner-node/log
	mkdir -p /var/lib/spinner-node/run
	mkdir -p /var/lib/spinner-node/etc
	cp -f spinner-node.conf /etc/init
	cp -f spinner-node /var/lib/spinner-node/bin
