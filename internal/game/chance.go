package game

import (
	"math/rand"
	"time"
)

func (g *GameState) processChance(plr *Player) (text string, amount int, calculateNewField, moneyToPlayer bool) {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	iChance := r.Intn(len(chanceCards))

	chance := chanceCards[iChance]
	text = chance.text

	if chance.jailGuard {
		plr.CardExit++

		return
	}

	if chance.repair {
		stars := 0

		g.RLock()
		for i := range g.Board {
			if g.Board[i].ownerID != plr.Index {
				continue
			}

			if g.Board[i].stars == 5 {
				stars += 4 // hotel == 100, not 125
				continue
			}

			stars += g.Board[i].stars
		}
		g.RUnlock()

		amount -= (stars * 25)

		return
	}

	if chance.needToGo {
		calculateNewField = true

		if len(chance.nearest) != 0 {
			for _, v := range chance.nearest {
				if plr.Position < v {
					plr.Position = v

					return
				}
			}

			// next on other round
			plr.Position = chance.nearest[0]

			return
		}

		if chance.goBack != 0 {
			plr.Position -= chance.goBack

			return
		}

		if chance.canGetGO {
			if plr.Position > chance.goTo {
				plr.Balance += 200
			}
		}

		plr.Position = chance.goTo

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
