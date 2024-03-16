package game

import (
	"sync"

	a "github.com/SaYaku64/business-game/internal/alert"
)

type (
	GameModule struct {
		Games map[string]*Game // key - lobbyID
		gMux  sync.RWMutex
	}

	Game struct {
		lobbyID string

		// freeProperties map[int]property
		// fpMux          sync.RWMutex

		allProperties map[int]*property
		apMux         sync.RWMutex

		player1 *player
		player2 *player
	}

	player struct {
		sessionID  string
		playerName string

		fieldIndex int // place, where he is now

		jailGuard int // card from chance
		jailed    bool
		jailDay   int // if 3 = needs to pay 50$

		playerID int
		money    int
		// properties map[int]*property
		// prMux      sync.RWMutex

		doubleCount int // if 3 = go to jail
	}

	property struct {
		ownerID int // 0 - none; 1 - player1; 2 - player2

		price int   // price of free property
		rent  []int // price of rent: 0 - clear rent, 1 - 1 star ... 5 - hotel

		index         int // 0-39 (field)
		uncommonIndex int // 0 - street, 1 - railway, 2 - utility
		streetIndex   int // 0 - brown ... 7 - blue

		isWholeStreetBought bool // shows if whole street bought (x2 to rent without stars)
		stars               int  // houses and hotels on each property
		starPrice           int  // price to buy 1 star

		isMortgaged bool // shows if property is mortgaged
	}
)

func (plr *player) checkDoubleToJail() bool {
	return plr.doubleCount >= 3
}

func (plr *player) jail() {
	plr.fieldIndex = 10
	plr.jailed = true
	plr.jailDay = 0
}

func (plr *player) unJail() {
	plr.doubleCount = 0
	plr.jailDay = 0
	plr.jailed = false
}

func (plr *player) processCornerJail() {
	if plr.fieldIndex == 30 {
		plr.jail()
	}
}

func (plr *player) checkJailDice(diceSum int, double bool) (free, needToGo bool) {
	plr.jailDay++

	if double {
		plr.unJail()

		return true, true
	}

	if plr.jailDay < 3 {
		return
	}

	if plr.jailGuard > 0 {
		plr.jailGuard--
		plr.unJail()

		return true, true
	}

	return false, true
}

func (g *Game) CalculateField(plr *player, diceSum int) (canBuy bool, payToOther bool, moneyFlow int, recalculateField bool, chanceCard string) {
	if isFieldProperty(plr.fieldIndex) { // can be bought
		prop, found := g.getPropertyByIndex(plr.fieldIndex)
		if !found {
			a.Error.Println("game.go -> Roll -> getPropertyByIndex: property not found; index = ", plr.fieldIndex)
			return
		}

		if prop.isFree() {
			canBuy = true

			return
		}

		if prop.belongsToPlayer(plr.playerID) {
			return
		}

		if prop.isMortgaged {
			return
		}

		if prop.isUncommon() {
			moneyFlow -= g.processUncommon(prop, plr.playerID, diceSum)
			payToOther = true

			return
		}

		moneyFlow -= prop.getAmountToPay()
		payToOther = true

		return
	}

	if isTax, taxAmount := isFieldTax(plr.fieldIndex); isTax {
		moneyFlow -= taxAmount

		return
	}

	if isFieldCorner(plr.fieldIndex) {
		plr.processCornerJail()

		return
	}

	if isFieldChance(plr.fieldIndex) {
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

func (g *Game) RollAction(plr *player) (result *RollResult) {
	diceSum, double := sumDices()

	var wasInJail bool

	if plr.jailed {
		wasInJail = true

		free, needToGo := plr.checkJailDice(diceSum, double)

		if !free {
			if needToGo {
				result = &RollResult{
					ContinueAction: true,
					DiceSum:        diceSum,
					JailPay50:      true,
				}

				return // pay 50$
			}

			result = &RollResult{
				Empty: true,
			}

			return // sits until next round
		}
	}

	if double && !wasInJail {
		plr.doubleCount++
		if plr.checkDoubleToJail() {
			plr.jail()

			result = &RollResult{
				Empty: true,
			}

			return // came to jail
		}
	} else {
		plr.doubleCount = 0
	}

	plr.fieldIndex += diceSum

	if plr.fieldIndex >= 40 {
		plr.money += 200     // for the round
		plr.fieldIndex -= 40 // new round
	}

	canBuy,
		payToOther,
		moneyFlow,
		recalculateField,
		chanceCard := g.CalculateField(plr, diceSum)

	if canBuy {
		prop, _ := g.getPropertyByIndex(plr.fieldIndex)
		result = &RollResult{
			BuyFieldInfo: prop,
			WasDouble:    double,
		}

		return
	}

	result = &RollResult{
		CardText:  chanceCard,
		WasDouble: double,
	}

	if moneyFlow < 0 { // player needs to pay
		if plr.money+moneyFlow < 0 { // player has less money than need to pay
			result.MoneyFlow = moneyFlow
			result.ContinueAction = true

			return // pay money alert
		}
	}

	if payToOther {
		plr2 := g.getAnotherPlayer(plr.playerID)
		plr2.money -= moneyFlow
	}

	plr.money += moneyFlow

	result.RecalculateField = recalculateField

	return
}
