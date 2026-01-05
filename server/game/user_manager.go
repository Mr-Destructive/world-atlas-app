package game

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"passwordHash"` // Simple SHA256(salt + password)
	Salt         string    `json:"salt"`
	CreatedAt    time.Time `json:"createdAt"`
	
	// Stats
	GamesPlayed  int `json:"gamesPlayed"`
	TotalScore   int `json:"totalScore"`
	Wins         int `json:"wins"`
}

type UserManager struct {
	filePath string
	users    map[string]*User // Username -> User
	mu       sync.RWMutex
}

func NewUserManager(path string) (*UserManager, error) {
	um := &UserManager{
		filePath: path,
		users:    make(map[string]*User),
	}
	if err := um.load(); err != nil {
		return nil, err
	}
	return um, nil
}

func (um *UserManager) load() error {
	um.mu.Lock()
	defer um.mu.Unlock()

	data, err := os.ReadFile(um.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	// File might be empty
	if len(data) == 0 {
		return nil
	}

	// We store as a list or map in JSON? Map is easier for lookup
	var loadedUsers map[string]*User
	if err := json.Unmarshal(data, &loadedUsers); err != nil {
		return err
	}
	um.users = loadedUsers
	return nil
}

func (um *UserManager) save() error {
	// Must be called under lock
	data, err := json.MarshalIndent(um.users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(um.filePath, data, 0644)
}

func (um *UserManager) Register(username, password string) (*User, error) {
	um.mu.Lock()
	defer um.mu.Unlock()

	if _, exists := um.users[username]; exists {
		return nil, errors.New("username already taken")
	}

	salt := uuid.New().String()
	hash := hashPassword(password, salt)

	user := &User{
		ID:           uuid.New().String(),
		Username:     username,
		PasswordHash: hash,
		Salt:         salt,
		CreatedAt:    time.Now(),
		GamesPlayed:  0,
		TotalScore:   0,
		Wins:         0,
	}

	um.users[username] = user
	if err := um.save(); err != nil {
		delete(um.users, username)
		return nil, fmt.Errorf("failed to save user: %v", err)
	}

	return user, nil
}

func (um *UserManager) Login(username, password string) (*User, error) {
	um.mu.RLock()
	user, exists := um.users[username]
	um.mu.RUnlock()

	if !exists {
		return nil, errors.New("invalid username or password")
	}

	if hashPassword(password, user.Salt) != user.PasswordHash {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

func (um *UserManager) UpdateStats(username string, score int, isWin bool) {
	um.mu.Lock()
	defer um.mu.Unlock()
	
	if user, ok := um.users[username]; ok {
		user.GamesPlayed++
		user.TotalScore += score
		if isWin {
			user.Wins++
		}
		_ = um.save()
	}
}

func hashPassword(password, salt string) string {
	h := sha256.New()
	h.Write([]byte(salt + password))
	return hex.EncodeToString(h.Sum(nil))
}
