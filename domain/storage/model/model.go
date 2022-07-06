package model

type File struct {
	Id     string `gorm:"primaryKey"`
	UserId string
}
