package models

import (
	"errors"
	"strconv"
	"we-chat/database"

	"github.com/jinzhu/gorm"
)

type GroupManager struct{}

type Group struct {
	gorm.Model
	Name        string `form:"name"`
	ManagerID   string
	Description string        `form:"description"`
	Childes     RelationUsers `sql:"TYPE:json"`
}

func NewGroupManager() *GroupManager {
	database.DB.AutoMigrate(&Group{})
	return &GroupManager{}
}

func (*GroupManager) ListGroup(userID string) ([]Group, error) {
	var groups []Group
	if err := database.DB.Where("manager_id = ?", userID).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (*GroupManager) CreateGroup(userID string, r *Group) error {
	r.ManagerID = userID
	if err := database.DB.Create(r).Error; err != nil {
		return err
	}
	return nil
}

func (*GroupManager) UpdateGroup(group *Group) error {
	if err := database.DB.Model(group).Update("description", group.Description).Error; err != nil {
		return err
	}
	return nil
}

func (*GroupManager) GetGroup(roomID string) (*Group, error) {
	var group Group
	if err := database.DB.Where("id = ?", roomID).Find(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (*GroupManager) DeleteGroup(id string) error {
	var group Group
	if err := database.DB.First(&group, id).Error; err != nil {
		return err
	}

	if group.ManagerID == id {
		return database.DB.Unscoped().Delete(&group).Error
	}
	return errors.New("this user is not allow to delete this room")
}

// action with group
func (userManager *UserManager) DeleteFromGroup(id, groupID, userID string) error {
	_, err := userManager.GetUser(userID)
	if err != nil {
		return err
	}

	group, err := ManageEnv.GroupManager.GetGroup(groupID)
	if err != nil {
		return err
	}

	if group.ManagerID == id {
		for i, user := range group.Childes {
			if strconv.Itoa(int(user.ID)) == id {
				group.Childes = append(group.Childes[:i], group.Childes[i+1:]...)
			}
		}
		err := database.DB.Model(group).Update("childrens", group.Childes).Error
		if err != nil {
			return err
		}
	}
	return errors.New("not allow to delte user from Room")
}

func (userManager *UserManager) AddUserToGroup(id, roomID, userID string) error {
	user, err := userManager.GetUser(userID)
	if err != nil {
		return err
	}

	group, err := ManageEnv.GroupManager.GetGroup(roomID)
	if err != nil {
		return err
	}

	if group.ManagerID == id {
		if group.Childes == nil {
			group.Childes = []User{}
		}
		group.Childes = append(group.Childes, *user)
		err := database.DB.Model(group).Update("Childes", group.Childes).Error
		if err != nil {
			return err
		}
	}
	return errors.New("not allow to delte user from Room")
}
