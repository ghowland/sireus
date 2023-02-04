package util

import (
	"github.com/gofiber/fiber/v2"
	"net/url"
)

func ParseContextBody(c *fiber.Ctx) map[string]string {
	values, err := url.ParseQuery(string(c.Body()))
	Check(err)

	obj := map[string]string{}
	for k, v := range values {
		if len(v) > 0 {
			obj[k] = v[0]
		}
	}

	return obj
}
