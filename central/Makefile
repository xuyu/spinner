spinner-central: *.go ../*.go
	go build -v -ldflags "-s" -o spinner-central

clean:
	rm -f spinner-central

install:
	mkdir -p /var/lib/spinner-central/bin
	mkdir -p /var/lib/spinner-central/log
	mkdir -p /var/lib/spinner-central/run
	mkdir -p /var/lib/spinner-central/etc/static
	cp -f spinner-central /var/lib/spinner-central/bin
	cp -rf ../webui/* /var/lib/spinner-central/etc/static