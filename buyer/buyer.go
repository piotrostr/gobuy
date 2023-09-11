package buyer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
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
	priceService := b.Client.NewAveragePriceService()
	balances := make(map[string]float64)
	for _, balance := range balance.Balances {
		freeBalance, err := strconv.ParseFloat(balance.Free, 32)
		if err != nil {
			return err
		}
		if freeBalance != 0 {
			if balance.Asset == "ETHW" {
				continue
			}
			if balance.Asset == "USDT" {
				balances[balance.Asset] = freeBalance
				continue
			}
			price, err := priceService.Symbol(balance.Asset + "USDT").Do(ctx)
			if err != nil {
				return err
			}
			// parse string to float
			priceFloat, err := strconv.ParseFloat(price.Price, 32)
			if err != nil {
				return err
			}
			balances[balance.Asset] = freeBalance * priceFloat
		}
	}
	fmt.Print("\n")
	price, err := b.Price()
	if err != nil {
		return err
	}
	fmt.Printf("Price (%s): %s\n\n", b.Symbol, price)

	// sum balances
	var sum float64
	for _, balance := range balances {
		sum += balance
	}

	// print total
	fmt.Printf("Total Balance: %f USDT\n", sum)

	// sort by highest percentage
	type kv struct {
		Key   string
		Value float64
	}
	var ss []kv
	for k, v := range balances {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	// print sorted balances as a table
	fmt.Print("\n")
	for _, kv := range ss {
		percent := (kv.Value / sum) * 100
		fmt.Printf("%s: %f USDT (%f%%)\n", kv.Key, kv.Value, percent)
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
