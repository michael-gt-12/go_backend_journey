package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Allowed roles.
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// User is the stored user document.
// Password holds ONLY the bcrypt hash. json:"-" guarantees it is never
// serialised into an API response, even if a handler returns the struct by mistake.
type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"-" bson:"password"`
	Role     string             `json:"role" bson:"role"`
}

// Credentials is the request body for /register and /login.
// It is separate from User so the plaintext password never lives in the stored model.
// bcrypt ignores bytes after 72, so we cap the length instead of silently truncating.
type Credentials struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}
