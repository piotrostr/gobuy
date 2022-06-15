package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/piotrostr/gobuy/buyer"
)

var ctx = context.Background()

func main() {
	quantity := flag.String("qty", "0.01", "Quantity of ETH")
	interval := flag.Int("interval", 420, "Interval in minutes")
	symbol := flag.String("symbol", "ETHUSDT", "Symbol of the pair to buy")
	run := flag.Bool("run", false, "Run the bot")
	buy := flag.Bool("buy", false, "Buy once")
	docker := flag.Bool("docker", false, "Include the flag if running in container")
	flag.Parse()

	if !*run && !*buy {
		flag.PrintDefaults()
		return
	}

	if !*docker {
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

	if *run {
		msg := "Buying %s ETH every %d minutes\n"
		fmt.Printf(msg, buyer.Quantity, buyer.Interval)
		if err = buyer.Run(); err != nil {
			fmt.Println(err.Error())
		}
	} else if *buy {
		msg := "Buying %s ETH once\n"
		fmt.Printf(msg, buyer.Quantity)
		res, err := buyer.Buy()
		if err != nil {
			fmt.Println(err.Error())
		}
		bytes, err := json.Marshal(res)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(string(bytes))
	}
}
