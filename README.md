# gobuy

Simple go service to keep started in the background and buy you $X of ethereum
during bear market.

Also supports buying from the command line.

No plans anytime soon to implement selling (redundant).

## Requires

- `.env` file containing `API_KEY` and `SECRET_KEY` from binance or setting the
  env variables

## Usage

```bash
$ gobuy -price
Balance: 628 USDT
Price: 1119.74556450
```

```bash
$ gobuy -buy [amount-of-eth]
Balance: 628 USDT
Buying 0.01 ETH once
{
  "symbol": "ETHUSDT",
  "orderId": "[secret]",
  "clientOrderId": "[secret]",
  "transactTime": 1655379729428,
  "price": "0.00000000",
  "origQty": "0.01000000",
  "executedQty": "0.01000000",
  "cummulativeQuoteQty": "11.20130000",
  "isIsolated": false,
  "status": "FILLED",
  "timeInForce": "GTC",
  "type": "MARKET",
  "side": "BUY",
  "fills": [
    {
      "price": "1120.13000000",
      "qty": "0.01000000",
      "commission": "0.00001000",
      "commissionAsset": "ETH"
    }
  ],
  "marginBuyBorrowAmount": "",
  "marginBuyBorrowAsset": ""
}
```

### Help message

```bash
$ gobuy

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
