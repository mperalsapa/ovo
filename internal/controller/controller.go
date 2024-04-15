package controller

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

// Function copied from templ docs: https://github.com/a-h/templ/blob/main/examples/integration-echo/main.go
func RenderView(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	ctx.Response().Writer.WriteHeader(statusCode)
	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}
