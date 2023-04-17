package scraper

import (
	"fmt"

	"github.com/Rosya-edwica/postupi-online/pkg/common/models"
	"github.com/Rosya-edwica/postupi-online/pkg/scraper/html"
	"github.com/gocolly/colly"
)

func scrapeBasic(h *colly.HTMLElement) (basic models.Basic) {
	if !BlockContainsImportantInfo(h) {
		fmt.Println("Нет важной инфы")
		return
	}
	basic.Description = html.GetFullDescription(h)
	basic.Scores = html.GetScores(h)
	basic.Url = h.ChildAttr("a", "href")
	basic.Image = h.ChildAttr("img.img-load", "data-dt")
	basic.Cost = html.RemoveSpaceBetweenDigits(h.ChildText("span.list__price b"))
	basic.Name = h.ChildText("h2.list__h")
	basic.Logo = h.ChildAttr("img.list__img-sm", "src")
	basic.Direction = h.ChildText("p.list__pre")
	return
}

func hasBasicInfo(basic models.Basic) bool {
	if basic == (models.Basic{}) {
		return false
	}
	return true
}

func BlockContainsImportantInfo(html *colly.HTMLElement) bool {
	if html.ChildText("p.list__score") == "" {
		return false
	}
	return true
}
