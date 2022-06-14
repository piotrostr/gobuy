package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	binance "github.com/adshao/go-binance/v2"
	"github.com/joho/godotenv"
	"github.com/kardianos/service"
)

var ctx = context.Background()

type Buyer struct {
	quantity string
	interval int
	client   *binance.Client
}

func GetBuyer(quantity string, interval int) (*Buyer, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return nil, errors.New("API_KEY not set")
	}
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return nil, errors.New("SECRET_KEY not set")
	}
	client := binance.NewClient(apiKey, secretKey)
	buyer := &Buyer{
		quantity: quantity,
		interval: interval,
		client:   client,
	}
	return buyer, nil
}

func (b *Buyer) Buy() (*binance.CreateOrderResponse, error) {
	order := b.client.NewCreateOrderService()
	order.Symbol("ETHUSDT").Side("BUY").Type("MARKET")
	res, err := order.Quantity(b.quantity).Do(ctx)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (b *Buyer) TestBuy() error {
	order := b.client.NewCreateOrderService()
	order.Symbol("ETHUSDT").Side("BUY").Type("MARKET")
	err := order.Quantity(b.quantity).Test(ctx)
	if err != nil {
		return err
	}
	return err
}

func (b *Buyer) Loop() error {
	for {
		err := b.TestBuy()
		if err != nil {
			return err
		}
		println("Bought (Test)")
		time.Sleep(time.Duration(b.interval) * time.Second)
	}
}

func (b *Buyer) Start(s service.Service) error {
	err := b.Loop()
	if err != nil {
		return err
	}
	return nil
}

func (b *Buyer) Stop(s service.Service) error {
	println("shutting down")
	return nil
}

func main() {
	config := &service.Config{
		Name:        "gobuy",
		DisplayName: "gobuy",
		Description: "Buy X ETH every N time.",
	}
	quantity := flag.String("qty", "0.001", "Quantity of ETH to buy")
	interval := flag.Int("interval", 1, "Interval in minutes")
	flag.Parse()
	buyer, err := GetBuyer(*quantity, *interval)
	if err != nil {
		println(err.Error())
		return
	}
	s, err := service.New(buyer, config)
	if err != nil {
		println(err.Error())
		return
	}
	msg := "Starting gobuy service, will buy %s ETH every %d minutes"
	println(fmt.Sprintf(msg, buyer.quantity, buyer.interval))
	err = s.Run()
	if err != nil {
		println(err.Error())
	}
}
