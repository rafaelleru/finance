package main

import (
	"strconv"
	"strings"
	"fmt"
	"time"
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

func build_transaction_from_line(transaction_line string) Transaction {
	dateLayout := "2006-01-02T15:04"
	var tr Transaction
	values := strings.Fields(transaction_line)

	// TODO: Check the if statement
	if values[0] != "" {
		tr.id = values[0]
	}
	if values[1] != "" {
		tr.value, _ = strconv.ParseFloat(values[1], 64)
	}

	if values[2] != "" {
		tr.date, _ = time.Parse(dateLayout, values[2])
	}

	if values[3] != "" {
		tr.description = values[3]
	}

	return tr
}
