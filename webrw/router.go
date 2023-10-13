package webrw

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
)

func New() {
	e := echo.New()
	a := func(c echo.Context) error {
		c.Cookie()
		c.HTML()
		c.Stream()
		c.Redirect()
		http.FileServer()
		mux := http.NewServeMux()
		http.FileServer()
		m := gin.Default()
		m.POST("/", func(context *gin.Context) {
			context.Bind()
		})
		http.StatusOK
		c.Bind()
		c.Response()
		e := echo.New()
		ge := gin.Error{}
		e.HTTPErrorHandler
	}
}
