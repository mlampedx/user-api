package main

import (
	"os"
)

func main() {
	a := API{}
	a.Initialize(
		os.Getenv("LOCAL_PSQL_DB_USERNAME"),
		os.Getenv("LOCAL_PSQL_DB_PASSWORD"),
		os.Getenv("LOCAL_USER_API_PSQL_DB_NAME"))

	a.Run(":8080")
}
