# MyfinsAPI v2.2.0  

## Endpoints

### Auth
POST:/login<br />
GET:/notify?name=XXXX&email=XXXX<br />

### Handle Transactions
POST:/api/myfins/v2/transactions<br />
PUT:/api/myfins/v2/transactions/:id<br />
DELETE:/api/myfins/v2/transactions/:id<br />

### Get Transactions
GET:/api/myfins/v2/transactions?limit=100&order=amount&desc=true _**UPDATED!**_<br />
GET:/api/myfins/v2/transactions/last<br />
GET:/api/myfins/v2/transactions/month?change=-1<br />
GET:/api/myfins/v2/transactions/dates?from=YYYY-MM-DD&to=YYYY-MM-DD<br />
GET:/api/myfins/v2/transactions/summary?change=-1&exclusions=between,transfers _**UPDATED!**_<br />
GET:/api/myfins/v2/transactions/summary/dates?from=YYYY-MM-DD&to=YYYY-MM-DD&exclusions=between,transfers _**UPDATED!**_<br />
GET:/api/myfins/v2/transactions/:id<br />

### Handle Stocks
POST:/api/myfins/v2/stocks<br />
PUT:/api/myfins/v2/stocks/:id<br />
DELETE:/api/myfins/v2/stocks/:id<br />

### Get Stocks
GET:/api/myfins/v2/stocks<br />
GET:/api/myfins/v2/stocks/:id<br />
GET:/api/myfins/v2/stocks/holdings<br /> 
GET:/api/myfins/v2/stocks/portfolio/daily _**NEW!**_<br />
GET:/api/myfins/v2/stocks/portfolio/daily?detailed=true` _**NEW!**_<br />

## Evironment Variables

From .env file in the root of the project folder.

### API

| Variable | Value | Description |
| ---------| ----- | ----------- |
| API_PORT | :8080 | Api port |
| API_JWT_SECRET | aaabbbccc | Secret to sign JWT |
| API_GOOGLE_CLIENT | 111-aaa.apps.googleusercontent.com | Google client ID (Google Cloud) |
| TB_ID | aaabbbccc | Telegram Bot Token |

### Database Postgres 

| Variable | Value | Description |
| ---------| ----- | ----------- |
| DB_HOST | 127.0.0.1 | Database host |
| DB_PORT | 5432 | Database port |
| DB_DRIVER | postgres | Database driver name | 
| DB_USER | steven | Database User name |
| DB_PASSWORD | Database Password |
| DB_NAME | fullstack_api | Database name |
| DB_DEBUG | false | Debug mode | 

### Database Mysql 

| Variable | Value | Description |
| ---------| ----- | ----------- |
| DB_HOST | 127.0.0.1 | Database host |
| DB_PORT | 3306 | Database port |
| DB_DRIVER | mysql | Database driver name | 
| DB_USER | steven | Database User name |
| DB_PASSWORD | mypass | Database Password |
| DB_NAME | fullstack_mysql | Database name |
| DB_DEBUG | false | Debug mode |

### Database Sqlite 

| Variable | Value | Description |
| ---------| ----- | ----------- |
| DB_DRIVER | sqlite3 | Database driver name |
| DB_NAME | fullstack.sqlite | Database name |
| DB_DEBUG | true | Debug mode |