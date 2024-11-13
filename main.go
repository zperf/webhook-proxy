package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

var listenAddr string
var httpProxy string
var headerForwardPrefix string

func init() {
	flag.StringVar(&listenAddr, "listen", ":4000", "Listen address")
	flag.StringVar(&httpProxy, "http-proxy", "http://127.0.0.1:8899", "HTTP proxy address")
	flag.StringVar(&headerForwardPrefix, "header-prefix", "X-Github-", "Header forward prefix")
}

func main() {
	flag.Parse()
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("WebHook Proxy")
	})

	app.Post("/*", func(c *fiber.Ctx) error {
		headers := make(map[string]string)
		for key, header := range c.GetReqHeaders() {
			if headerForwardPrefix == "" || strings.HasPrefix(key, headerForwardPrefix) {
				headers[key] = strings.Join(header, ",")
			}
		}

		dst := fmt.Sprintf("https://%s", c.Params("*"))
		log.Printf("Posting to `%s`", dst)
		rsp, err := resty.New().SetProxy(httpProxy).R().
			SetHeaders(headers).
			SetBody(c.BodyRaw()).
			Post(dst)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if rsp.StatusCode() != 200 {
			return c.SendStatus(rsp.StatusCode())
		}
		return c.Send(rsp.Body())
	})

	err := app.Listen(listenAddr)
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
}
