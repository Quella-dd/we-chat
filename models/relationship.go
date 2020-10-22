package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"we-chat/database"
	"we-chat/message"

	"github.com/jinzhu/gorm"
)

// type 属性用于区分申请添加好友还是申请加入群聊
type Request struct {
	gorm.Model
	Type          int
	SourceID      string
	DestinationID string
	Status        bool
}

type RelationShipManager struct {
	Requests []Request
}

func NewRelationShipManager() *RelationShipManager {
	database.DB.AutoMigrate(&Request{})
	database.DB.AutoMigrate(&RelationShips{})
	return &RelationShipManager{}
}

type RelationShips struct {
	gorm.Model
	UserID       string
	RelationShip RelationUsers `sql:"TYPE:json"`
}

type RelationUsers []User

func (c RelationUsers) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c *RelationUsers) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), c)
}

func (r *RelationShipManager) GetFriends(id string) ([]User, error) {
	result, err := r.GetRelation(id)
	if err != nil {
		return nil, err
	}
	return result.RelationShip, nil
}

func (r *RelationShipManager) AddFriend(requestID, destinationID string) error {
	user, err := ManageEnv.UserManager.GetUser(requestID)
	if err != nil {
		return errors.New(fmt.Sprintf("user %s not found", requestID))
	}

	// 如果被添加的用户不需要进行认证时，需要修改添加方和被添加方的relations_ship
	if user.Validator {
		if err := r.UpdateRelationShip(requestID, destinationID); err != nil {
			return err
		}
		if err := r.UpdateRelationShip(destinationID, requestID); err != nil {
			return err
		}
	} else {
		// 如果用户开启了添加好友认证，先向确认放发送一条确认信息，当确认方进行确认后才能添加好友
		request := &Request{
			Type:          message.USERMESSAGE,
			SourceID:      requestID,
			DestinationID: destinationID,
			Status:        false,
		}
		if err := database.DB.Save(request).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *RelationShipManager) DeleteFriend(requestID, destinationID string) error {
	result, err := r.GetRelation(requestID)
	if err != nil {
		return err
	}

	for index, user := range result.RelationShip {
		if string(user.ID) == destinationID {
			result.RelationShip = append(result.RelationShip[:index], result.RelationShip[index:]...)
			break
		}
	}

	err = database.DB.Model(&result).Update("relation_ship", result.RelationShip).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *RelationShipManager) UpdateRelationShip(requestID, destinationID string) error {
	user, err := ManageEnv.UserManager.GetUser(destinationID)
	if err != nil {
		return err
	}

	result, err := r.GetRelation(requestID)
	if err != nil {
		result = &RelationShips{
			UserID: requestID,
		}
		result.RelationShip = RelationUsers{*user}
		return database.DB.Create(&result).Error
	}

	for _, user := range result.RelationShip {
		if string(user.ID) == destinationID {
			return errors.New(fmt.Sprintf("don't add repeat: userName: %s\n", user.Name))
		}
	}
	result.RelationShip = append(result.RelationShip, *user)
	err = database.DB.Model(&result).Update("relation_ship", result.RelationShip).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *RelationShipManager) ListAddFriendRequest(id string) ([]Request, error) {
	var requests, middleRequests []Request
	if err := database.DB.Where("source_id = ?", id).Find(&middleRequests).Error; err == nil {
		requests = append(requests, middleRequests...)
	}

	if err := database.DB.Where("destination_id = ?", id).Find(&middleRequests).Error; err == nil {
		requests = append(requests, middleRequests...)
	}
	return requests, nil
}

func (r *RelationShipManager) AckRequest(requestID, destinationID string) error {
	var request Request
	err := database.DB.Where("id = ?", destinationID).First(&request).Error
	if err != nil {
		return err
	}

	if request.DestinationID != requestID {
		return errors.New("not allow to ack this request")
	}
	err = database.DB.Model(&request).Update("status", true).Error
	if err != nil {
		return err
	}
	if err := r.UpdateRelationShip(requestID, destinationID); err != nil {
		return err
	}
	if err := r.UpdateRelationShip(destinationID, requestID); err != nil {
		return err
	}
	return nil
}

func (r *RelationShipManager) GetRelation(id string) (*RelationShips, error) {
	var result *RelationShips
	if err := database.DB.Where("user_id = ?", id).First(result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
