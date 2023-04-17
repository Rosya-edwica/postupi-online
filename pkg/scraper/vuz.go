package scraper

import (
	"github.com/Rosya-edwica/postupi-online/pkg/common/logger"
	"github.com/Rosya-edwica/postupi-online/pkg/common/models"
	"github.com/Rosya-edwica/postupi-online/pkg/scraper/html"
	"github.com/gocolly/colly"
)

func (s *Scraper) ScrapeVuz(h *colly.HTMLElement) (vuz models.Vuz, err error) {
	basic := scrapeBasic(h)
	if !hasBasicInfo(basic) {
		logger.Log.Printf("[Ошибка] Нет базовой инфы")
		return
	}
	bodyHTML, err := html.GetBody(basic.Url)
	checkErr(err)
	vuz.Base = basic
	vuz.VuzId = getId(vuz.Base.Url)
	vuz.City = html.GetCity(bodyHTML)
	vuz.Description = html.GetDescription(bodyHTML) + html.GetFacts(bodyHTML)
	html.SetFullName(&basic)
	return
}
