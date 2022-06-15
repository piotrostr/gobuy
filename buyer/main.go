package buyer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
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
		Quantity: quantity,
		Interval: interval,
		Client:   client,
		Symbol:   symbol,
	}
	return buyer, nil
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

func (b *Buyer) Run() error {
	for {
		res, err := b.Buy()
		if err != nil {
			return err
		}
		fmt.Println(res)
		time.Sleep(time.Duration(b.Interval) * time.Second)
	}

	// put a chan on struct and share it to other methods
	// of a struct? so for example, Stop could have the b.chan <-1
	// possibility to stop the loop
	// kind of no point since here time is blocking, unless
	// select case with time.Timeout and case <-chan
}
