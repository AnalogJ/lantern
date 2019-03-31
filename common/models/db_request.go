package models

import "time"


// Gorm is used for saving/retrieving data from database and mapping to Go structs.
// MapStructure is used for converting DB notify events to Go structs
type DbRequest struct {
	Id            uint                   `gorm:"PRIMARY_KEY;AUTO_INCREMENT" mapstructure:"id"`
	Method        string                 `gorm:"type:varchar(10);NOT NULL" mapstructure:"method"`
	Url           string                 `gorm:"NOT NULL" mapstructure:"url"`
	Headers       HeaderMap 			 `gorm:"type:jsonb;NOT NULL" mapstructure:"headers"`
	Body          string                 `gorm:"NOT NULL" mapstructure:"body"`
	ContentLength int64                  `gorm:"NOT NULL" mapstructure:"content_length"`
	Host          string                 `gorm:"NOT NULL" mapstructure:"host"`
	RequestedOn   time.Time              `gorm:"NOT NULL" mapstructure:"requested_on"`
	CreatedAt     time.Time              `gorm:"NOT NULL" mapstructure:"created_at"`
}
