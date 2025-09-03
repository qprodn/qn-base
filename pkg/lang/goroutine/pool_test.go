// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package goroutine

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/bytedance/gg/gslice"
	"github.com/stretchr/testify/assert"
)

func TestNewPool(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		wantErr bool
	}{
		{
			name:    "valid pool size",
			size:    10,
			wantErr: false,
		},
		{
			name:    "invalid pool size",
			size:    -1,
			wantErr: true,
		},
		{
			name:    "zero pool size",
			size:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, err := NewPool(tt.size)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, pool)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, pool)
			}
		})
	}
}

func TestPool_Exec(t *testing.T) {
	t.Run("successful execution", func(t *testing.T) {
		ctx := context.Background()
		var result int

		pool, err := NewPool(2)
		assert.NoError(t, err)

		pool.Add(func() error {
			result = 1
			return nil
		})

		err = pool.Exec(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 1, result)
	})

	t.Run("error execution", func(t *testing.T) {
		ctx := context.Background()
		expectedErr := errors.New("test error")

		pool, err := NewPool(2)
		assert.NoError(t, err)

		pool.Add(func() error {
			return expectedErr
		})

		err = pool.Exec(ctx)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		pool, err := NewPool(2)
		assert.NoError(t, err)

		pool.Add(func() error {
			time.Sleep(100 * time.Millisecond)
			return nil
		})

		err = pool.Exec(ctx)
		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)
	})
}

func TestPool_ExecAll(t *testing.T) {
	t.Run("continue on error", func(t *testing.T) {
		ctx := context.Background()
		var (
			results []int
			mutex   = sync.Mutex{}
		)

		pool, err := NewPool(2)
		assert.NoError(t, err)

		pool.Add(func() error {
			mutex.Lock()
			results = append(results, 1)
			mutex.Unlock()
			return errors.New("first error")
		})

		pool.Add(func() error {
			mutex.Lock()
			results = append(results, 2)
			mutex.Unlock()
			return nil
		})

		err = pool.ExecAll(ctx)
		assert.Error(t, err)
		assert.Equal(t,
			gslice.ToMap([]int{1, 2}, func(v int) (int, bool) { return v, true }),
			gslice.ToMap(results, func(v int) (int, bool) { return v, true }),
		)
	})
}

func Test_pool_execute(t *testing.T) {
	t.Run("execute tasks with pool size equal to task count", func(t *testing.T) {
		ctx := context.Background()
		size := 5

		p, err := NewPool(size)
		assert.NoError(t, err)

		for i := 0; i < size; i++ {
			p.Add(func() error {
				time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
				return nil
			})
		}

		err = p.Exec(ctx)
		assert.NoError(t, err)
	})

	t.Run("execute tasks with pool size greater than task count", func(t *testing.T) {
		ctx := context.Background()
		size := 5

		p, err := NewPool(size)
		assert.NoError(t, err)

		for i := 0; i < size>>1; i++ {
			p.Add(func() error {
				time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
				return nil
			})
		}

		err = p.Exec(ctx)
		assert.NoError(t, err)
	})

	t.Run("execute tasks with pool size less than task count", func(t *testing.T) {
		ctx := context.Background()
		size := 5

		p, err := NewPool(size)
		assert.NoError(t, err)

		for i := 0; i < size<<1; i++ {
			p.Add(func() error {
				time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
				return nil
			})
		}

		err = p.Exec(ctx)
		assert.NoError(t, err)
	})

	t.Run("execute tasks with error", func(t *testing.T) {
		ctx := context.Background()
		size := 5

		p, err := NewPool(size)
		assert.NoError(t, err)

		for i := 0; i < size; i++ {
			idx := i

			p.Add(func() error {
				time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
				if idx == 2 {
					return errors.New("test err")
				}
				return nil
			})
		}

		err = p.Exec(ctx)
		assert.Error(t, err)
		assert.Equal(t, "test err", err.Error())
	})
}
