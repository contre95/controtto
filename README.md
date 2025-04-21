# 📊 Controtto

A self-hosted, P&L tracker made with Go, HTMX and *no JavaScript*. Controtto, keeps track of your transaction saving it in a sqlite file, and returns all sorts of calculations including:
* Avg. Buy price
* Current asset value
* Transaction history
* Import / Export transaction (see [sample](./public/assets/export_USDT_ARS.csv))
* Profit & Loss 

In order to fetch the price of an asset, Controtto relies on mainly 4 APIs ([see code](https://github.com/contre95/controtto/tree/main/src/gateways/markets)). 
* [Binance](https://api.binance.com/api/v3/ticker/price) - Public API, no token needed.
* [BingX](https://open-api.bingx.com/openApi/swap/v2/quote/price) - Public API, no token needed.
* [Alpha Vantage](https://www.alphavantage.co/) - Stocks, free but short rate limit. ([get](https://www.alphavantage.co/support/#api-key) an token and set `CONTROTTO_AVANTAGE_TOKEN`)
* [Tiingo](https://www.tiingo.com/documentation/) - Stocks, crypto and Forex. ([create](https://www.tiingo.com/) account/token and set `CONTROTTO_TIINGO_TOKEN`) 

## Demo
You can also check out the demo at [demo.contre.io](https://demo.contre.io), the database resets every hour.


https://github.com/user-attachments/assets/28bf8782-6118-47c0-a94b-085e1a7045b2




All configurations are set in the `.env` file and passed as environment variables. Variables `CONTROTTO_PORT` and `CONTROTTO_DB_PATH` are available.
```sh
# Install the dependencies
go mod tidy
# Set the .env
mv .env.example .env
# Source the env variables
. <(cat .env | grep -v -e '^$' | grep -v "#" | awk '{}')
```

## Build and Run 
```sh
go run ./cmd/main.go # go build ./cmd/main.go to just build it
```

## Development env
```sh
go install github.com/cosmtrek/air@latest # Download air
air -c air.toml
```
and access [localhost:3000](http://localhost:3000)

## Run with Podman
A [Container image](https://hub.docker.com/r/contre95/controtto) is available on Docker's public registry.
If you want to use Docker, simply replace `podman` with `docker`. 

```sh
mkdir data
podman container run --rm -p 8000:8000 -v $(pwd)/data:/data contre95/controtto
```

## Run tests
```sh
go test -cover ./...
#   Expected result
#   ?       controtto/cmd   [no test files]
#   ?       controtto/src/app/managing      [no test files]
#   ?       controtto/src/domain/pnl        [no test files]
#   ?       controtto/src/gateways/markets  [no test files]
#   ?       controtto/src/gateways/sqlite   [no test files]
#   ?       controtto/src/presenters        [no test files]
#   ok      controtto/src/app/querying      0.003s  coverage: 40.7% of statements
```
### TODO
* More tests
* Wrappers for logging and metrics would be nice as well.
* Add Accounts to keep track of the total net worth.
* Remove all the CSS and use custom `style.css` + Tailwind CDN.
