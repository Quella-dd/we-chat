package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"webchat/common"
	"webchat/database"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// type 属性用于区分申请添加好友还是申请加入群聊
type Request struct {
	gorm.Model

	Type          int
	SourceID      int
	DestinationID int
	Status        int
}

type RelationShipManager struct {
	Requests []Request
}

func NewRelationShipManager() *RelationShipManager {
	database.DB.AutoMigrate(&Request{})
	database.DB.AutoMigrate(&RelationShips{})
	return &RelationShipManager{}
}

func (r *RelationShipManager) GetFriends(ctx *gin.Context, id string) {
	userID, err := strconv.Atoi(id)
	if err != nil {
		common.Http404Response(ctx, err)
		return
	}

	var result RelationShips
	if err := database.DB.Where("user_id = ?", userID).First(&result).Error; err != nil {
		common.Http404Response(ctx, "404 not found")
		return
	}
	common.HttpSuccessResponse(ctx, result.RelationShip)
}

func (r *RelationShipManager) AddFriend(ctx *gin.Context, userID string, id string) {
	relationRequestID, err := strconv.Atoi(userID)
	if err != nil {
		common.Http404Response(ctx, err)
		return
	}

	destinationID, err := strconv.Atoi(id)
	if err != nil {
		common.Http404Response(ctx, err)
		return
	}

	if user, err := ManageEnv.UserManager.getUserByID(destinationID); err != nil {
		panic(err)
	} else {
		// 如果被添加的用户不需要进行认证时，需要修改添加方和被添加放的relations_ship
		if user.NoValidator {
			if err := r.UpdateRelationShip(relationRequestID, destinationID); err != nil {
				common.HttpServerError(ctx, err)
				return
			} else {
				if err := r.UpdateRelationShip(destinationID, relationRequestID); err != nil {
					common.HttpServerError(ctx, err)
					return
				}
			}
		} else {
			// 如果用户开启了添加好友认证，先向确认放发送一条确认信息，当确认方进行确认后才能添加好友
			request := &Request{
				Type:          common.USERMESSAGE,
				SourceID:      relationRequestID,
				DestinationID: destinationID,
				Status: common.UnConfirmd,
			}
			if err := database.DB.Save(request).Error; err != nil {
				panic(err)
			}
		}
	}
}

func (r *RelationShipManager) UpdateRelationShip(userID, destinationID int) error {
	if user, err := ManageEnv.UserManager.getUserByID(destinationID); err != nil {
		return errors.New("user not found")
	} else {
		result, err := r.GetRelationByID(userID)
		if err == nil {
			for _, user := range result.RelationShip {
				if int(user.ID) == destinationID {
					return errors.New(fmt.Sprintf("don't add repeat: userName: %s\n", user.Name))
				}
			}
			result.RelationShip = append(result.RelationShip, *user)
			database.DB.Model(&result).Update("relation_ship", result.RelationShip)
		} else {
			result = &RelationShips{
				UserID: userID,
			}
			result.RelationShip = RelationUsers{*user}
			database.DB.Create(&result)
		}
	}
	return nil
}

func (r *RelationShipManager) DeleteFriend(c *gin.Context, userID string, id string) {

}


/////////////////////////////////////////
func (r *RelationShipManager) ListRequest(ctx *gin.Context, userID string) ([]Request, error) {
	var requests, middleRequests []Request
	if err := database.DB.Where("source_id = ?", userID).Find(&middleRequests).Error; err == nil {
		requests = append(requests, middleRequests...)
	}

	if err := database.DB.Where("destination_id = ?", userID).Find(&middleRequests).Error; err == nil {
		requests = append(requests, middleRequests...)
	}
	return requests, nil
}

func (r *RelationShipManager) AckRequest(ctx *gin.Context, userID, id string) error {
	var request Request
	if err := database.DB.Where("id = ?", id).First(&request).Error; err != nil {
		return err
	}

	i, _ := strconv.Atoi(userID)
	if request.DestinationID != i {
		return errors.New("not allow to ack this request")
	}
	if err := database.DB.Model(&request).Update("status", common.Confirmd).Error; err != nil {
		return err
	}

	sourceID, _ := strconv.Atoi(userID)
	destinationID, _ := strconv.Atoi(id)

	r.UpdateRelationShip(sourceID, destinationID)
	r.UpdateRelationShip(destinationID, sourceID)
	return nil
}
/////////////////////////////////////////////

func (r *RelationShipManager) GetRelationByID(id int) (*RelationShips, error) {
	var result RelationShips
	if err := database.DB.Where("user_id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

type RelationShips struct {
	gorm.Model

	UserID       int
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
