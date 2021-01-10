package main

import (
	"os"
	"os/exec"
	"crypto/md5"
	"fmt"
	"time"
	"flag"
)

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

func commit_transaction(description string) int {
	home_dir := os.Getenv("HOME")
	os.Chdir(home_dir + "/.expenses/")
	cmd := exec.Command("git", "add", "expenses.txt")
	_, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Print(err)
		return -1
	}

	cmd = exec.Command("git", "commit", "-m", "Tracks : \"" + description + "\"")
	_, err = cmd.CombinedOutput()

	if err != nil {
		fmt.Print(err)
		return -1
	}

	return 0
}

func main() {

	transactions_file := os.Getenv("EXPENSES_FILE") 
	if len(transactions_file) == 0 {
		home_dir := os.Getenv("HOME")
		transactions_file = home_dir + "/.expenses/expenses.txt"
	}

	var addFlag = flag.Bool("add", false, "Create a new transaction")
	var checkFlag = flag.Bool("check", false, "Print a summary of the transactions")
	var valueFlag = flag.Float64("value", 0.0, "Value for the transaction")
	var mFlag = flag.String("m", "", "Description for the transaction")

	flag.Parse()

	if *checkFlag == true {
		cmd := exec.Command("cat", transactions_file)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
	}

	if *addFlag == true {
		var tr = Transaction{value: *valueFlag, description: *mFlag, date: time.Now()}
		tr.id = fmt.Sprintf("%x",  md5.Sum([]byte(tr.date.String())))

		var transaction_line = []byte(transaction_to_line(tr))
		file, err := os.OpenFile(transactions_file, os.O_WRONLY|os.O_APPEND, 0644)

		if err != nil {
			fmt.Println("Error opening " + transactions_file)
			fmt.Println(err)
			os.Exit(-1)
		}

		defer file.Close()
		file.WriteString(string(transaction_line))

		status := commit_transaction(tr.description)

		if status == -1 {
			fmt.Println("Error tracking the new transaction in git")
			os.Exit(-1)
		}
	}
}
