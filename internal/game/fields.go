package game

import "fmt"

var notPropertyList = []int{0, 2, 4, 7 /**/, 10, 17 /**/, 20, 22 /**/, 30, 33, 36, 38}
var chanceList = []int{2, 7, 17, 22, 33, 36}
var railwayList = []int{5, 15, 25, 35}
var cornerList = []int{0, 10, 20, 30}
var utilityList = []int{12, 28}

// var taxList = []int{4, 38}

func isFieldProperty(index int) (found bool) {
	for _, i := range notPropertyList {
		if index == i {
			return false
		}
	}

	return true
}

func isFieldChance(index int) (found bool) {
	for _, i := range chanceList {
		if index == i {
			return true
		}
	}

	return false
}

func isFieldCorner(index int) (found bool) {
	for _, i := range cornerList {
		if index == i {
			return true
		}
	}

	return false
}

func isFieldTax(index int) (found bool, taxAmount int) {
	switch index {
	case 4:
		return true, 200
	case 38:
		return true, 100
	default:
		return
	}
}

func (g *GameState) getAnotherPlayer(playerID int) (plr *Player) {
	if g.player1.Index == playerID {
		plr = g.player2
	} else {
		plr = g.player1
	}

	return
}

func (g *GameState) processUncommon(prop *property, playerID, diceSum int) (amount int) {
	if prop.isRailway() {
		return g.getRentRailway(playerID)
	}

	if prop.isUtility() {
		return g.getRentUtility(playerID, diceSum)
	}

	return
}

func (g *GameState) getRentRailway(playerID int) (amount int) {
	multiplier := 0
	for _, v := range railwayList {
		railway, _ := g.getPropertyByIndex(v)
		if railway.ownerID != playerID && railway.ownerID != 0 {
			multiplier++
		}
	}

	switch multiplier {
	case 0, 1, 2:
		return multiplier * 25
	case 3:
		return 100
	case 4:
		return 200
	}

	return
}

func (g *GameState) getRentUtility(playerID, diceSum int) (amount int) {
	multiplier := 0
	for _, v := range utilityList {
		utility, _ := g.getPropertyByIndex(v)
		if utility.ownerID != playerID && utility.ownerID != 0 {
			multiplier++
		}
	}

	switch multiplier {
	case 1:
		return diceSum * 4
	case 2:
		return diceSum * 10
	}

	return
}

func (g *GameState) getPropertyByIndex(index int) (prop *property, found bool) {
	g.RLock()
	prop, found = g.Board[index]
	g.RUnlock()

	return
}

func initBoard() (board map[int]*property) {
	board = make(map[int]*property, 40)

	for i := 0; i < 40; i++ {
		board[i] = &property{
			Index:       i,
			Buyable:     true,
			StreetIndex: 0,
			Name:        fmt.Sprintf("Street-%d", i),
			Price:       10,
			Rent:        []int{1, 3, 5, 7, 9, 11},
		}
	}

	return
}

func (prop *property) belongsToPlayer(playerID int) (belongs bool) {
	return prop.ownerID == playerID
}

func (prop *property) isFree() (free bool) {
	return prop.Owner == ""
}

func (prop *property) getAmountToPay() (amount int) {
	if len(prop.Rent) == 0 {
		return 0
	}

	if prop.isMortgaged {
		return
	}

	if !prop.isWholeStreetBought {
		return prop.Rent[0]
	}

	if prop.stars == 0 {
		return prop.Rent[0] * 2
	}

	return prop.Rent[prop.stars]
}

func (prop *property) isUncommon() (uncommon bool) {
	return prop.uncommonIndex != 0
}

func (prop *property) isRailway() (railway bool) {
	return prop.uncommonIndex == 1
}

func (prop *property) isUtility() (utility bool) {
	return prop.uncommonIndex == 2
}
