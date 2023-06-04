package scraper

import (
	"fmt"
	"time"

	"github.com/Rosya-edwica/postupi-online/pkg/common/logger"
	"github.com/Rosya-edwica/postupi-online/pkg/common/models"
	"github.com/Rosya-edwica/postupi-online/pkg/scraper/html"
	"github.com/gocolly/colly"
)

func (scr *Scraper) ScrapeContacts(url string) (err error) {
	fmt.Println(url)
	contact := models.Contacts{}
	c := colly.NewCollector()
	badGateway := html.CheckBadGateway(c)
	if badGateway {
		return
	}
	c.SetRequestTimeout(30 * time.Second)

	c.OnHTML("span.contact-icon.contact-icon_sm.mail", func(h *colly.HTMLElement) {
		contact.Email = h.Text
	})
	c.OnHTML("span.contact-icon.contact-icon_sm.phone", func(h *colly.HTMLElement) {
		contact.Phone = h.Text
	})
	c.OnHTML("span.contact-icon.contact-icon_sm.address", func(h *colly.HTMLElement) {
		contact.Address = h.Text
	})
	c.OnHTML("span.contact-icon.contact-icon_sm.site", func(h *colly.HTMLElement) {
		contact.WebSite = h.Text
	})
	contact.VuzId = getVuzId(url)
	err = c.Post(url+"contacts/", html.Headers)
	if err != nil {
		fmt.Println("Catched the error. Program stopped to sleep of 10 seconds.")
		time.Sleep(10 * time.Second)
		err = scr.ScrapeContacts(url)
		checkErr(err)
	}
	scr.Db.SaveContacts(contact)
	checkErr(err)
	logger.Log.Printf("Contact:%s\n", contact.VuzId)
	return
}
