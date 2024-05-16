package middleware

import (
	"log"
	"ovo-server/internal/model"
	"ovo-server/internal/session"

	"github.com/labstack/echo/v4"
	"github.com/mileusna/useragent"
)

func UpdateDeviceActivity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var device model.Device
		var err error

		userSession := session.GetUserSession(c)
		deviceID := userSession.DeviceID

		if deviceID == 0 {
			parsedUA := useragent.Parse(c.Request().UserAgent())

			device = model.CreateDevice(model.GetUserByUsername(userSession.Username).ID, parsedUA.Name)
			userSession.DeviceID = device.ID
			userSession.SaveUserSession(c)
		} else {
			device, err = model.GetDeviceById(deviceID)
			if err != nil {
				log.Println("Error getting device: ", err)
				return next(c)
			}
		}

		device.UpdateDeviceActivity()
		return next(c)
	}
}
