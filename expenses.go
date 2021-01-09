package main

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

func main() {
	var addFlag = flag.Bool("add", false, "Create a new transaction")
	var valueFlag = flag.Float64("value", 0.0, "Value for the transaction")
	var mFlag = flag.String("m", "", "Description for the transaction")

	flag.Parse()


	if *addFlag == true {
		var tr = Transaction{value: *valueFlag, description: *mFlag, date: time.Now()}
		tr.id = fmt.Sprintf("%x",  md5.Sum([]byte(tr.description)))

	}


}
