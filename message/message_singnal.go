///*
//	如果接收方接受了视频通话，则在视频通话结束的时候记录这条日志
//*/
//
package message

import (
	"time"

	"github.com/jinzhu/gorm"
)

type SingnalMesage struct {
}

func (e *SingnalMesage) FormatMessage() {

}

// 存储于redis中，并设置过期时间，若到了过期事件，则将该事件进行删除，并且向前端进行提示
type RTCEvent struct {
	ID string
	gorm.Model
	SourceID      int `json: "sourceID"`
	DestinationID int `json: "destinationID"`
	Type          string
	StartTime time.Time
	EndTime   time.Time
	Status    string
}

//func (e *RTCEvent) NewSessioMessage() *MessageInfo {
//	content, err := json.Marshal(struct{
//		Status string
//		SubTime time.Duration
//	}{e.Status, e.EndTime.Sub(e.StartTime)})
//
//	if err != nil {
//		return nil
//	}
//	return &MessageInfo{
//		SourceID: e.SourceID,
//		DestinationID: e.DestinationID,
//		Create_At: e.CreatedAt.Unix(),
//		Content: string(content),
//	}
//}
