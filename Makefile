build: finance

finance: finance.go
	go build finance.go

clean:
	rm -rf finance

install: finance
	cp ./finance /usr/bin/finance

