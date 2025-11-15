package context

import (
	"context"
)

type key string

const (
	userIDKey key = "userID"
)

// SetUserID добавляет userID в контекст
func SetUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserID извлекает userID из контекста
func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDKey).(int64)
	return userID, ok
}

// MustGetUserID извлекает userID из контекста или паникует если не найден
// Использовать только в middleware где точно знаем что userID есть
func MustGetUserID(ctx context.Context) int64 {
	userID, ok := ctx.Value(userIDKey).(int64)
	if !ok {
		panic("userID not found in context")
	}
	return userID
}
