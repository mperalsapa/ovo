package model

import (
	"fmt"
	"log"
	db "ovo-server/internal/database"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Role int

const (
	Visitor Role = iota
	Editor
	Admin
)

type User struct {
	gorm.Model
	Username string `form:"username" json:"username" gorm:"not null"`
	Password string `form:"password" json:"password" gorm:"not null"`
	Role     Role   `json:"role"`
	// WatchedMovies  []Movie   `gorm:"many2many:user_watched_movies;"`
	// WatchedEpisode []Episode `gorm:"many2many:user_watched_episodes;"`
	WatchedItems  []Item `gorm:"many2many:user_watched_items;"`
	FavoriteItems []Item `gorm:"many2many:user_favorite_items;"`
	Enabled       bool   `json:"enabled" gorm:"default:false"`
}

func (u *User) HashPassword() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hashedPassword)
}

// Compares the hashed password (stored in the database) with the password provided by the user
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func NewUser(username, password string) User {
	user := User{
		Username: username,
		Password: password,
	}

	user.HashPassword()
	return user
}

func (u *User) Save() {
	db.GetDB().Save(u)
}

func (u *User) Delete() {
	db.GetDB().Delete(u)
}

func UserCount() int64 {
	var count int64
	db.GetDB().Model(&User{}).Count(&count)
	return count
}

func GetUserByID(id uint) User {
	user := User{}
	db.GetDB().Where("id = ?", id).First(&user)
	fmt.Println(user)
	return user
}

func GetUserByUsername(username string) User {
	user := User{}
	db.GetDB().Where("username = ?", username).First(&user)
	return user
}

func GetUserExists(username string) bool {
	user := GetUserByUsername(username)
	return user.Username != ""
}

func (u *User) ToggleFavoriteItem(itemID uint) bool {
	item, err := GetItemById(itemID)
	if item.ID == 0 || err != nil {
		return false
	}

	u.FetchFavoriteItems()
	for _, favoriteItem := range u.FavoriteItems {
		if favoriteItem.ID == item.ID {
			db.GetDB().Model(&u).Association("FavoriteItems").Delete(item)
			u.Save()
			return false
		}
	}

	u.FavoriteItems = append(u.FavoriteItems, item)
	u.Save()
	return true
}

func (u *User) FetchFavoriteItems() {
	db.GetDB().Model(&u).Association("FavoriteItems").Find(&u.FavoriteItems)
}

func (u *User) ItemIsFavorite(itemID uint) bool {
	var favoriteID uint
	db.GetDB().Model(&u).Where("item_id = ?", itemID).Association("FavoriteItems").Find(&u.FavoriteItems)
	if len(u.FavoriteItems) > 0 {
		favoriteID = u.FavoriteItems[0].ID
	}
	log.Println(favoriteID)
	log.Println(len(u.FavoriteItems))
	return favoriteID != 0
}

func (u *User) ToggleWatchedItem(itemID uint) bool {
	item, err := GetItemById(itemID)
	if item.ID == 0 || err != nil {
		return false
	}

	u.FetchWatchedItems()
	for _, watchedItem := range u.WatchedItems {
		if watchedItem.ID == item.ID {
			db.GetDB().Model(&u).Association("WatchedItems").Delete(item)
			u.Save()
			return false
		}
	}

	u.WatchedItems = append(u.WatchedItems, item)
	u.Save()
	return true
}

func (u *User) FetchWatchedItems() {
	db.GetDB().Model(&u).Association("WatchedItems").Find(&u.WatchedItems)
}

func (u *User) ItemIsWatched(itemID uint) bool {
	var watchedID uint
	db.GetDB().Model(&u).Where("item_id = ?", itemID).Association("WatchedItems").Find(&u.WatchedItems)
	if len(u.WatchedItems) > 0 {
		watchedID = u.WatchedItems[0].ID
	}
	log.Println(watchedID)
	log.Println(len(u.WatchedItems))
	return watchedID != 0
}
