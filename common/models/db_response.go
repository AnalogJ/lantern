package models

import "time"

// Gorm is used for saving/retrieving data from database and mapping to Go structs.
// MapStructure is used for converting DB notify events to Go structs
type DbResponse struct {
	Id            uint                   `gorm:"PRIMARY_KEY;AUTO_INCREMENT" mapstructure:"id"`
	RequestId     uint                   `gorm:"NOT NULL" mapstructure:"request_id"`
	Status        string                 `gorm:"type:varchar(50);NOT NULL" mapstructure:"status"`
	StatusCode    int                    `gorm:"NOT NULL" mapstructure:"status_code"`
	Headers       HeaderMap				 `gorm:"type:jsonb;NOT NULL" mapstructure:"headers"`
	Body          string                 `gorm:"NOT NULL" mapstructure:"body"`
	ContentLength int64                  `gorm:"NOT NULL" mapstructure:"content_length"`
	MimeType      string                 `gorm:"type:varchar(50);NOT NULL" mapstructure:"mime_type"`
	Protocol      string				 `gorm:"type:varchar(50);NOT NULL" mapstructure:"protocol"`
	RespondedOn   time.Time              `gorm:"NOT NULL" mapstructure:"responded_on"`
	CreatedAt     time.Time              `gorm:"NOT NULL" mapstructure:"created_at"`
}

func (db DbResponse) TableName() string {
	return "responses"
}