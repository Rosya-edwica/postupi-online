package scraper

import (
	"fmt"

	"github.com/gocolly/colly"

	"github.com/Rosya-edwica/postupi-online/pkg/common/models"
	"github.com/Rosya-edwica/postupi-online/pkg/scraper/html"
)

func (s *Scraper) ScrapeProfessions(connectionId int, programUrl string) {
	programId := getId(programUrl)
	pages := GetPagesCount(programUrl + "professii")
	for page := 1; page <= pages; page++ {
		url := fmt.Sprintf("%sprofessii/?page_num=%d", programUrl, page)
		s.ScrapePageProfessions(connectionId, url, programId)
	}
}

func (s *Scraper) ScrapePageProfessions(connectionId int, pageUrl string, programId string) {
	body, err := html.GetBody(pageUrl)
	checkErr(err)

	body.ForEach("li.list-col", func(i int, h *colly.HTMLElement) {
		profession := models.Profession{}
		profession.Name = h.ChildText("h2")
		profession.Image = h.ChildAttr("img.img-load", "data-dt")
		profession.ProgramId = programId
		db.SaveProfession(profession, connectionId)
	})
	return
}
