package game

import (
	"fmt"
	"math/rand"
	"sync"

	a "github.com/SaYaku64/business-game/internal/alert"
	lp "github.com/SaYaku64/business-game/internal/lobby"
	"github.com/gin-gonic/gin"
)

type (
	GameModule struct {
		Games map[string]*GameState // key - lobbyID
		sync.RWMutex
	}

	GameState struct {
		LobbyID string

		Board map[int]*property // key - 0-39
		sync.RWMutex

		Players       []*Player // unused now, todo
		CurrentPlayer int       // index in Players slice

		player1 *Player // deprecated
		player2 *Player // deprecated
	}

	property struct {
		Index     int    // 0-39 (field)
		Buyable   bool   // if it can be bought
		Owner     string // SessionID of player which bought it
		OwnerName string // name of  player which bought it

		StreetIndex   int // 0 - brown ... 7 - blue
		uncommonIndex int // 0 - street, 1 - railway, 2 - utility

		Name  string // name to display
		Price int    // price of free property
		Rent  []int  // price of rent: 0 - clear rent, 1 - 1 star ... 5 - hotel

		ownerID int // 0 - none; 1 - player1; 2 - player2

		isWholeStreetBought bool // shows if whole street bought (x2 to rent without stars)
		stars               int  // houses and hotels on each property
		starPrice           int  // price to buy 1 star

		isMortgaged bool // shows if property is mortgaged
	}
)

func NewGameModule() *GameModule {
	return &GameModule{
		Games: make(map[string]*GameState),
	}
}

func (g *GameModule) GetGame(lobbyID string) (game *GameState, found bool) {
	g.RLock()
	defer g.RUnlock()

	game, found = g.Games[lobbyID]

	return
}

func (g *GameModule) SetGame(lobby lp.Lobby) {
	game := &GameState{
		LobbyID:       lobby.LobbyID,
		CurrentPlayer: rand.Intn(len(lobby.SessionIDs)),
		Players:       initPlayersFromLobby(lobby),
		Board:         initBoard(),
	}

	g.Lock()
	g.Games[game.LobbyID] = game
	g.Unlock()
}

func (g *GameModule) CheckActiveGame(lobbyID, playerName, sessionID string) (active, plrTurn bool, current int) {
	game, found := g.GetGame(lobbyID)
	if !found {
		return
	}

	for i, plr := range game.Players {
		if plr.SessionID == sessionID && plr.Name == playerName {
			active = true
			plrTurn = game.CurrentPlayer == i
		}
	}

	current = game.CurrentPlayer

	return
}

func (g *GameModule) UpdatePlates(lobbyID string) (players []*Player, ok bool) {
	game, found := g.GetGame(lobbyID)
	if !found {
		return
	}

	return game.Players, true
}

func (g *GameState) GetCurrentPlayer() (plr *Player) {
	return g.Players[g.CurrentPlayer]
}

func (g *GameState) NextPlayerTurn() (indexBefore int) {
	indexBefore = g.CurrentPlayer
	g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.Players)

	return
}

func (g *GameState) CalculateField(plr *Player, diceSum int) (canBuy bool, payToOther bool, moneyFlow int, recalculateField bool, chanceCard string) {
	if isFieldProperty(plr.Position) { // can be bought
		prop, found := g.getPropertyByIndex(plr.Position)
		if !found {
			a.Error.Println("game.go -> Roll -> getPropertyByIndex: property not found; index = ", plr.Position)
			return
		}

		if prop.isFree() {
			canBuy = true

			return
		}

		if prop.belongsToPlayer(plr.Index) {
			return
		}

		if prop.isMortgaged {
			return
		}

		if prop.isUncommon() {
			moneyFlow -= g.processUncommon(prop, plr.Index, diceSum)
			payToOther = true

			return
		}

		moneyFlow -= prop.getAmountToPay()
		payToOther = true

		return
	}

	if isTax, taxAmount := isFieldTax(plr.Position); isTax {
		moneyFlow -= taxAmount

		return
	}

	if isFieldCorner(plr.Position) {
		plr.processCornerJail()

		return
	}

	if isFieldChance(plr.Position) {
		chanceText, amount, calculateNewField, moneyToPlayer := g.processChance(plr)

		moneyFlow = amount
		payToOther = moneyToPlayer
		chanceCard = chanceText
		recalculateField = calculateNewField

		return
	}

	return
}

type RollResult struct {
	Empty bool

	ContinueAction bool
	DiceSum        int // sum of dice to continue
	WasDouble      bool
	JailPay50      bool

	RecalculateField bool // recalculate new field

	BuyFieldInfo *property // property info to buy
	CardText     string    // text on Chance Card

	MoneyFlow int // amount of money to pay
}

const (
	RollAction_Buy = iota + 1
	RollAction_PayRent
)

type RollActionResult struct {
	FirstDice  int `json:"firstDice"`
	SecondDice int `json:"secondDice"`
	Status     int `json:"status"`
}

func (g *GameState) RollAction(plr *Player) (result RollActionResult) {
	first, second := genDices()
	result.FirstDice = first
	result.SecondDice = second

	diceSum, _ /*double*/ := sumDices(first, second)

	// plr.CheckJail(double)

	// var wasInJail bool

	// if plr.Locked {
	// 	wasInJail = true

	// 	free, needToGo := plr.checkJailDice(diceSum, double)

	// 	if !free {
	// 		if needToGo {
	// 			result = &RollResult{
	// 				ContinueAction: true,
	// 				DiceSum:        diceSum,
	// 				JailPay50:      true,
	// 			}

	// 			return // pay 50$
	// 		}

	// 		result = &RollResult{
	// 			Empty: true,
	// 		}

	// 		return // sits until next round
	// 	}
	// }

	// if double && !wasInJail {
	// 	plr.DoubleCount++
	// 	if plr.checkDoubleToJail() {
	// 		plr.jail()

	// 		result = &RollResult{
	// 			Empty: true,
	// 		}

	// 		return // came to jail
	// 	}
	// } else {
	// 	plr.DoubleCount = 0
	// }

	plr.Position += diceSum

	if plr.Position >= len(g.Board) {
		plr.Balance += 200           // for the round
		plr.Position -= len(g.Board) // new round
	}

	field, _ := g.getPropertyByIndex(plr.Position)

	if field.Buyable {
		if field.isFree() {
			result.Status = RollAction_Buy

			return
			// // Send message to WebSocket connection
			// msg := fmt.Sprintf("%s landed on an unowned property. Would you like to buy %s for %d?", plr.Name, field.Name, field.Price)
			// err := wsConn.WriteMessage(websocket.TextMessage, []byte(msg))
			// if err != nil {
			// 	log.Println("Failed to write message to WebSocket:", err)
			// 	return
			// }
			// // Read response from WebSocket connection
			// input, err := readFromWebSocket()
			// if err != nil {
			// 	log.Println("Failed to read input from WebSocket:", err)
			// 	return
			// }
			// if input == "y" {
			// 	g.Players[player].Balance -= property.Price
			// 	property.Owner = g.Players[player].Name
			// }
		} else {
			result.Status = RollAction_PayRent

			return
			// // Send message to WebSocket connection
			// msg := fmt.Sprintf("%s landed on %s, owned by %s. Pay %d.", plr.Name, field.Name, field.Owner, field.Rent)
			// err := wsConn.WriteMessage(websocket.TextMessage, []byte(msg))
			// if err != nil {
			// 	log.Println("Failed to write message to WebSocket:", err)
			// 	return
			// }
			// g.Players[player].Balance -= property.Rent
			// g.Players[property.Owner].Balance += property.Rent
		}
	}

	// canBuy,
	// 	payToOther,
	// 	moneyFlow,
	// 	recalculateField,
	// 	chanceCard := g.CalculateField(plr, diceSum)

	// if canBuy {
	// 	prop, _ := g.getPropertyByIndex(plr.Position)
	// 	result = &RollResult{
	// 		BuyFieldInfo: prop,
	// 		WasDouble:    double,
	// 	}

	// 	return
	// }

	// result = &RollResult{
	// 	CardText:  chanceCard,
	// 	WasDouble: double,
	// }

	// if moneyFlow < 0 { // Player needs to pay
	// 	if plr.Balance+moneyFlow < 0 { // Player has less money than need to pay
	// 		result.MoneyFlow = moneyFlow
	// 		result.ContinueAction = true

	// 		return // pay money alert
	// 	}
	// }

	// if payToOther {
	// 	plr2 := g.getAnotherPlayer(plr.Index)
	// 	plr2.Balance -= moneyFlow
	// }

	// plr.Balance += moneyFlow

	// result.RecalculateField = recalculateField

	return
}

func (g *GameState) CalculateRollActionResult(plr *Player, result RollActionResult) gin.H {
	switch result.Status {
	case RollAction_Buy:
		field, _ := g.getPropertyByIndex(plr.Position)

		msg := fmt.Sprintf("<b>%s</b> став на поле <b>%s</b>. Його ціна %d.", plr.Name, field.Name, field.Price)

		return gin.H{"msg": msg, "result": result}
	case RollAction_PayRent:
		field, _ := g.getPropertyByIndex(plr.Position)

		pay := field.getAmountToPay()

		msg := fmt.Sprintf("<b>%s</b> став на поле <b>%s</b> гравця <b>%s</b>. Ціна аренди складає: <b>%d</b>.", plr.Name, field.Name, field.OwnerName, pay)

		return gin.H{"msg": msg, "result": result}
	}

	return gin.H{"result": result}
}

func (g *GameState) Buy(plr *Player) (answer gin.H, ok bool) {
	field, exists := g.getPropertyByIndex(plr.Position)
	if !exists {
		answer = gin.H{"error": "field does not exists"}

		return
	}

	if !field.Buyable || !field.isFree() {
		answer = gin.H{"error": "cannot buy this field"}

		return
	}

	if field.Price > plr.Balance {
		answer = gin.H{"error": "not enough money"}

		return
	}

	plr.Balance -= field.Price
	field.Owner = plr.SessionID
	ok = true

	msg := fmt.Sprintf("<b>%s</b> придбав <b>%s</b>.", plr.Name, field.Name)
	answer = gin.H{"msg": msg, "index": plr.Position, "plr": plr.Index}

	answer["balUpd"] = []gin.H{
		{
			"index":   plr.Index,
			"balance": plr.Balance,
		},
	}

	return
}

func (g *GameState) PayRent(plr *Player) (answer gin.H, ok bool) {
	field, exists := g.getPropertyByIndex(plr.Position)
	if !exists {
		answer = gin.H{"error": "field does not exists"}

		return
	}

	anotherPlr, found := g.getPlayerByID(field.Owner)
	if !found {
		answer = gin.H{"error": "unknown player to pay"}

		return
	}

	pay := field.getAmountToPay()
	if pay > plr.Balance {
		answer = gin.H{"error": "not enough money"}

		return
	}

	plr.Balance -= pay
	anotherPlr.Balance += pay
	ok = true

	msg := fmt.Sprintf("<b>%s</b> заплатив аренду <b>%s</b> у розмірі <b>%d</b>.", plr.Name, anotherPlr.Name, pay)
	answer = gin.H{"msg": msg}

	answer["balUpd"] = []gin.H{
		{
			"index":   plr.Index,
			"balance": plr.Balance,
		},
		{
			"index":   anotherPlr.Index,
			"balance": anotherPlr.Balance,
		},
	}

	return
}
