build: finance.go
	go build finance.go

clean:
	rm -rf finance

install: expenses
	cp ./finance /usr/bin/finance

