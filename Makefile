build: finance

finance: finance.go movement.go
	go build

clean:
	rm -rf finance

install: finance
	cp ./finance /usr/bin/finance

