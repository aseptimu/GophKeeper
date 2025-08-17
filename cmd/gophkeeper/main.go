package main

import (
	"github.com/aseptimu/GophKeeper/internal/app/server/http"
	"log/slog"
)

func main() {
	slog.Info("Starting gophkeeper server")

	http.NewServer().Run()
}
