package PatientCare

import (
	"clinic_server/database"
	"clinic_server/logger"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	maxLoginAttempts   = 5
	loginLockoutPeriod = 15 * time.Minute
	sessionTimeout     = 24 * time.Hour
	userCacheDuration  = 1 * time.Hour
)

// Thread-safe cache implementations
var (
	loginAttempts = &LoginAttemptCache{attempts: make(map[string]*AttemptInfo)}
	userCache     = &UserCache{users: make(map[string]*CachedUser)}

	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountLocked      = errors.New("account locked")
)

// Clean up expired cache entries periodically
func init() {
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			cleanupCaches()
		}
	}()
}

func cleanupCaches() {
	now := time.Now()

	// Cleanup login attempts
	loginAttempts.Lock()
	for username, info := range loginAttempts.attempts {
		if now.After(info.LockedUntil) {
			delete(loginAttempts.attempts, username)
		}
	}
	loginAttempts.Unlock()

	// Cleanup user cache
	userCache.Lock()
	for username, cached := range userCache.users {
		if now.After(cached.ExpiresAt) {
			delete(userCache.users, username)
		}
	}
	userCache.Unlock()
}

func ValidateLogin(c *gin.Context) {
	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	// Check for rate limiting
	if isRateLimited(loginReq.UserName) {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many login attempts"})
		return
	}

	// Try to get user from cache first
	user := getUserFromCache(loginReq.UserName)
	if user == nil {
		// If not in cache, get from database
		var err error
		user, err = getUserFromDB(loginReq.UserName)
		if err != nil {
			handleLoginError(c, loginReq.UserName, err)
			return
		}
		// Cache the user data
		cacheUser(user)
	}

	// Validate password and handle login
	if err := validateAndHandleLogin(user, loginReq.Password); err != nil {
		handleLoginError(c, loginReq.UserName, err)
		return
	}

	// Generate session token and create response
	token, expiresAt, err := createSession(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	response := LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		UserInfo: UserInfo{
			ID:       user.ID,
			UserName: user.UserName,
			Role:     user.Role,
		},
	}

	// Update last login time asynchronously
	go updateLastLogin(user.ID)

	c.JSON(http.StatusOK, response)
}

func getUserFromCache(username string) *LoginModel {
	userCache.RLock()
	defer userCache.RUnlock()

	if cached, exists := userCache.users[username]; exists && time.Now().Before(cached.ExpiresAt) {
		return cached.User
	}
	return nil
}

func cacheUser(user *LoginModel) {
	userCache.Lock()
	defer userCache.Unlock()

	userCache.users[user.UserName] = &CachedUser{
		User:      user,
		ExpiresAt: time.Now().Add(userCacheDuration),
	}
}

func isRateLimited(username string) bool {
	loginAttempts.RLock()
	defer loginAttempts.RUnlock()

	if info, exists := loginAttempts.attempts[username]; exists {
		return info.Count >= maxLoginAttempts && time.Now().Before(info.LockedUntil)
	}
	return false
}

func incrementLoginAttempts(username string) {
	loginAttempts.Lock()
	defer loginAttempts.Unlock()

	now := time.Now()
	info, exists := loginAttempts.attempts[username]
	if !exists {
		info = &AttemptInfo{}
		loginAttempts.attempts[username] = info
	}

	info.Count++
	info.LastTry = now
	if info.Count >= maxLoginAttempts {
		info.LockedUntil = now.Add(loginLockoutPeriod)
	}
}

func resetLoginAttempts(username string) {
	loginAttempts.Lock()
	defer loginAttempts.Unlock()
	delete(loginAttempts.attempts, username)
}

func validateAndHandleLogin(user *LoginModel, password string) error {
	if user.Status == "locked" {
		return ErrAccountLocked
	}

	if !validatePassword(password, user.Password) {
		incrementLoginAttempts(user.UserName)
		return ErrInvalidCredentials
	}

	resetLoginAttempts(user.UserName)
	return nil
}

func validatePassword(inputPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	return err == nil
}

func getUserFromDB(username string) (*LoginModel, error) {
	db := database.GetDBContext()
	if db == nil {
		return nil, errors.New("database connection failed")
	}

	var user LoginModel
	err := db.Where("user_name = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	return &user, nil
}

func updateLastLogin(userID uint) {
	db := database.GetDBContext()
	if db == nil {
		logger.Error("Failed to update last login time", errors.New("database connection failed"))
		return
	}

	if err := db.Model(&LoginModel{}).Where("id = ?", userID).
		Update("last_login", time.Now()).Error; err != nil {
		logger.Error("Failed to update last login time", err)
	}
}

func handleLoginError(c *gin.Context, username string, err error) {
	switch err {
	case ErrInvalidCredentials:
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
	case ErrAccountLocked:
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "account locked"})
	default:
		logger.Error("Login error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

func createSession(user *LoginModel) (string, time.Time, error) {
	expiresAt := time.Now().Add(sessionTimeout)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     expiresAt.Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}
