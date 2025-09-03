package idgen

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/sony/sonyflake"
	"qn-base/pkg/util/idgen"
)

// IDGenerator is a service for generating unique IDs
type IDGenerator struct {
	logger *log.Helper
}

// NewIDGenerator creates a new ID generator service
func NewIDGenerator(logger log.Logger) (*IDGenerator, error) {
	// Initialize the ID generator with default settings
	idgen.Init(sonyflake.Settings{})

	return &IDGenerator{
		logger: log.NewHelper(logger),
	}, nil
}

// NextID generates and returns the next unique ID
func (ig *IDGenerator) NextID() (uint64, error) {
	id, err := idgen.NextID()
	if err != nil {
		ig.logger.Errorf("failed to generate snowflake ID: %v", err)
		return 0, fmt.Errorf("failed to generate snowflake ID: %w", err)
	}
	return id, nil
}

// NextStringID generates and returns the next unique ID as a string
func (ig *IDGenerator) NextStringID() (string, error) {
	id, err := ig.NextID()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", id), nil
}
