package idgen

import (
	"fmt"
	"sync"

	"github.com/sony/sonyflake"
)

var (
	once     sync.Once
	instance *sonyflake.Sonyflake
)

// Init initializes the ID generator with the given settings
func Init(settings sonyflake.Settings) {
	once.Do(func() {
		instance = sonyflake.NewSonyflake(settings)
		if instance == nil {
			panic("failed to initialize sonyflake")
		}
	})
}

// GetInstance returns the singleton instance of the ID generator
func GetInstance() *sonyflake.Sonyflake {
	if instance == nil {
		// Initialize with default settings if not initialized
		Init(sonyflake.Settings{})
	}
	return instance
}

// NextID generates and returns the next unique ID
func NextID() (uint64, error) {
	sf := GetInstance()
	if sf == nil {
		return 0, fmt.Errorf("id generator not initialized")
	}
	return sf.NextID()
}

// NextStringID generates and returns the next unique ID as a string
func NextStringID() (string, error) {
	id, err := NextID()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", id), nil
}
