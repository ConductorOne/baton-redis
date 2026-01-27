package main

import (
	cfg "github.com/conductorone/baton-redis/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("redis", cfg.Config)
}
