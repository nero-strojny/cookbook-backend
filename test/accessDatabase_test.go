package test

import (
	"flag"
	"server/controller"
	"testing"
)

var dbPointer = flag.String("DB_STRING", "", "Database connection string")
var envPointer = flag.String("ENV", "", "Environment string")

func TestDatabaseSetupForCalorieLog(t *testing.T) {
	controller.SetClients(*dbPointer, *envPointer)
}
