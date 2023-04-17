package scraper

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"

	"github.com/Rosya-edwica/postupi-online/pkg/common/database"
	"github.com/Rosya-edwica/postupi-online/pkg/common/logger"
	"github.com/Rosya-edwica/postupi-online/pkg/common/telegram"
	"github.com/Rosya-edwica/postupi-online/pkg/scraper/html"
)

type Scraper struct {
	Domain         string
	FormEducations []string
}

var db *database.DB

func Run(scr Scraper, db *database.DB) {

	pagesCount := GetPagesCount(scr.Domain)
	for i := 1; i <= pagesCount; i++ {
		url := fmt.Sprintf("%s?page_num=%d", scr.Domain, i)
		vuzesBlock := GetVuzesBlockInPage(url)
		for _, block := range vuzesBlock {
			vuz, err := scr.ScrapeVuz(block)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = db.SaveVuz(vuz)
			checkErr(err)
			scr.ScrapeAllSpecializations(vuz.Base.Url)
		}
	}
}

func GetVuzesBlockInPage(page string) (blocks []*colly.HTMLElement) {
	c := colly.NewCollector()
	c.SetRequestTimeout(30 * time.Second)

	c.OnHTML("div.list-cover li.list", func(h *colly.HTMLElement) {
		blocks = append(blocks, h)
	})

	c.OnError(func(r *colly.Response, err error) {
		if r.StatusCode >= 500 {
			logger.Log.Printf("[Ошибка] При подключении к %s произошла ошибка %s", page, err)
			return
		}
	})

	err := c.Post(page, html.Headers)
	if err != nil {
		fmt.Println("Catched the error. Program stopped to sleep of 10 seconds.")
		time.Sleep(10 * time.Second)
		err = c.Post(page, html.Headers)
	}
	return
}

func GetPagesCount(url string) (count int) {
	body, err := html.GetBody(url)
	checkErr(err)

	body.ForEach("a.paginator", func(i int, h *colly.HTMLElement) {
		count, _ = strconv.Atoi(h.Text)
	})
	if count == 0 {
		return 1
	}
	return
}

func getId(url string) (id string) {
	return strings.Split(url, "/")[len(strings.Split(url, "/"))-2]
}

func getVuzId(url string) (id string) {
	return strings.Split(url, "/")[4]
}

func checkErr(err error) {
	if err != nil {
		if strings.Contains(err.Error(), "Невозможно сохранить пустую") {
			return
		} else {
			telegram.Mailing(err.Error())
			logger.Log.Fatal(err)
		}
	}
}
