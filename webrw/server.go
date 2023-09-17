package webrw

import "github.com/labstack/echo/v4"

func New() {
	e := echo.New()

	a := func(c echo.Context) error {
		c.Cookie()
	}
}
