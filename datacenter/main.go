package datacenter

import "time"

type DataCenterManager struct{
	Count int
	TimeOut time.Time
	Publisher *pubsub.Publisher
}

func NewDataCenterManager() *DataCenterManager{
	database.DB.AutoMigrate(&SessionMessage{})
	return &DataManager{
		Publisher: pubsub.NewPublisher(timeout, buffer),
		Count: 10,
		TimeOut: time.Hour,
	}
}

type SessionMessage struct {
	ID int
	Count int
	SourceID int
	DestinationID int
	Create_At time.Time
	Latest_Visit time.Time
	MessageBody MessagesBody ` sql:"TYPE:json"` 
}

type MessagesBody []MesageBody

type MesageBody struct {
	Content string
}

func (dataCenter *DataCenterManager) Save(msg SessionMessage) {

}
