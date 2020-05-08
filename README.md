## Stonks-Service (example app)


### About
stonk-service is an example grpc application that allows CRUD operations on stock data, grabbing the latest pricing information on stocks from yahoo finance API.


### local setup

1. create a RapidAPI account to gain access to yahoo finance's API (https://rapidapi.com/apidojo/api/yahoo-finance1)
2. update the `TRADING_API_KEY` env variable for the stonks-service in docker-compose
3. run the `make run` command
4. use the script in `cmd/consumer/main.go` to make requests to your locally running service