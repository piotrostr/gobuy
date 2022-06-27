package buyer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	binance "github.com/adshao/go-binance/v2"
)

var ctx = context.Background()

type Buyer struct {
	Quantity string
	Interval int
	Client   *binance.Client
	Symbol   string
	ch       chan int
}

func Get(quantity string, interval int, symbol string) (*Buyer, error) {
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
		Quantity: quantity,
		Interval: interval,
		Client:   client,
		Symbol:   symbol,
	}
	return buyer, nil
}

func (b *Buyer) Top() error {
	balance, err := b.Client.NewGetAccountService().Do(ctx)
	if err != nil {
		return err
	}
	fmt.Println("Balances:")
	for _, b := range balance.Balances {
		freeBalance, err := strconv.ParseFloat(b.Free, 32)
		if err != nil {
			return err
		}
		if freeBalance != 0 {
			fmt.Printf("%s: %s\n", b.Asset, b.Free)
		}
	}
	fmt.Print("\n")
	return nil
}

func (b *Buyer) Price() (string, error) {
	service := b.Client.NewAveragePriceService()
	price, err := service.Symbol(b.Symbol).Do(ctx)
	if err != nil {
		return "", err
	}
	return price.Price, nil
}

func (b *Buyer) Buy() (*binance.CreateOrderResponse, error) {
	order := b.Client.NewCreateOrderService()
	order.Symbol(b.Symbol).Side("BUY").Type("MARKET")
	res, err := order.Quantity(b.Quantity).Do(ctx)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (b *Buyer) Sell(symbol string, quantity string) (*binance.CreateOrderResponse, error) {
	order := b.Client.NewCreateOrderService()
	order.Symbol(symbol).Side("SELL").Type("MARKET")
	res, err := order.Quantity(b.Quantity).Do(ctx)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (b *Buyer) TestBuy() error {
	order := b.Client.NewCreateOrderService()
	order.Symbol(b.Symbol).Side("BUY").Type("MARKET")
	err := order.Quantity(b.Quantity).Test(ctx)
	if err != nil {
		return err
	}
	return err
}

func (b *Buyer) PrintRes(res any) {
	barr, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(barr))
}

func (b *Buyer) Run() error {
	for {
		res, err := b.Buy()
		if err != nil {
			return err
		}
		b.PrintRes(res)
		time.Sleep(time.Duration(b.Interval) * time.Second)
	}
	// TODO run multiple channels of bots and collect logs from every bot
}
