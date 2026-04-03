package utils

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"backend-typing-trainer/internal/domain/models"
	input "backend-typing-trainer/internal/domain/ports/input"
)

func TestWithActorAndActorFromContext(t *testing.T) {
	actor := input.Actor{
		UserID: uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		Role:   models.UserRoleUser,
	}

	ctx := WithActor(context.Background(), actor)

	got, ok := ActorFromContext(ctx)
	require.True(t, ok)
	require.Equal(t, actor, got)
}

func TestActorFromContext_MissingActor(t *testing.T) {
	got, ok := ActorFromContext(context.Background())
	require.False(t, ok)
	require.Equal(t, input.Actor{}, got)
}
