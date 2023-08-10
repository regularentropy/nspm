main :	make

make:
	go build -a -gcflags=all="-l -B" -ldflags="-w -s"

install: make
	sudo cp -f nspm /usr/local/bin/

clean:
	sudo rm /usr/local/bin/nspm
