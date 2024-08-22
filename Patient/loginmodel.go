package Patient

import "clinic_server/types"

type Login struct {
	Id       types.NullInt64  `gorm:"PrimaryKey"`
	UserName types.NullString `json:"UserName"`
	Password types.NullString `json:"Password"`
}
