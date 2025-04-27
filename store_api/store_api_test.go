package store_api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStore_SetGet(t *testing.T) {
	s := NewStore()
	key := "foo"
	value := "bar"

	s.Set(key, value, 0)

	result := s.Get(key)
	assert.Equal(t, value, result)
}

func TestStore_SetWithTTL(t *testing.T) {
	s := NewStore()
	key := "temp"
	value := "value"

	s.Set(key, value, 1*time.Second)

	// Should be available immediately
	result := s.Get(key)
	assert.Equal(t, value, result)

	// Wait for TTL to expire
	time.Sleep(2 * time.Second)

	result = s.Get(key)
	assert.Nil(t, result)
}

func TestStore_Update(t *testing.T) {
	s := NewStore()
	key := "foo"
	initial := "bar"
	updated := "baz"

	s.Set(key, initial, 0)

	success := s.Update(key, updated)
	assert.True(t, success)

	result := s.Get(key)
	assert.Equal(t, updated, result)
}

func TestStore_UpdateNonExistentKey(t *testing.T) {
	s := NewStore()
	key := "missing"

	success := s.Update(key, "value")
	assert.False(t, success)
}

func TestStore_Remove(t *testing.T) {
	s := NewStore()
	key := "delete_me"
	value := "gone"

	s.Set(key, value, 0)
	s.Remove(key)

	result := s.Get(key)
	assert.Nil(t, result)
}
