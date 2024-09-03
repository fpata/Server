package Patient

import "clinic_server/types"

type LoginModel struct {
	ID       types.NullInt64  `gorm:"PrimaryKey"`
	UserName types.NullString `json:"UserName"`
	Password types.NullString `json:"Password"`
	Role     types.NullString `json:"Role"`
	FullName types.NullString `json:"FullName"`
}
