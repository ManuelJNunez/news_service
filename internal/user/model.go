package user

import "go.mongodb.org/mongo-driver/v2/bson"

// User stored in database
type User struct {
	ID       bson.ObjectID `bson:"_id,omitempty"`
	Username string        `bson:"username"`
	Password string        `bson:"password"`
}

// Data received from the login form
type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserOutput for API responses without confidential data
type UserOutput struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// Convert User to UserOutput
func (u *User) ToOutput() UserOutput {
	return UserOutput{
		ID:       u.ID.Hex(),
		Username: u.Username,
	}
}
