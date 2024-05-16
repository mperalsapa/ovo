package model

import (
	"ovo-server/internal/database"
	"time"
)

type Device struct {
	ID       uint       `json:"id"`
	UserID   uint       `json:"user_id"`
	User     User       `json:"user"`
	Name     string     `json:"name"`
	Activity *time.Time `json:"activity"`
}

func CreateDevice(userId uint, name string) Device {
	device := Device{}
	device.UserID = userId
	device.Name = name
	currentTime := time.Now()
	device.Activity = &currentTime
	device.Save()
	return device
}

func (d *Device) Save() error {
	return database.GetDB().Save(d).Error
}

func GetDevices() []Device {
	devices := []Device{}
	database.GetDB().Find(&devices)
	return devices
}

func GetDeviceById(id uint) (Device, error) {
	device := Device{}
	err := database.GetDB().First(&device, id).Error
	return device, err
}

func GetDevicesByUserId(userId uint) []Device {
	devices := []Device{}
	database.GetDB().Where("user_id = ?", userId).Find(&devices)
	return devices
}

func (d *Device) UpdateDeviceActivity() error {
	return database.GetDB().Model(&d).Update("activity", time.Now()).Error
}
