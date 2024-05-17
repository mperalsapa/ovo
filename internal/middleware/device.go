package middleware

import (
	"log"
	"ovo-server/internal/controller"
	"ovo-server/internal/model"
	"ovo-server/internal/session"

	"github.com/labstack/echo/v4"
)

func UpdateDeviceActivity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var device model.Device
		var err error

		userSession := session.GetUserSession(c)
		deviceID := userSession.DeviceID

		if deviceID != 0 {
			device, err = model.GetDeviceById(deviceID)
			if err != nil {
				log.Printf("User %s tried to acces using invalid deviceID %d. Loggin out.", userSession.Username, deviceID)
				// As we could not find the device, we will log out the user
				// In a future we could add the session into the database, instead
				// of using cookies. This way we could invalidate the session directly
				controller.Logout(c)
				return nil
			}
		}

		device.UpdateDeviceActivity()
		return next(c)
	}
}
