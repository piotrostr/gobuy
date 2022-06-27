package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/piotrostr/gobuy/buyer"
)

var ctx = context.Background()

func envExists() bool {
	exists := false

	files, err := os.ReadDir("./")
	if err != nil {
		fmt.Println(err)
		return exists
	}

	for _, f := range files {
		if f.Name() == ".env" {
			exists = true
		}
	}
	return exists
}

/*
* this is all kind of wonky so would replace it with cobra
* nice completion and help output and can all args in a type-safe way
 */
func main() {
	quantity := flag.String("qty", "0.01", "Quantity of ETH")
	interval := flag.Int("interval", 420, "Interval in minutes")
	symbol := flag.String("symbol", "ETHUSDT", "Symbol of the pair to buy")
	run := flag.Bool("run", false, "Run the bot")
	buy := flag.Bool("buy", false, "Buy once")
	docker := flag.Bool("docker", false, "Include the flag if running in container")
	flag.Parse()

	if !*docker && envExists() {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file")
			return
		}
	}

	buyer, err := buyer.Get(*quantity, *interval, *symbol)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if !*run && !*buy {
		err = buyer.Top()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		flag.PrintDefaults()
		return
	}

	err = buyer.Top()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if *run {
		msg := "Buying %s ETH every %d minutes\n"
		fmt.Printf(msg, buyer.Quantity, buyer.Interval)
		if err = buyer.Run(); err != nil {
			fmt.Println(err.Error())
		}
	} else if *buy {
		fmt.Printf("Buying %s ETH once\n", buyer.Quantity)
		res, err := buyer.Buy()
		if err != nil {
			fmt.Println(err.Error())
		}
		b, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(string(b))
	}
}
