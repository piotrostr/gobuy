package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
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
	balance, err := client.NewGetAccountService().Do(ctx)
	if err != nil {
		return nil, err
	}
	for _, b := range balance.Balances {
		if b.Asset == "USDT" {
			fmt.Printf("Balance: %s USDT\n", strings.Split(b.Free, ".")[0])
		}
	}
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

func (b *Buyer) Start(s service.Service) error {
	ch := make(chan error)
	go b.Run(ch)
	return <-ch
}

func (b *Buyer) Run(ch chan error) {
	for {
		res, err := b.Buy()
		if err != nil {
			ch <- err
			return
		}
		fmt.Println(res)
		time.Sleep(time.Duration(b.interval) * time.Second)
	}
}

func (b *Buyer) Stop(s service.Service) error {
	fmt.Println("Stopping")
	return nil
}

func main() {
	config := &service.Config{
		Name:        "gobuy",
		DisplayName: "gobuy",
		Description: "Buy X ETH every N time.",
	}
	quantity := flag.String("qty", "0.01", "Quantity of ETH to buy")
	interval := flag.Int("interval", 420, "Interval in minutes")
	flag.Parse()

	buyer, err := GetBuyer(*quantity, *interval)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	s, err := service.New(buyer, config)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	logger, err := s.Logger(nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	msg := "Buying %s ETH every %d minutes"
	err = logger.Info(fmt.Sprintf(msg, buyer.quantity, buyer.interval))
	if err != nil {
		fmt.Println(err.Error())
	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
