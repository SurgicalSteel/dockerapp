package token

import "github.com/satori/go.uuid"

func GenerateToken(name string) string {
	u := uuid.NewV5(uuid.UUID{}, name)
	return u.String()
}
