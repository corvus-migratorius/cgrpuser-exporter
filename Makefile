# PREFIX is an environment variable, use the default value if it's not set
ifeq ($(PREFIX),)
    PREFIX := /usr/bin
endif

compile:
	@printf "> "
	go mod tidy
	@printf "> "
	GOOS=linux GOARCH=amd64 go build -v -o cgrpuser-exporter-linux-amd64

install: compile
	@printf "> "
	sudo install -d $(PREFIX)
	@printf "> "
	sudo install -m 770 cgrpuser-exporter-linux-amd64 $(PREFIX)
	@printf "> "
	sudo ln --force --verbose --symbolic $(PREFIX)/cgrpuser-exporter-linux-amd64 $(PREFIX)/cgrpuser-exporter