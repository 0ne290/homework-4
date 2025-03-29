package internal

import "github.com/google/uuid"

type UuidProvider interface {
	Random() uuid.UUID
}

type RealUuidProvider struct{}

func NewRealUuidProvider() *RealUuidProvider {
	return &RealUuidProvider{}
}

func (*RealUuidProvider) Random() uuid.UUID {
	return uuid.New()
}
