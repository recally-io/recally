package handlers

import (
	"context"
	"fmt"
	"strings"
	"vibrain/internal/pkg/constant"
)

func (h *Handler) getOrCreateUser(ctx context.Context) (*User, error) {
	userID := ctx.Value(constant.ContextKey(constant.ContextKeyUserID)).(string)

	user, err := h.repository.GetUser(ctx, userID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			userName := ctx.Value(constant.ContextKey(constant.ContextKeyUserName)).(string)
			user, err = h.repository.CreateUser(ctx, userName, userID)
			if err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
	}
	return user, nil
}
