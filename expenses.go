package main

import "os"
import "os/exec"
import "crypto/md5"
import "fmt"
import "time"
import "flag"

type Transaction struct {
	// TODO: use md5 as type??
	id string
	value float64
	date time.Time
	description string
}

func transaction_to_line(tr Transaction) string {
	var date_string = fmt.Sprintf("%d-%02d-%02dT%02d:%02d",
									tr.date.Year(), tr.date.Month(), tr.date.Day(),
									tr.date.Hour(), tr.date.Minute())
	return fmt.Sprintf("%s\t%.2f\t%s\t%s\n", tr.id, tr.value, date_string, tr.description)
}

func main() {

	transactions_file := os.Getenv("EXPENSES_FILE") 
	if len(transactions_file) == 0 {
		transactions_file = "~/.expenses/expenses.txt"
	}

	var addFlag = flag.Bool("add", false, "Create a new transaction")
	var checkFlag = flag.Bool("check", false, "Print a summary of the transactions")
	var valueFlag = flag.Float64("value", 0.0, "Value for the transaction")
	var mFlag = flag.String("m", "", "Description for the transaction")

	flag.Parse()

	if *checkFlag == true {
		cmd := exec.Command("cat", "transactions")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
	}

	if *addFlag == true {
		fmt.Println("Adding new transaction")
		var tr = Transaction{value: *valueFlag, description: *mFlag, date: time.Now()}
		tr.id = fmt.Sprintf("%x",  md5.Sum([]byte(tr.date.String())))

		var transaction_line = []byte(transaction_to_line(tr))
		file, err := os.OpenFile(transactions_file, os.O_WRONLY|os.O_APPEND, 0644)

		if err != nil {
			fmt.Println("Error writing to file")
		}

		defer file.Close()
		file.WriteString(string(transaction_line))
	}
}
