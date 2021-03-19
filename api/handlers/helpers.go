package handlers

import (
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

func (server *Server) Help(c *fiber.Ctx) error {
	version := apiVersion[0:1]
	var msg = `MyfinsAPI v%v
	
Auth
POST:/login
GET:/notify?name=XXXX&email=XXXX

Handle Transactions
POST:/api/myfins/v%[1]v/transactions
PUT:/api/myfins/v%[1]v/transactions/:id
DELETE:/api/myfins/v%[1]v/transactions/:id

Get Transactions
GET:/api/myfins/v%[1]v/transactions?limit=100&order=amount&desc=true
GET:/api/myfins/v%[1]v/transactions/last
GET:/api/myfins/v%[1]v/transactions/month?change=-1
GET:/api/myfins/v%[1]v/transactions/dates?from=YYYY-MM-DD&to=YYYY-MM-DD
GET:/api/myfins/v%[1]v/transactions/summary?change=-1&exclusions=between,transfers
GET:/api/myfins/v%[1]v/transactions/summary/dates?from=YYYY-MM-DD&to=YYYY-MM-DD&exclusions=between,transfers
GET:/api/myfins/v%[1]v/transactions/:id

Handle Stocks
POST:/api/myfins/v%[1]v/stocks
PUT:/api/myfins/v%[1]v/stocks/:id
DELETE:/api/myfins/v%[1]v/stocks/:id

Get Stocks
GET:/api/myfins/v%[1]v/stocks
GET:/api/myfins/v%[1]v/stocks/:id
GET:/api/myfins/v%[1]v/stocks/holdings
GET:/api/myfins/v%[1]v/stocks/portfolio/daily
GET:/api/myfins/v%[1]v/stocks/portfolio/daily?detailed=true
`

	msg = fmt.Sprintf(msg, version)
	return c.SendString(msg)
}

func (server *Server) Notify(c *fiber.Ctx) error {
	if c.Query("name") == "" || c.Query("email") == "" {
		return nil
	}
	msg := fmt.Sprintf("%s (%s) is trying to sign-in on Myfins-Web", c.Query("name"), c.Query("email"))
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?text=%s&chat_id=%s", os.Getenv("TB_ID"), msg, "165270556")
	client := resty.New()
	client.
		SetRetryCount(3).
		SetRetryWaitTime(3 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second).
		SetTimeout(10 * time.Second).
		R().
		Get(url)
	return nil
}
