package main

import (
	"fmt"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
	agnostypes "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

func main() {
	// 1. Build deps via the standard adapter and inject them into the lib.
	l := agnoslib.New(agnosadapter.New("trackerdata"))

	// 2. Record a few transactions across two categories.
	l.AddCategory("salary")
	l.AddReceived("salary", "august paycheck", 250000)

	l.AddCategory("groceries")
	l.AddSpend("groceries", "weekly shopping", 8450)
	l.AddSpend("groceries", "weekly shopping", 7910)
	mistake, _ := l.AddSpend("groceries", "wrong amount", 99900)

	// 3. A transaction removes itself from the database it was read from.
	if mistake.Remove() {
		fmt.Println("removed:", mistake.Description)
	}

	// 4. ListTransactions walks every category and hands each record back
	//    with its dependencies already wired in.
	fmt.Println("\nall transactions:")
	spent := int64(0)
	for _, transaction := range l.ListTransactions() {
		fmt.Println(" ", transaction.String())
		if transaction.Kind == agnostypes.Spend {
			spent += transaction.Amount
		}
	}

	fmt.Println("\nspent (cents):", spent)
	fmt.Println("balance (cents):", l.Balance())
}
