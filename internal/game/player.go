package game

import (
	lp "github.com/SaYaku64/business-game/internal/lobby"
)

type (
	Player struct {
		SessionID string // unique identifier
		Index     int    // index in current game
		Name      string // name to display

		Position int // place, where he is now
		Balance  int

		DoubleCount int // if 3 = go to jail

		Jail
	}

	Jail struct {
		CardExit int // card from chance
		Day      int // if 3 = needs to pay 50$
		Locked   bool
	}
)

func (g *GameState) getPlayerByID(sessionID string) (plr *Player, found bool) {
	for i := range g.Players {
		if g.Players[i].SessionID == sessionID {
			return g.Players[i], true
		}
	}

	return
}

func initPlayersFromLobby(lobby lp.Lobby) (players []*Player) {
	players = make([]*Player, len(lobby.SessionIDs))

	for i := range lobby.SessionIDs {
		players[i] = &Player{
			SessionID: lobby.SessionIDs[i],
			Index:     i,
			Name:      lobby.PlayerNames[i],
			Balance:   1500,
		}
	}

	return
}

func (plr *Player) checkDoubleToJail() bool {
	return plr.DoubleCount >= 3
}

func (plr *Player) jail() {
	plr.Position = 10
	plr.Locked = true
	plr.Day = 0
}

func (plr *Player) unJail() {
	plr.DoubleCount = 0
	plr.Day = 0
	plr.Locked = false
}

func (plr *Player) processCornerJail() {
	if plr.Position == 30 {
		plr.jail()
	}
}

func (plr *Player) checkJailDice(diceSum int, double bool) (free, needToGo bool) {
	plr.Day++

	if double {
		plr.unJail()

		return true, true
	}

	if plr.Day < 3 {
		return
	}

	if plr.CardExit > 0 {
		plr.CardExit--
		plr.unJail()

		return true, true
	}

	return false, true
}
