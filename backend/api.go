package main

import (
	_ "backend/handlers"
	"google.golang.org/appengine/v2"
)

func main() {
	appengine.Main()
}
