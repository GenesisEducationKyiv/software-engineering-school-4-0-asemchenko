# Architecture

## Currency rate provider

I choose [Monobank Open API](https://api.monobank.ua/docs/index.html) as a currency rate
provider. That decision has been taken because comparing to other currency rate providers
monobank has some benefits:

1) It's free, no payments, no subscriptions
2) No API key is required

For sure, there are some drawbacks as well:

1) Monobank does not allow to send as many requests as you want. There is a limit

So because of that drawback, I decided to implement a cache mechanism.
The cache is simple, and it stores the currency rate for 1 hour.
So it means that it does not matter how many requests you send to the service,
an actual request to the monobank API will be sent only once per hour. 
