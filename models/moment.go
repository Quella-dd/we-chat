package models

import (
	"strconv"
	"we-chat/common"

	"github.com/jinzhu/gorm"
)

type MomentManager struct{}

type Moment struct {
	gorm.Model
	OwnerID string
	Content string
	Stars   common.RelationStruct `gorm:"type:json"`
}

type MomentInfo struct {
	*Moment
	OwnerName string
	StarsMap  map[int]string
}

type UpdateSpec struct {
	Action string
}

func NewMomentManager() *MomentManager {
	return &MomentManager{}
}

func (m *MomentManager) List(userID string) ([]*MomentInfo, error) {
	var moments []*Moment
	var momentInfo []*MomentInfo

	user, err := ManagerEnv.UserManager.GetUser(userID, "id")
	if err != nil {
		return nil, err
	}

	relations := append(user.Relations, strconv.Itoa(int(user.ID)))

	err = ManagerEnv.DB.Where("owner_id IN (?)", []string(relations)).Find(&moments).Error
	if err != nil {
		return nil, err
	}

	for _, v := range moments {
		info, err := m.GetInfo(strconv.Itoa(int(v.ID)))
		if err != nil {
			continue
		}
		momentInfo = append(momentInfo, info)
	}

	return momentInfo, nil
}

func (m *MomentManager) Create(moment Moment) error {
	err := ManagerEnv.DB.Create(&moment).Error
	if err != nil {
		return err
	}
	return nil
}

func (m *MomentManager) Get(id string) (*Moment, error) {
	var moment Moment
	err := ManagerEnv.DB.Where("id = ?", id).Find(&moment).Error

	if err != nil {
		return nil, err
	}

	return &moment, nil
}

func (m *MomentManager) GetInfo(id string) (*MomentInfo, error) {
	var info MomentInfo
	moment, err := m.Get(id)
	if err != nil {
		return nil, nil
	}

	user, err := ManagerEnv.UserManager.GetUser(moment.OwnerID, "id")
	if err != nil {
		return nil, err
	}

	info.OwnerName = user.Name
	info.Moment = moment

	for _, id := range moment.Stars {
		user, err := ManagerEnv.UserManager.GetUser(id, "id")
		if err != nil {
			continue
		}

		if info.StarsMap == nil {
			info.StarsMap = make(map[int]string)
		}
		info.StarsMap[int(user.ID)] = user.Name
	}
	return &info, nil
}

func (m *MomentManager) Update(id, userID string, spec UpdateSpec) error {
	moment, err := m.Get(id)
	if err != nil {
		return err
	}

	var stars common.RelationStruct
	switch spec.Action {
	case "ADD":
		stars = append(moment.Stars, userID)
	case "DELETE":
		var tmpStars common.RelationStruct
		for _, v := range moment.Stars {
			if v != userID {
				tmpStars = append(tmpStars, v)
			}
		}
		stars = tmpStars
	}
	return ManagerEnv.DB.Where("id = ?", id).Update("stars", stars).Error
}

func (m *MomentManager) Delete(id string) error {
	return ManagerEnv.DB.Delete(&Moment{}, id).Error
}
