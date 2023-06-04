package scraper

import (
	"fmt"

	"github.com/Rosya-edwica/postupi-online/pkg/common/models"
	"github.com/Rosya-edwica/postupi-online/pkg/scraper/html"
	"github.com/gocolly/colly"
)

func (s *Scraper) ScrapeAllSpecializations(vuzUrl string) {
	for _, form := range s.FormEducations {
		url := vuzUrl + form
		pagesCount := GetPagesCount(url)
		for page := 1; page <= pagesCount; page++ {
			pageUrl := fmt.Sprintf("%s?page_num=%d", url, page)
			specsBlocks := html.FindHtmlBlocks(pageUrl)
			for _, spec := range specsBlocks {
				s.ScrapeOneSpecialization(spec)
			}
		}
	}
}

func (s *Scraper) ScrapeOneSpecialization(HTMLBlock *colly.HTMLElement) {
	specialization := models.Specialization{}
	basic := scrapeBasic(HTMLBlock)
	if !hasBasicInfo(basic) {
		return
	}
	bodyHTML, err := html.GetBody(basic.Url)
	checkErr(err)

	specialization.Base = basic
	specialization.SpecId = getId(basic.Url)
	specialization.VuzId = getVuzId(basic.Url)
	specialization.Base.Direction = html.GetSpecializationDirection(HTMLBlock, specialization.SpecId)
	specialization.Description = html.GetSpecializationDescription(bodyHTML)

	s.Db.SaveSpecialization(specialization)
	s.ScrapeAllPrograms(basic.Url)
}
