package repository

import (
	"context"
	"github.com/draco121/common/constants"
	"github.com/draco121/common/jwt"
	"github.com/draco121/common/models"
	"go.mongodb.org/mongo-driver/mongo"
)

var actionRoleMap = map[constants.Role][]constants.Action{
	constants.Tenant: {constants.Read, constants.Write},
	constants.Root:   {constants.All},
}

type IAuthorizationRepository interface {
	GetActions(role constants.Role) []constants.Action
	GetTokenClaims(token string) (*models.JwtCustomClaims, error)
	InsertAuthorizationLog(ctx context.Context, log models.AuthorizationLog) error
}

type authorizationRepository struct {
	IAuthorizationRepository
	db *mongo.Database
}

func NewAuthorizationRepo(database *mongo.Database) IAuthorizationRepository {
	return &authorizationRepository{
		db: database,
	}
}

func (c *authorizationRepository) GetActions(role constants.Role) []constants.Action {
	return actionRoleMap[role]
}

func (c *authorizationRepository) GetTokenClaims(token string) (*models.JwtCustomClaims, error) {
	claims, err := jwt.VerifyJwtToken(token)
	if err != nil {
		return nil, err
	} else {
		return &claims.JwtCustomClaims, nil
	}
}

func (c *authorizationRepository) InsertAuthorizationLog(ctx context.Context, log models.AuthorizationLog) error {
	_, err := c.db.Collection("authorization-log").InsertOne(ctx, log)
	return err
}
