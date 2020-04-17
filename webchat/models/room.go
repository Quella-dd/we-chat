package models

import (
	"webchat/database"

	"github.com/jinzhu/gorm"
)

type RoomManager struct{}

type Room struct {
	gorm.Model
	Name        string `form:"name"`
	ManagerID   string
	Description string
}

func NewRoomManager() *RoomManager {
	database.DB.AutoMigrate(&Room{})
	return &RoomManager{}
}

func (*RoomManager) CreateRoom(userID string, r *Room) error {
	r.ManagerID = userID
	if err := database.DB.Create(r).Error; err != nil {
		return err
	}
	return nil
}

func (*RoomManager) ListRooms(userID string) ([]Room, error) {
	var roomes []Room
	if err := database.DB.Where("manager_id = ?", userID).Find(&roomes).Error; err != nil {
		return nil, err
	}
	return roomes, nil
}

func (*RoomManager) UpdateRoom() {

}

func (*RoomManager) GetRoom(roomID string) (*Room, error) {
	var room Room
	if err := database.DB.Where("id = ?", roomID).Find(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (*RoomManager) DeleteRoom() {

}
