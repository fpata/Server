package PatientCare

import (
	"sync"
	"time"
)

type LoginModel struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserName  string    `json:"username" binding:"required" gorm:"uniqueIndex"`
	Password  string    `json:"-"` // Never send password in response
	Role      string    `json:"role"`
	LastLogin time.Time `json:"lastLogin"`
	Status    string    `json:"status"`
}

type LoginRequest struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	UserInfo  UserInfo  `json:"userInfo"`
}

type LoginAttemptCache struct {
	sync.RWMutex
	attempts map[string]*AttemptInfo
}

type UserCache struct {
	sync.RWMutex
	users map[string]*CachedUser
}

type AttemptInfo struct {
	Count       int
	LastTry     time.Time
	LockedUntil time.Time
}

type CachedUser struct {
	User      *LoginModel
	ExpiresAt time.Time
}

type UserInfo struct {
	ID       uint   `json:"id"`
	UserName string `json:"username"`
	Role     string `json:"role"`
}
