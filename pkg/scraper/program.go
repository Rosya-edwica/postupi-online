package scraper

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/Rosya-edwica/postupi-online/pkg/common/models"
	"github.com/Rosya-edwica/postupi-online/pkg/scraper/html"
	"github.com/gocolly/colly"
)

func (s *Scraper) ScrapeAllPrograms(specUrl string) {
	specId := getId(specUrl)
	pagesCount := GetPagesCount(specUrl)
	for page := 1; page <= pagesCount; page++ {
		programsBlocks := html.FindHtmlBlocks(fmt.Sprintf("%s?page_num=%d", specUrl, page))
		var wg sync.WaitGroup
		wg.Add(len(programsBlocks))
		for _, item := range programsBlocks {
			go s.ScrapeProgram(item, specId, &wg)
		}
		wg.Wait()
	}
}

func (s *Scraper) ScrapeProgram(h *colly.HTMLElement, specId string, wg *sync.WaitGroup) {
	defer wg.Done()
	program := models.Program{}
	basic := scrapeBasic(h)
	if !hasBasicInfo(basic) {
		return
	}

	bodyHTMl, err := html.GetBody(basic.Url)
	checkErr(err)
	program.Base = basic
	program.VuzId = getVuzId(basic.Url)
	program.SpecId = specId
	program.ProgramId = getId(basic.Url)
	program.Description = html.GetSpecializationDescription(bodyHTMl)
	program.Base.Direction = html.GetProgramDirection(h, specId)
	program.Form = html.GetFormEducation(bodyHTMl)
	program.Exams = html.GetSubjects(bodyHTMl)
	program.HasProfessions = checkProfessionsExist(basic.Url)
	setNameToProgram(&program)

	programId, err := s.Db.SaveProgram(program)
	checkErr(err)
	s.ScrapeProfessions(programId, program.Base.Url)
}

func setNameToProgram(program *models.Program) {
	html.SetFullName(&program.Base)
	program.Base.Name = strings.Split(program.Base.Name, ":")[0]
	if program.Form == "Магистратура" { // Названия программ у магистратуры выглядят следующим образом: Профиль магистратуры "Управление свойствами нетканых материалов" РГУ им. А.Н. Косыгина, Москва
		re := regexp.MustCompile(`"(.*?)"`)
		program.Base.Name = re.FindString(program.Base.Name)
	}

}

func checkProfessionsExist(programUrl string) (exists bool) {
	body, err := html.GetBody(programUrl + "professii/")
	checkErr(err)
	body.ForEach("div.list-cover li", func(i int, h *colly.HTMLElement) {
		exists = true
	})
	return
}
