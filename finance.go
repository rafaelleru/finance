package main

import (
	"strconv"
	"bufio"
	"os"
	"os/exec"
	"crypto/md5"
	"fmt"
	"time"
	"flag"
)

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

func get_balance(dateStart string, dateEnd string) (float64, []string) {
	value := 0.0
	var expensesInDateRange []string
	layout := "2006-01-02"

	home_dir := os.Getenv("HOME")
	os.Chdir(home_dir + "/.expenses/")

	tStart, err := time.Parse(layout, dateStart)

	if err != nil {
		fmt.Println("Error parsing date start")
		fmt.Println(err)
		os.Exit(-1)
	}

	tEnd, err := time.Parse(layout, dateEnd)

	if err != nil {
		fmt.Println("Error parsing date start")
		os.Exit(-1)
	}

	//open transactions file 
	transactions, err := os.Open("expenses.txt")
	if err != nil {
		fmt.Println("Could not open transactions file")
		os.Exit(-1)
	}

	// TODO: Improve how we are filtering the lines
	scanner := bufio.NewScanner(transactions)
	for scanner.Scan() {
		line := scanner.Text()
		tr := build_transaction_from_line(line)

		if tStart.Before(tr.date) && tEnd.After(tr.date) {
			value += tr.value
			expensesInDateRange = append(expensesInDateRange, line)
		}
	}

	return value, expensesInDateRange
}

func main() {

	const colorRed = "\033[31m"
	const colorGreen = "\033[32m"

	transactions_file := os.Getenv("EXPENSES_FILE") 
	if len(transactions_file) == 0 {
		home_dir := os.Getenv("HOME")
		transactions_file = home_dir + "/.expenses/expenses.txt"
	}

	var addFlag = flag.Bool("add", false, "Create a new transaction")
	var checkFlag = flag.Bool("check", true, "Print a summary of the transactions")
	var balanceFlag = flag.Bool("balance", false, "Get a balance between dates")
	var findFlag = flag.String("find", "", "Find matching money movements (Using grep)")
	var nLinesFlag = flag.Int("n", 50, "Number of transactions to show in resume")
	var valueFlag = flag.Float64("value", 0.0, "Value for the transaction")
	var mFlag = flag.String("m", "", "Description for the transaction")
	var dateStart = flag.String("start", time.Now().String(), "Get a balance between dates")
	var dateEnd = flag.String("end", time.Now().String(), "Get a balance between dates")

	flag.Parse()

	if *findFlag != "" {
		cmd := exec.Command("rg", "-i", "-e", *findFlag, transactions_file)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
		os.Exit(0)
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

		os.Exit(0)
	}

	if *balanceFlag == true {
		var color string

		if dateStart == nil || dateEnd == nil {
			fmt.Println("start and end parameters must be provided to get a balance")
			os.Exit(-1)
		}

		total_balance, expensesInBetweenRange := get_balance(*dateStart, *dateEnd)

		for _, ex := range expensesInBetweenRange {
			fmt.Println(ex)
		}

		if total_balance > 0 {
			color = colorGreen
		} else {
			color = colorRed
		}

		fmt.Printf("Total balance: %s%.2f\n", color, total_balance)
		os.Exit(0)
	}

	if *checkFlag == true {
		nLines := strconv.Itoa(*nLinesFlag)
		cmd := exec.Command("tail",  "-n", nLines, transactions_file)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
		os.Exit(0)
	}

}
