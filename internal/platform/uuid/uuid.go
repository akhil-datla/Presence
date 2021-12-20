package uuid

import (
	"github.com/gofrs/uuid"
)

//New creates a string UUID
func New() string{
	return uuid.Must(uuid.NewV4()).String()
}