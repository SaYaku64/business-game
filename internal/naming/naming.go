package naming

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetRandName(c *gin.Context) {
	c.String(http.StatusOK, genName())
}

func genName() string {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	iAd := r.Intn(len(adjectives))
	iAn := r.Intn(len(animals))

	return fmt.Sprintf("%s %s", adjectives[iAd], animals[iAn])
}

var adjectives = []string{
	"Фантастичний",
	"Блискучий",
	"Мальовничий",
	"Дивовижний",
	"Магічний",
	"Екзотичний",
	"Чудовий",
	"Барвистий",
	"Чарівний",
	"Вишуканий",
	"Прекрасний",
	"Зачарований",
	"Феєричний",
	"Оригінальний",
	"Вражаючий",
	"Виразний",
	"Унікальний",
	"Розкішний",
	"Елегантний",
	"Романтичний",
	"Повітряний",
	"Легкий",
	"Спокусливий",
	"Сяючий",
	"Витончений",
	"Загадковий",
	"Драйвовий",
	"Безпечний",
	"Веселий",
	"Забавний",
	"Живописний",
	"Сюрреалістичний",
	"Епічний",
	"Атмосферний",
	"Весняний",
	"Осяйний",
	"Захоплюючий",
	"Свіжий",
	"Жвавий",
	"Буйний",
	"Відчайдушний",
	"Драйвовий",
	"Витівливий",
	"Палкий",
	"Веселий",
	"Безтурботний",
	"Чуттєвий",
	"Шалений",
	"Розумний",
}

var animals = []string{
	"Собака",
	"Корова",
	"Кіт",
	"Кінь",
	"Осел",
	"Тигр",
	"Лев",
	"Пантера",
	"Леопард",
	"Гепард",
	"Ведмідь",
	"Слон",
	"Білий ведмідь",
	"Черепаха",
	"Земноводна черепаха",
	"Крокодил",
	"Кролик",
	"Дикобраз",
	"Заєць",
	"Курка",
	"Голуб",
	"Альбатрос",
	"Ворона",
	"Риба",
	"Дельфін",
	"Жаба",
	"Кит",
	"Алігатор",
	"Орел",
	"Летюча білка",
	"Страус",
	"Лисиця",
	"Коза",
	"Шакал",
	"Ему",
	"Броненосець",
	"Угорь",
	"Гусак",
	"Полярна лисиця",
	"Вовк",
	"Бігль",
	"Горила",
	"Шимпанзе",
	"Мавпа",
	"Бобер",
	"Орангутанг",
	"Антилопа",
	"Кажан",
	"Барсук",
	"Жирафа",
	"Краб-ерміт",
	"Великий панда",
	"Ховрах",
	"Кобра",
	"Молотоголовий акула",
	"Верблюд",
	"Яструб",
	"Олень",
	"Хамелеон",
	"Гіпопотам",
	"Ягуар",
	"Чихуахуа",
	"Королівська кобра",
	"Ібекс",
	"Ящірка",
	"Коала",
	"Кенгуру",
	"Ігуана",
	"Лама",
	"Шиншила",
	"Додо",
	"Медуза",
	"Носоріг",
	"Зебра",
	"Опосум",
	"Вомбат",
	"Бізон",
	"Бик",
	"Буйвол",
	"Вівця",
	"Сурикат",
	"Миша",
	"Видра",
	"Лінивець",
	"Сова",
	"Ловець",
	"Фламінго",
	"Єнот",
	"Кріт",
	"Качка",
	"Лебідь",
	"Рись",
	"Монітор",
	"Лось",
	"Кабан",
	"Лемур",
	"Мул",
	"Бавійка",
	"Мамонт",
	"Синій кит",
	"Миша",
	"Змія",
	"Павич",
}
