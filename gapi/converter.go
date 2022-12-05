package gapi

import (
	db "github.com/bagashiz/Simple-Bank/db/sqlc"
	"github.com/bagashiz/Simple-Bank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// convertUser converts a database user field into a protobuf user field
func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
