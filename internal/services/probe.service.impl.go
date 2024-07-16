package services

import (
	"context"
)

type ProbeServiceImpl struct {
	ctx context.Context
}

func NewProbeService(ctx context.Context) ProbeService {
	return &ProbeServiceImpl{
		ctx: ctx,
	}
}
