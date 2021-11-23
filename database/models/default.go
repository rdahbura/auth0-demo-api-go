package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Login struct {
	Email    string `bson:"email,omitempty" json:"email,omitempty" validate:"required,email"`
	Password string `bson:"password,omitempty" json:"password,omitempty" validate:"required"`
}

type User struct {
	Id            primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Email         string             `bson:"email,omitempty" json:"email,omitempty"`
	EmailVerified *bool              `bson:"email_verified,omitempty" json:"email_verified,omitempty"`
	Username      string             `bson:"username,omitempty" json:"username,omitempty"`
	Password      string             `bson:"password,omitempty" json:"password,omitempty"`
	PasswordHash  string             `bson:"password_hash,omitempty" json:"password_hash,omitempty"`
	FamilyName    string             `bson:"family_name,omitempty" json:"family_name,omitempty"`
	GivenName     string             `bson:"given_name,omitempty" json:"given_name,omitempty"`
	CreatedAt     time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt     time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
