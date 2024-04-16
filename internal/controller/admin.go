package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func AdminDashboard(context echo.Context) error {
	return context.String(http.StatusOK, "Admin Dashboard")
}
