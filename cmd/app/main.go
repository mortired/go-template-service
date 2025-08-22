package main

import (
	"users/internal/bootstrap"
)

func main() {
	bootstrap.Setup().
		Run()
}
