package main

import (
	"fmt"
	"url-shortener/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	fmt.Print(cfg.Env)
}
