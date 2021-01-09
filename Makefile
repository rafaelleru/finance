build: expenses.go
	go build expenses.go

clean:
	rm -rf expenses

install: expenses
	cp ./expenses /usr/bin/expenses

