package engine

import (
	"NotRedis/internal/myError"
	"go.uber.org/zap"
	"testing"
)

func TestMyMap(t *testing.T) {
	logger := zap.NewNop()
	m := NewMyMap(10, logger)

	tests := []struct {
		name      string
		action    func(*MyMap) error
		key       string
		value     string
		wantValue string
		wantErr   error
	}{
		{
			name: "Set and Get success",
			action: func(m *MyMap) error {
				return m.Set("key1", "value1")
			},
			key:       "key1",
			value:     "value1",
			wantValue: "value1",
			wantErr:   nil,
		},
		{
			name: "Get non-existent key",
			action: func(m *MyMap) error {
				_, err := m.Get("nonexistent")
				return err
			},
			key:     "nonexistent",
			wantErr: myError.KeyNotFound,
		},
		{
			name: "Set with empty key",
			action: func(m *MyMap) error {
				return m.Set("", "value")
			},
			wantErr: myError.EmptyKeyOrValue,
		},
		{
			name: "Set with empty value",
			action: func(m *MyMap) error {
				return m.Set("key", "")
			},
			wantErr: myError.EmptyKeyOrValue,
		},
		{
			name: "Delete success",
			action: func(m *MyMap) error {
				m.Set("key2", "value2")
				return m.Del("key2")
			},
			key:     "key2",
			wantErr: nil,
		},
		{
			name: "Delete non-existent key",
			action: func(m *MyMap) error {
				return m.Del("nonexistent")
			},
			key:     "nonexistent",
			wantErr: myError.KeyNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.action(m)

			if err != tt.wantErr {
				t.Errorf("%s error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if err == nil && tt.name == "Set and Get success" {
				got, err := m.Get(tt.key)
				if err != nil || got != tt.wantValue {
					t.Errorf("Get() = %v, err = %v; want %v, nil", got, err, tt.wantValue)
				}
			}
			if err == nil && tt.name == "Delete success" {
				_, err := m.Get(tt.key)
				if err != myError.KeyNotFound {
					t.Errorf("Key %v should be deleted, but Get() returned err = %v", tt.key, err)
				}
			}
		})
	}
}

func TestMyMapConcurrency(t *testing.T) {
	logger := zap.NewNop()
	m := NewMyMap(10, logger)

	done := make(chan bool)
	go func() {
		for i := 0; i < 100; i++ {
			m.Set("key", "value")
		}
		done <- true
	}()
	go func() {
		for i := 0; i < 100; i++ {
			m.Get("key")
		}
		done <- true
	}()
	go func() {
		for i := 0; i < 100; i++ {
			m.Del("key")
		}
		done <- true
	}()

	for i := 0; i < 3; i++ {
		<-done
	}
}
