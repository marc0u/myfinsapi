package handlers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/marc0u/myfinsapi/api/models"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

func (server *Server) Help(c *fiber.Ctx) error {
	version := apiVersion[0:1]
	var msg = `MyfinsAPI v%v

Helpers
GET:/api/help
GET:/api/stack

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
GET:/api/myfins/v%[1]v/stocks/summary
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

func (server *Server) MirrorProductionTables() error {
	trans := []models.Transaction{}
	stocks := []models.Stock{}
	urlTrans := "http://192.168.1.15:7001/api/myfins/v2/transactions"
	urlStocks := "http://192.168.1.15:7001/api/myfins/v2/stocks"
	client := resty.New().SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).SetRetryCount(3).SetRetryWaitTime(3 * time.Second).SetRetryMaxWaitTime(5 * time.Second).SetTimeout(10 * time.Second)
	client.SetAuthToken(os.Getenv("API_CLIENT_TOKEN"))
	resp, err := client.R().Get(urlTrans)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp.Body(), &trans)
	if err != nil {
		return err
	}
	resp, err = client.R().Get(urlStocks)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp.Body(), &stocks)
	if err != nil {
		return err
	}
	for i, j := 0, len(trans)-1; i < j; i, j = i+1, j-1 {
		trans[i], trans[j] = trans[j], trans[i]
	}
	for _, item := range trans {
		// Saving data
		_, err := item.SaveTransaction(server.DB)
		if err != nil {
			return err
		}
	}
	for i, j := 0, len(stocks)-1; i < j; i, j = i+1, j-1 {
		stocks[i], stocks[j] = stocks[j], stocks[i]
	}
	for _, item := range stocks {
		// Saving data
		_, err := item.SaveAStock(server.DB)
		if err != nil {
			return err
		}
	}
	return nil
}
