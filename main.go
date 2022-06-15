package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/piotrostr/gobuy/buyer"
)

var ctx = context.Background()

func main() {
	quantity := flag.String("qty", "0.01", "Quantity of ETH to buy")
	interval := flag.Int("interval", 420, "Interval in minutes")
	flag.Parse()

	buyer, err := buyer.Get(*quantity, *interval)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	msg := "Buying %s ETH every %d minutes\n"
	fmt.Printf(msg, buyer.Quantity, buyer.Interval)
	err = buyer.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}
