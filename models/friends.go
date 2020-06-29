package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"webchat/common"
	"webchat/database"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type RelationShipManager struct{}

func NewRelationShipManager() *RelationShipManager {
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
	retID, err := strconv.Atoi(userID)
	if err != nil {
		common.Http404Response(ctx, err)
		return
	}

	ret1ID, err := strconv.Atoi(id)
	if err != nil {
		common.Http404Response(ctx, err)
		return
	}

	user, err := ManageEnv.UserManager.getUser(id)
	result, err := r.GetRelationByID(retID)

	if err == nil {
		for _, user := range result.RelationShip {
			if int(user.ID) == ret1ID {
				common.Http404Response(ctx, fmt.Sprintf("don't add repeat: userName: %s\n", user.Name))
				return
			}
		}
		result.RelationShip = append(result.RelationShip, *user)
		database.DB.Model(&result).Update("relation_ship", result.RelationShip)
	} else {
		result = &RelationShips{
			UserID: retID,
		}
		result.RelationShip = RelationUsers{*user}
		database.DB.Create(&result)
	}
}

func (r *RelationShipManager) DeleteFriend(c *gin.Context, userID string, id string) {

}

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
