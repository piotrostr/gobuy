# gobuy

Simple go service to keep started in the background and buy you $X of ethereum
during bear market.

Also supports buying from the command line

## Usage

With `.env` file containing `API_KEY` and `SECRET_KEY` from binance

```bash
$ go run ./main.go

  -buy
        Buy once
  -docker
        Include the flag if running in container
  -interval int
        Interval in minutes (default 420)
  -qty string
        Quantity of ETH (default "0.01")
  -run
        Run the bot
  -symbol string
        Symbol of the pair to buy (default "ETHUSDT")

```

Dollar-cost average into ether!
