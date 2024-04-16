package config

import "os"

var Port = "7540"
var DbFile = "data/scheduler.db"
var SchemaPath = "data/schema.sql"
var WebDir = "web"

func AcceptEnvironments() {
	updateVariable(&Port, "TODO_PORT")
	updateVariable(&DbFile, "TODO_DBFILE")
}

func updateVariable(variable *string, env string) {
	value := os.Getenv(env)
	if len(value) > 0 {
		*variable = value
	}
}
