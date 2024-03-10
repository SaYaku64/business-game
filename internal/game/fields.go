package game

import (
	"math/rand"
	"time"
)

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

func (g *Game) getAnotherPlayer(playerID int) (plr *player) {
	if g.player1.playerID == playerID {
		plr = g.player2
	} else {
		plr = g.player1
	}

	return
}

func (g *Game) processUncommon(prop *property, playerID, diceSum int) (amount int) {
	if prop.isRailway() {
		return g.getRentRailway(playerID)
	}

	if prop.isUtility() {
		return g.getRentUtility(playerID, diceSum)
	}

	return
}

func (g *Game) getRentRailway(playerID int) (amount int) {
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

func (g *Game) getRentUtility(playerID, diceSum int) (amount int) {
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

func (g *Game) getPropertyByIndex(index int) (prop *property, found bool) {
	g.apMux.RLock()
	prop, found = g.allProperties[index]
	g.apMux.RUnlock()

	return
}

func (prop *property) belongsToPlayer(playerID int) (belongs bool) {
	return prop.ownerID == playerID
}

func (prop *property) isFree() (free bool) {
	return prop.ownerID == 0
}

func (prop *property) getAmountToPay() (amount int) {
	if prop.isMortgaged {
		return
	}

	if !prop.isWholeStreetBought {
		return prop.rent[0]
	}

	if prop.stars == 0 {
		return prop.rent[0] * 2
	}

	return prop.rent[prop.stars]
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

func (g *Game) processChance(plr *player) (text string, amount int, calculateNewField, moneyToPlayer bool) {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	iChance := r.Intn(len(chanceCards))

	chance := chanceCards[iChance]
	text = chance.text

	if chance.jailGuard {
		plr.jailGuard++

		return
	}

	if chance.repair {
		stars := 0

		g.apMux.RLock()
		for i := range g.allProperties {
			if g.allProperties[i].ownerID != plr.playerID {
				continue
			}

			if g.allProperties[i].stars == 5 {
				stars += 4 // hotel == 100, not 125
				continue
			}

			stars += g.allProperties[i].stars
		}
		g.apMux.RUnlock()

		amount -= (stars * 25)

		return
	}

	if chance.needToGo {
		calculateNewField = true

		if len(chance.nearest) != 0 {
			for _, v := range chance.nearest {
				if plr.fieldIndex < v {
					plr.fieldIndex = v

					return
				}
			}

			// next on other round
			plr.fieldIndex = chance.nearest[0]

			return
		}

		if chance.goBack != 0 {
			plr.fieldIndex -= chance.goBack

			return
		}

		if chance.canGetGO {
			if plr.fieldIndex > chance.goTo {
				plr.money += 200
			}
		}

		plr.fieldIndex = chance.goTo

		return
	}

	moneyToPlayer = chance.moneyToPlayer

	amount += chance.moneyFlow

	return
}

type chanceCard struct {
	text string

	moneyFlow     int
	moneyToPlayer bool // if true and moneyFlow + = get money from other; if true and moneyFlow - = give money to others

	jailGuard bool
	repair    bool

	needToGo bool
	canGetGO bool
	nearest  []int
	goTo     int
	goBack   int
}

var chanceCards = []chanceCard{
	{
		text:     "Відправляйтесь на \"Вперед\" (Отримайте $200)",
		needToGo: true,
		canGetGO: true,
		goTo:     0,
	},
	{
		text:     "Відправляйтесь на Іллінойс авеню. Якщо ви пройдете \"Вперед\", отримайте $200",
		needToGo: true,
		canGetGO: true,
		goTo:     24,
	},
	{
		text:     "Відправляйтесь на Сент-Чарльз авеню. Якщо ви пройдете \"Вперед\", отримайте $200",
		needToGo: true,
		canGetGO: true,
		goTo:     11,
	},
	{
		text:     "Прийдіть на найближчу комунальну службу. Якщо вона нікому не належить, ви можете купити її в банку.",
		needToGo: true,
		nearest:  []int{12, 28},
	},
	{
		text:     "Прийдіть на найближчу залізницю. Якщо вона нікому не належить, ви можете купити її в банку.",
		needToGo: true,
		nearest:  railwayList,
	},
	{
		text:      "Банк виплачує вам дивіденди у розмірі $50",
		moneyFlow: 50,
	},
	{
		text:      "Вийдіть з в'язниці безкоштовно. Цю картку можна використати в будь-який момент вашого ходу після ув'язнення.",
		jailGuard: true,
	},
	{
		text:     "Поверніться на 3 ходи назад",
		needToGo: true,
		goBack:   3,
	},
	{
		text:     "Вас запроторили до в'язниці. Не проходьте \"Вперед\", не отримуйте $200",
		needToGo: true,
		goTo:     10,
	},
	{
		text:   "Зробіть загальний ремонт всіх ваших будівель. За кожний будинок заплатіть $25, за кожний готель $100",
		repair: true,
	},
	{
		text:      "Заплатіть маленький податок у розмірі $15",
		moneyFlow: -15,
	},
	{
		text:     "Подорожуйте до залізниці «Рідінг». Якщо ви пройдете \"Вперед\", отримайте $200",
		needToGo: true,
		canGetGO: true,
		goTo:     5,
	},
	{
		text:     "Зробіть прогулянку по Бордволку",
		needToGo: true,
		goTo:     39,
	},
	{
		text:          "Ви обрані головою правління. Виплатіть кожному гравцеві $50",
		moneyFlow:     -50,
		moneyToPlayer: true,
	},
	{
		text:      "Кешбек з швидкого погашення вашого будівництва та кредиту. Отримайте $150",
		moneyFlow: 150,
	},
	{
		text:      "Ви виграли конкурс кросворду. Отримайте $100",
		moneyFlow: 100,
	},
	{
		text:     "Відправляйтесь на \"Вперед\" (Отримайте $200)",
		needToGo: true,
		canGetGO: true,
		goTo:     0,
	},
	{
		text:      "Помилка банку на вашу користь. Отримайте $200",
		moneyFlow: 200,
	},
	{
		text:      "Лікарські витрати. Сплатіть $50",
		moneyFlow: -50,
	},
	{
		text:      "За продаж акцій ви отримуєте $50",
		moneyFlow: 50,
	},
	{
		text:      "Вийдіть з в'язниці безкоштовно. Цю картку можна використати в будь-який момент вашого ходу після ув'язнення.",
		jailGuard: true,
	},
	{
		text:     "Вас запроторили до в'язниці. Не проходьте \"Вперед\", не отримуйте $200",
		needToGo: true,
		goTo:     10,
	},
	{
		text:          "Ніч великої опери. Зберіть $50 з кожного гравця за найкращі місця",
		moneyFlow:     50,
		moneyToPlayer: true,
	},
	{
		text:      "Кешбек з відпустки. Отримайте $100",
		moneyFlow: 100,
	},
	{
		text:      "Кешбеку податку на дохід. Отримайте $20",
		moneyFlow: 20,
	},
	{
		text:          "Сьогодні ваш день народження. Зберіть $10 від кожного гравця",
		moneyFlow:     10,
		moneyToPlayer: true,
	},
	{
		text:      "Завершується термін страхування на життя. Отримайте $100",
		moneyFlow: 100,
	},
	{
		text:      "Лікарські витрати. Сплатіть $100",
		moneyFlow: -100,
	},
	{
		text:      "Витрати на освіту. Сплатіть $150",
		moneyFlow: -150,
	},
	{
		text:      "Отримайте консультаційний гонорар у розмірі $25",
		moneyFlow: 25,
	},
	{
		text:   "Вас призначили відповідальним за ремонт вулиць. За кожний будинок заплатіть $25, за кожний готель $100",
		repair: true,
	},
	{
		text:      "Ви посіли друге місце в конкурсі краси. Отримайте $10",
		moneyFlow: 10,
	},
}
