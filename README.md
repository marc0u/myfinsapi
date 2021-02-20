## MyfinsAPI v2.2.0  

## Handle Transactions
POST:/api/myfins/v2/transactions
PUT:/api/myfins/v2/transactions/:id
DELETE:/api/myfins/v2/transactions/:id

## Get Transactions
GET:/api/myfins/v2/transactions?limit=100&order=amount&desc=true _**UPDATED!**_
GET:/api/myfins/v2/transactions/last
GET:/api/myfins/v2/transactions/month?change=-1
GET:/api/myfins/v2/transactions/dates?from=YYYY-MM-DD&to=YYYY-MM-DD
GET:/api/myfins/v2/transactions/summary?change=-1&exclusions=between,transfers _**UPDATED!**_
GET:/api/myfins/v2/transactions/summary/dates?from=YYYY-MM-DD&to=YYYY-MM-DD&exclusions=between,transfers _**UPDATED!**_
GET:/api/myfins/v2/transactions/:id

## Handle Stocks
POST:/api/myfins/v2/stocks
PUT:/api/myfins/v2/stocks/:id
DELETE:/api/myfins/v2/stocks/:id

## Get Stocks
GET:/api/myfins/v2/stocks
GET:/api/myfins/v2/stocks/:id
GET:/api/myfins/v2/stocks/holdings 
GET:/api/myfins/v2/stocks/portfolio/daily _**NEW!**_
GET:/api/myfins/v2/stocks/portfolio/daily?detailed=true` _**NEW!**_