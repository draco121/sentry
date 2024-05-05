package core

import (
	"context"
	"github.com/draco121/horizon/constants"
	"github.com/draco121/horizon/models"
	"github.com/draco121/horizon/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"sentry/repository"
	"time"
)

type IAuthorizationService interface {
	Authorize(ctx context.Context, authorizationInput *models.AuthorizationInput) *models.AuthorizationOutput
}

type authorizationService struct {
	IAuthorizationService
	repo   repository.IAuthorizationRepository
	client *mongo.Client
}

func NewAuthorizationService(client *mongo.Client, repository repository.IAuthorizationRepository) IAuthorizationService {
	return &authorizationService{
		repo:   repository,
		client: client,
	}
}

func (c *authorizationService) Authorize(ctx context.Context, authorizationInput *models.AuthorizationInput) *models.AuthorizationOutput {
	session, err := c.client.StartSession()
	if err != nil {
		utils.Logger.Error("failed to start mongo session", "error: ", err.Error())
		return nil
	}
	defer session.EndSession(ctx)
	err = session.StartTransaction()
	if err != nil {
		utils.Logger.Error("failed to start mongo transaction", "error: ", err.Error())
		return nil
	}
	claims, err := c.repo.GetTokenClaims(authorizationInput.Token)
	if err != nil {
		log := models.AuthorizationLog{
			Grant:       constants.Rejected,
			RequestedAt: time.Now(),
			Actions:     authorizationInput.Actions,
			Reason:      err.Error(),
		}
		err = c.repo.InsertAuthorizationLog(ctx, log)
		if err != nil {
			utils.Logger.Error("failed to insert authorization log", "error: ", err.Error())
			return nil
		}
		utils.Logger.Info("inserted authorization log")
		utils.Logger.Info("authorization failed ", log)
		_ = session.CommitTransaction(ctx)
		return &models.AuthorizationOutput{
			Grant:  constants.Rejected,
			UserId: primitive.NilObjectID,
		}
	} else {
		allowedActions := c.repo.GetActions(claims.Role)
		if authorizationEngine(allowedActions, authorizationInput.Actions) {
			log := models.AuthorizationLog{
				Grant:       constants.Allowed,
				RequestedAt: time.Now(),
				Actions:     authorizationInput.Actions,
				UserId:      claims.UserId,
				Role:        claims.Role,
				Reason:      "permissions matched",
			}
			err = c.repo.InsertAuthorizationLog(ctx, log)
			if err != nil {
				utils.Logger.Error("failed to insert authorization log", "error: ", err.Error())
				return nil
			}
			utils.Logger.Info("inserted authorization log")
			utils.Logger.Info("authorization completed ", log)
			_ = session.CommitTransaction(ctx)
			return &models.AuthorizationOutput{
				Grant:  constants.Allowed,
				UserId: claims.UserId,
			}
		} else {
			log := models.AuthorizationLog{
				Grant:       constants.Rejected,
				RequestedAt: time.Now(),
				Actions:     authorizationInput.Actions,
				UserId:      claims.UserId,
				Role:        claims.Role,
				Reason:      "permissions does not match",
			}
			err = c.repo.InsertAuthorizationLog(ctx, log)
			if err != nil {
				utils.Logger.Error("failed to insert authorization log", "error: ", err.Error())
				return nil
			}
			utils.Logger.Info("inserted authorization log")
			utils.Logger.Info("authorization failed ", log)
			_ = session.CommitTransaction(ctx)
			return &models.AuthorizationOutput{
				Grant:  constants.Rejected,
				UserId: claims.UserId,
			}
		}
	}
}
