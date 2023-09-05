# Controtto

An, of course, blazingly fast, self-hosted profit and loss tracker made with Go, HTMX and *no JavaScript*.
You give them two assets (e.g BTC and USDT) and you'll get a nice UI with you Avg. buy price and your transactions, that's it (for now).

You can try it [here](pnl.contre.io)

# Run

## Development env
```bash
# Download air (go install github.com/cosmtrek/air@latest)
go install github.com/cosmtrek/air@latest 
air -c air.toml
```
and access [localhost:8721](http://localhost:8721)

## Run with docker
```bash
# TODO: Create a docker image 
```


# TODO
* Domain validations
* Testing 
* Wrappers for logging and metrics would be nice as well.
