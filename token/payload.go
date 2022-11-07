package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Different types of error returned by verifyToken function
var (
  ErrExpiredToken = errors.New("token has expired")
  ErrInvalidToken = errors.New("token is invalid")
)

// Payload contains the payload data of the token
type Payload struct {
  ID uuid.UUID `json:"id"`
  Username string `json:"username"`
  IsUsedAt time.Time `json:"isused_at"`
  ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error)  {
  tokenID, err := uuid.NewRandom()
  if err != nil {
    return nil, err
  }

  payload := &Payload{
    ID: tokenID,
    Username: username,
    IsUsedAt: time.Now(),
    ExpiredAt: time.Now().Add(duration),
  }

  return payload, nil
}

// Valid is a function that checks if the token payload is valid or not
func (payload *Payload) Valid() error {
  if time.Now().After(payload.ExpiredAt) {
    return ErrExpiredToken
  }

  return nil
}
