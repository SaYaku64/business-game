package main

import (
	"github.com/SaYaku64/business-game/internal/alert"
	"github.com/SaYaku64/business-game/internal/lobby"
	"github.com/SaYaku64/business-game/internal/router"
)

func main() {
	lobbyModule := lobby.CreateLobbyModule()

	r := router.NewRouter(lobbyModule)
	r.Load()

	if err := r.RunRouter(); err != nil {
		alert.Warning(err)
	}
}
