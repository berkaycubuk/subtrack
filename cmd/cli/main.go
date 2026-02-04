package main

import (
	"fmt"
	"log"
	"os"

	"github.com/berkaycubuk/subtrack/internal/cli"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	c, err := cli.New()
	if err != nil {
		log.Fatalf("Failed to initialize CLI: %v", err)
	}

	switch command {
	case "add":
		if len(os.Args) < 7 {
			fmt.Println("Usage: subtrack add <name> <price> <currency> <cycle> <payment_date>")
			fmt.Println("Example: subtrack add \"Netflix\" 15.99 USD monthly 15-02-2025")
			os.Exit(1)
		}
		name := os.Args[2]
		price := os.Args[3]
		currency := os.Args[4]
		cycle := os.Args[5]
		paymentDate := os.Args[6]
		if err := c.Add(name, price, currency, cycle, paymentDate); err != nil {
			log.Fatalf("Error: %v", err)
		}

	case "list":
		if err := c.List(); err != nil {
			log.Fatalf("Error: %v", err)
		}

	case "update":
		if len(os.Args) < 3 {
			fmt.Println("Usage: subtrack update <id> [name] [price] [currency] [cycle] [payment_date]")
			fmt.Println("Example: subtrack update 1 \"Netflix\" 19.99 USD monthly 15-03-2025")
			os.Exit(1)
		}
		id := os.Args[2]
		var name, price, currency, cycle, paymentDate string
		if len(os.Args) > 3 {
			name = os.Args[3]
		}
		if len(os.Args) > 4 {
			price = os.Args[4]
		}
		if len(os.Args) > 5 {
			currency = os.Args[5]
		}
		if len(os.Args) > 6 {
			cycle = os.Args[6]
		}
		if len(os.Args) > 7 {
			paymentDate = os.Args[7]
		}
		if err := c.Update(id, name, price, currency, cycle, paymentDate); err != nil {
			log.Fatalf("Error: %v", err)
		}

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: subtrack delete <id>")
			fmt.Println("Example: subtrack delete 1")
			os.Exit(1)
		}
		id := os.Args[2]
		if err := c.Delete(id); err != nil {
			log.Fatalf("Error: %v", err)
		}

	case "check":
		if err := c.Check(); err != nil {
			log.Fatalf("Error: %v", err)
		}

	case "health":
		if err := c.Health(); err != nil {
			log.Fatalf("Error: %v", err)
		}

	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("SubTrack CLI - Subscription Tracker")
	fmt.Println("\nUsage:")
	fmt.Println("  subtrack add <name> <price> <currency> <cycle> <payment_date>")
	fmt.Println("  subtrack list")
	fmt.Println("  subtrack update <id> [name] [price] [currency] [cycle] [payment_date]")
	fmt.Println("  subtrack delete <id>")
	fmt.Println("  subtrack check")
	fmt.Println("  subtrack health")
	fmt.Println("\nExamples:")
	fmt.Println("  subtrack add \"Netflix\" 15.99 USD monthly 15-02-2025")
	fmt.Println("  subtrack list")
	fmt.Println("  subtrack update 1 \"Netflix\" 19.99 USD monthly 15-03-2025")
	fmt.Println("  subtrack delete 1")
	fmt.Println("  subtrack check")
	fmt.Println("  subtrack health")
}
