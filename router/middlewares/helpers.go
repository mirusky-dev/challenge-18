package middlewares

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrValueNotFound = errors.New("VALUE not found")
)

type Extractor func(c *fiber.Ctx) (string, error)

// FromHeader returns a function that extracts VALUE from the request header.
func FromHeader(header string, prefix string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		header := c.Get(header)
		l := len(prefix)
		if len(header) > l+1 && strings.EqualFold(header[:l], prefix) {
			return header[l+1:], nil
		}
		return "", ErrValueNotFound
	}
}

// FromHeader returns a function that extracts VALUE from the request hostname.
func FromSubDomain(param string) func(c *fiber.Ctx) (string, error) {
	paramInt, err := strconv.Atoi(param)
	if err != nil {
		panic("FromSubDomain: the param MUST be a valid int")
	}

	return func(c *fiber.Ctx) (string, error) {
		split := strings.Split(c.Hostname(), ".")

		if paramInt >= 0 && paramInt < len(split) {
			return split[paramInt], nil
		}

		return "", ErrValueNotFound
	}
}

// FromQuery returns a function that extracts VALUE from the query string.
func FromQuery(param string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		value := c.Query(param)
		if value == "" {
			return "", ErrValueNotFound
		}
		return value, nil
	}
}

// FromParam returns a function that extracts VALUE from the url param string.
func FromParam(param string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		value := c.Params(param)
		if value == "" {
			return "", ErrValueNotFound
		}
		return value, nil
	}
}

// FromCookie returns a function that extracts VALUE from the named cookie.
func FromCookie(name string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		value := c.Cookies(name)
		if value == "" {
			return "", ErrValueNotFound
		}
		return value, nil
	}
}
