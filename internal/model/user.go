package model

import (
	"fmt"
	db "ovo-server/internal/database"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Role string

const (
	Admin   Role = "admin"
	Editor  Role = "editor"
	Visitor Role = "visitor"
)

type User struct {
	gorm.Model
	Username       string    `form:"username" json:"username" gorm:"not null"`
	Password       string    `form:"password" json:"password" gorm:"not null"`
	Role           Role      `json:"role" gorm:"enum('admin', 'editor', 'visitor');default:'visitor'"`
	WatchedMovies  []Movie   `gorm:"many2many:user_watched_movies;"`
	WatchedEpisode []Episode `gorm:"many2many:user_watched_episodes;"`
	Enabled        bool      `json:"enabled" gorm:"default:false"; default:false`
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

func CreateUser(username, password string) User {
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
	fmt.Println(user)
	return user
}
