package main

import (
	a "github.com/SaYaku64/business-game/internal/alert"
	"github.com/SaYaku64/business-game/internal/game"
	"github.com/SaYaku64/business-game/internal/lobby"
	"github.com/SaYaku64/business-game/internal/router"
)

func main() {
	lobbyModule := lobby.CreateLobbyModule()

	gameModule := game.NewGameModule()

	r := router.NewRouter(lobbyModule, gameModule)
	r.Load()

	if err := r.RunRouter(); err != nil {
		a.Warning.Println(err)
	}
}
