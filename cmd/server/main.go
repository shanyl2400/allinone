package main

import (
	"gomssbuilder/internal/api"
	"gomssbuilder/internal/repository/boltdb"
)

func main() {
	err := boltdb.GetClient().Open()
	if err != nil {
		panic(err)
	}
	defer boltdb.GetClient().Close()

	api.Start()
}
