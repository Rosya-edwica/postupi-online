package scraper

import (
	"fmt"
	"log"
	"time"

	"github.com/Rosya-edwica/postupi-online/pkg/common/models"
	"github.com/Rosya-edwica/postupi-online/pkg/scraper/html"
	"github.com/gocolly/colly"
)

func scrapeContacts(url string) (err error) {
	contact := models.Contacts{}
	c := colly.NewCollector()
	badGateway := html.CheckBadGateway(c)
	if badGateway {
		return
	}
	c.SetRequestTimeout(30 * time.Second)

	// Ищем 4 span блока с контактами вуза
	c.OnHTML("section.section-box", func(h *colly.HTMLElement) {
		contactList := [4]string{}
		h.ForEach("span", func(i int, e *colly.HTMLElement) {
			contactList[i] = e.Text
		})

		if contactList[0] != "" {
			contact.WebSite = contactList[0]
			contact.Email = contactList[1]
			contact.Phone = contactList[2]
			contact.Address = contactList[3]
			contact.VuzId = getVuzId(url)
		}
	})
	err = c.Post(url, html.Headers)
	if err != nil {
		fmt.Println("Catched the error. Program stopped to sleep of 10 seconds.")
		time.Sleep(10 * time.Second)
		err = scrapeContacts(url + "contacts/")
		checkErr(err)
	}
	db.SaveContacts(contact)
	checkErr(err)
	log.Printf("Contact:%s", contact.VuzId)
	return
}
