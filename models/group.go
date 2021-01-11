package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type GroupManager struct{}

type Group struct {
	gorm.Model
	Name        string `form:"name"`
	ManagerID   string
	Description string        `form:"description"`
	Childes RelationStruct `gorm:"type:json"`
}

func NewGroupManager() *GroupManager {
	return &GroupManager{}
}

func (*GroupManager) ListGroups(userID string) ([]Group, error) {
	var groups []Group
	if err := ManagerEnv.DB.Where("manager_id = ?", userID).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (*GroupManager) CreateGroup(userID string, r *Group) error {
	r.ManagerID = userID
	return ManagerEnv.DB.Create(r).Error
}

func (*GroupManager) UpdateGroup(group *Group) error {
	return ManagerEnv.DB.Model(group).Update("description", group.Description).Error
}

func (*GroupManager) GetGroup(id string) (*Group, error) {
	var group Group
	if err := ManagerEnv.DB.Where("id = ?", id).Find(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (*GroupManager) DeleteGroup(id string) error {
	var group Group
	if err := ManagerEnv.DB.First(&group, id).Error; err != nil {
		return err
	}

	if group.ManagerID == id {
		return ManagerEnv.DB.Unscoped().Delete(&group).Error
	}
	return errors.New("this user is not allow to delete this room")
}

// action with group
func (userManager *UserManager) JoinGroup(id, roomID, userID string) error {
	_, err := userManager.GetUser(userID, "id")
	if err != nil {
		return err
	}

	group, err := ManagerEnv.GroupManager.GetGroup(roomID)
	if err != nil {
		return err
	}

	if group.ManagerID == id {
		if group.Childes == nil {
			group.Childes = RelationStruct{}
		}
		group.Childes = append(group.Childes, userID)
		err := ManagerEnv.DB.Model(group).Update("Childes", group.Childes).Error
		if err != nil {
			return err
		}
	}
	return errors.New("not allow to delete user from Room")
}

func (userManager *UserManager) LeaveGroup(id, groupID, userID string) error {
	_, err := userManager.GetUser(userID, "id")
	if err != nil {
		return err
	}

	group, err := ManagerEnv.GroupManager.GetGroup(groupID)
	if err != nil {
		return err
	}

	if group.ManagerID == id {
		for i, userID := range group.Childes {
			if userID == id {
				group.Childes = append(group.Childes[:i], group.Childes[i+1:]...)
			}
		}
		err := ManagerEnv.DB.Model(group).Update("childrenes", group.Childes).Error
		if err != nil {
			return err
		}
	}
	return errors.New("not allow to delete user from Group")
}
