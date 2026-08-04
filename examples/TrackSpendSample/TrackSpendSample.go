package main

import (
	"fmt"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
	// 1. Build deps via the standard adapter (real clock + a Keep database
	//    on disk under "trackerdata").
	deps := agnosadapter.New("trackerdata")

	// 2. Inject deps into the pure library — a financial tracker.
	l := agnoslib.New(deps)

	// 3. Money coming in and money going out, both under one category.
	//    Amounts are in the smallest currency unit, so 250000 is 2500.00.
	l.AddCategory("salary")
	l.AddReceived("salary", "august paycheck", 250000)

	groceries, _ := l.AddCategory("groceries")
	groceries.AddSpend("weekly shopping", 8450)
	groceries.AddSpend("coffee beans", 1290)

	// 4. A transaction can also be recorded straight from the lib, as long
	//    as the category already exists.
	if _, ok := l.AddSpend("holidays", "flight", 42000); !ok {
		fmt.Println("holidays: no such category, nothing recorded")
	}

	// 5. Each category keeps its own running balance.
	fmt.Println("\nper category:")
	for _, category := range l.ListCategories() {
		fmt.Println(" ", category.String())
	}

	// 6. And the tracker sums every category into one balance: received
	//    money counts up, spent money counts down.
	fmt.Println("\ngroceries balance (cents):", groceries.Balance())
	fmt.Println("total balance (cents):", l.Balance())
}
