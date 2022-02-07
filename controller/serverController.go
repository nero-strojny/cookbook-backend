package controller

import (
	"context"
	"fmt"
	"server/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type ServerControl interface {
	HealthCheck() models.HealthStatus
}

type ServerController struct {
	client *mongo.Client
}

func NewServerController(client *mongo.Client) ServerController {
	return ServerController{client}
}

func (sc ServerController) HealthCheck() models.HealthStatus {
	mongoStatus := "OK"
	err := sc.client.Ping(context.TODO(), nil)
	if err != nil {
		mongoStatus = fmt.Sprintf("Mongo ERROR: %s", err)
	}
	return models.HealthStatus{DB: mongoStatus}
}
