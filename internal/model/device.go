package model

import (
	"log"
	"ovo-server/internal/database"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Device struct {
	ID        uuid.UUID `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserID    uint           `json:"user_id"`
	User      User           `json:"user"`
	Name      string         `json:"name"`
	Activity  *time.Time     `json:"activity"`
}

func (d *Device) BeforeCreate(tx *gorm.DB) error {
	newUuid, err := uuid.NewV7()
	if err != nil {
		log.Printf("Error generating UUID: %s", err)
		return err
	}
	log.Println("Generated UUID: ", newUuid)
	d.ID = newUuid
	return nil
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

func GetDevice(id uuid.UUID) (Device, error) {
	device := Device{}
	err := database.GetDB().Where("id = ?", id).First(&device).Error
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
