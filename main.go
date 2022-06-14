package main

import (
	"flag"
	"fmt"

	"github.com/kardianos/service"
)

type Client struct{}

type Buyer struct {
	quantity float64
	interval int
	client   *Client
}

func (b *Buyer) Start(s service.Service) error {
	fmt.Println("asdf")
	return nil
}

func (b *Buyer) Stop(s service.Service) error {
	fmt.Println("fdsa")
	return nil
}

func main() {
	config := &service.Config{
		Name:        "gobuy",
		DisplayName: "gobuy",
		Description: "Buy X ETH every N time.",
	}
	quantity := flag.Float64("qty", 0.01, "Quantity of ETH to buy")
	interval := flag.Int("interval", 1, "Interval")
	buyer := &Buyer{
		quantity: *quantity,
		interval: *interval,
	}
	s, err := service.New(buyer, config)
	if err != nil {
		println(err)
		return
	}
	msg := "Starting gobuy service, will buy %d ETH every %s"
	println(fmt.Sprintf(msg, buyer.quantity, buyer.interval))
	s.Run()
}
