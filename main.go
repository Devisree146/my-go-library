package main

import (
	"unified/multicache"
)

func main() {
	router := multicache.SetupRouter()
	router.Run(":8080")
}
