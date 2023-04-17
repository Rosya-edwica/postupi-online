package html

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Rosya-edwica/postupi-online/pkg/common/models"
	"github.com/gocolly/colly"
)

var Headers = map[string]string{
	"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
	"sec-ch-ua":    `Google Chrome";v="105", "Not)A;Brand";v="8", "Chromium";v="105"`,
	"accept":       "image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	"cookie":       "yandexuid=6850906421666216763; yabs-sid=1696581601666216766; yuidss=6850906421666216763; ymex=1981576766.yrts.1666216766#1981576766.yrtsi.1666216766; gdpr=0; _ym_uid=1666216766168837185; _ym_d=1666216766; yandex_login=rosya-8; i=Peh4utbtslQvge42D7cbDtH7CwXIiDs5Yp6IXWYsxx/SEQD1HtUncw/qqJV7NXqNqOS81fsaJSedcq/Ds9+yOfVKCNQ=; is_gdpr=0; skid=6879224341667473690; ys=udn.cDrQr9GA0L7RgdC70LDQsg%3D%3D#c_chck.841052032; is_gdpr_b=CIyaHxCclAE=; Session_id=3:1668355426.5.0.1666216795333:P19ouQ:2f.1.2:1|711384492.0.2|3:10261113.753043.lm80KKusrHll2DmXDLpHMjsmBYY; sessionid2=3:1668355426.5.0.1666216795333:P19ouQ:2f.1.2:1|711384492.0.2|3:10261113.753043.fakesign0000000000000000000; _ym_isad=1; _ym_visorc=b",
	"Content-Type": "text/html",
}

func SetFullName(basic *models.Basic) {
	body, err := GetBody(basic.Url)
	checkErr(err)

	fullname := body.ChildText("h1[id=prTitle]")
	if fullname != "" {
		basic.Name = fullname
	}
}

// func GetFullName(body *colly.HTMLElement) (name string) {
// name = body.ChildText("h1[id=prTitle]")
// return
// }

func GetCity(body *colly.HTMLElement) (city string) {
	city = body.ChildText("p.bg-nd__pre")
	return
}

func GetBody(url string) (body *colly.HTMLElement, err error) {
	c := colly.NewCollector()
	badGateway := CheckBadGateway(c)
	if badGateway {
		return
	}
	c.SetRequestTimeout(30 * time.Second)
	c.OnHTML("body", func(h *colly.HTMLElement) {
		body = h
	})
	err = c.Post(url, Headers)
	if body == nil {
		fmt.Println("Catched the error. Program stopped to sleep of 10 seconds. Url: ", url)
		time.Sleep(10 * time.Second)
		body, err = GetBody(url)
		checkErr(err)
		return
	}
	return
}

func GetFacts(body *colly.HTMLElement) (facts string) {
	var facts_list []string
	body.ForEach("ul.facts-list-nd li", func(i int, h *colly.HTMLElement) {
		facts_list = append(facts_list, h.Text)
	})
	facts = strings.Join(facts_list, ".")
	if facts != "" {
		return "\nФакты\n:" + facts
	}
	return
}

func FindHtmlBlocks(blocksUrl string) (blocks []*colly.HTMLElement) {
	c := colly.NewCollector()
	badGateway := CheckBadGateway(c)
	if badGateway {
		return
	}
	c.SetRequestTimeout(30 * time.Second)
	c.OnHTML("div.list-cover", func(h *colly.HTMLElement) {
		h.ForEach("li.list", func(i int, h *colly.HTMLElement) {
			blocks = append(blocks, h)
		})
	})
	err := c.Post(blocksUrl, Headers)
	checkErr(err)
	return
}

func GetFullDescription(body *colly.HTMLElement) (description string) {
	body.ForEach("div.list__info p", func(i int, e *colly.HTMLElement) {
		if i == 1 {
			description = e.Text
		}
	})
	return
}

func GetDescription(body *colly.HTMLElement) (description string) {
	description = GetFullDescription(body)
	if description == "" {
		return GetMiniDescription(body)
	}
	return
}

func GetMiniDescription(body *colly.HTMLElement) (description string) {
	description = body.ChildText("div.descr-min")
	return
}

func GetScores(body *colly.HTMLElement) (scores models.Scores) {
	body.ForEach("div.list__score-wrap p", func(i int, e *colly.HTMLElement) {
		switch true {
		case strings.Contains(e.Text, "бал.бюджет"):
			digit, err := strconv.ParseFloat(e.ChildText("b"), 64)
			checkErr(err)
			scores.PointsBudget = digit
		case strings.Contains(e.Text, "бал.платно"):
			digit, err := strconv.ParseFloat(e.ChildText("b"), 64)
			checkErr(err)
			scores.PointsPayment = digit
		case strings.Contains(e.Text, "бюджетных мест") && !strings.Contains("нет", e.ChildText("b")):
			scores.PlacesBudget = RemoveSpaceBetweenDigits(e.ChildText("b"))
		case strings.Contains(e.Text, "платных мест") && !strings.Contains("нет", e.ChildText("b")):
			scores.PlacesPayment = RemoveSpaceBetweenDigits(e.ChildText("b"))
		}
	})
	return
}

func GetSpecializationDirection(body *colly.HTMLElement, specId string) (direction string) {
	direction = strings.ReplaceAll(body.ChildText("p.list__pre"), specId, "")
	return
}
func GetProgramDirection(body *colly.HTMLElement, specId string) (direction string) {
	direction = strings.Split(body.ChildText("p.list__pre"), specId)[1]
	return
}

func GetSpecializationDescription(body *colly.HTMLElement) (description string) {
	description = body.ChildText("div.descr-max")
	return
}

func GetFormEducation(body *colly.HTMLElement) (form string) {
	re := regexp.MustCompile("Бакалавриат|Специалитет|Магистратура|Подготовка специалистов среднего звена|Подготовка квалифицированных рабочих (служащих)")
	form = re.FindString(body.ChildText("div.detail-box"))
	return
}

func GetSubjects(body *colly.HTMLElement) (subjects []string) {
	body.ForEach("div.score-box-wrap div.score-box", func(i int, h *colly.HTMLElement) {
		if i == 1 {
			h.ForEach("div.score-box__item", func(i int, h *colly.HTMLElement) {
				// score := h.ChildText("span.score-box__score") Не берем т.к без авторизации не показывается количество баллов
				exam := strings.Split(h.Text, "или")[0] //  Обрезаем строку и оставляем только матешу: Математика или другиеили Иностранный языкили Обществознание
				subjects = append(subjects, exam)
			})
		}
	})
	return
}

func CheckBadGateway(c *colly.Collector) (badGateway bool) {
	c.OnError(func(r *colly.Response, err error) {
		if r.StatusCode >= 500 {
			badGateway = true
			return
		}
	})
	return
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func RemoveSpaceBetweenDigits(text string) int {
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\u00a0", "")
	digit, err := strconv.Atoi(text)
	if err != nil {
		if err.Error() == `strconv.Atoi: parsing "": invalid syntax` {
			return 0
		}
	}
	return digit
}
