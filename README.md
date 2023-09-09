# 📊 Controtto

A self-hosted, P&L tracker made with Go, HTMX and *no JavaScript*. Controtto, keeps track of your transaction saving it in a sqlite file, and returns all sorts of calculations including:
* Avg. buy price
* Current asset value
* Transaction history
* Profit & Loss 


## Screenshots
See some illustrative screenshorts or just try it it at [pnl.contre.io](https://pnl.contre.io) )
Trading pair | Dashboard
:-------------------------:|:-------------------------:
![accounts-dashboard](./public/assets/img/pairpnl.png) | ![kpi-dashboard](./public/assets/img/pairList.png)

## Configurations

All configurations are set in the `.env` file and passed as environment variables

```sh
# Install the dependencies
go mod tidy
# Set the .env
mv .env.example .env
# Source the env variables
. <(cat .env | grep -v -e '^$' | grep -v "#" | awk '{print "export " $1}')
```

## Build and Run 
```sh
go run ./cmd/main.go # go build ./cmd/main.go to just build it
```

## Development env
```bash
# Download air (go install github.com/cosmtrek/air@latest)
go install github.com/cosmtrek/air@latest 
air -c air.toml # go run ./cmd/main.go
```
and access [localhost:8721](http://localhost:8721)

## Run with docker
```bash
# TODO: Create a docker image 
```

<!-- ### TODO -->
* Testing 
* Wrappers for logging and metrics would be nice as well.
* Remove all the CSS and use custom `style.css` + Tailwind CDN.
