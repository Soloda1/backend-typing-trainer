package utils

import (
	"context"

	input "backend-typing-trainer/internal/domain/ports/input"
)

type actorContextKey string

var key actorContextKey = "actorContextKey"

func WithActor(ctx context.Context, actor input.Actor) context.Context {
	return context.WithValue(ctx, key, actor)
}

func ActorFromContext(ctx context.Context) (input.Actor, bool) {
	actor, ok := ctx.Value(key).(input.Actor)
	if !ok {
		return input.Actor{}, false
	}

	return actor, true
}
