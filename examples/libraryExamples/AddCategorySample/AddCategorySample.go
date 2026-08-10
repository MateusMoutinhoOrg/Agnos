package main

import (
	"fmt"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
	// 1. Build deps via an adapter (real clock + a Keep database on disk).
	deps := agnosadapter.New("trackerdata", ".")

	// 2. Inject deps into the pure library — a financial tracker.
	l := agnoslib.New(deps)

	// 3. Create the categories transactions are tracked under. Creating a
	//    category is idempotent, so re-running the sample is safe.
	for _, name := range []string{"groceries", "salary", "rent"} {
		category, ok := l.AddCategory(name)
		if !ok {
			fmt.Println("could not create category:", name)
			continue
		}
		fmt.Println("category ready:", category.Name)
	}

	// 4. Read them back out of the database. The library never knows which
	//    adapter is behind it.
	fmt.Println("\nstored categories:")
	for _, category := range l.ListCategories() {
		fmt.Println(" ", category.String())
	}

	// 5. A category that was never created is a miss.
	if _, ok := l.GetCategory("holidays"); !ok {
		fmt.Println("\nholidays: no such category")
	}
}
