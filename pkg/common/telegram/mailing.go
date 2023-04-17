package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tidwall/gjson"
)

const appName = "Сообщение из программы: Postupi-Online:\n\n"

func GetUpdates() (countMessages int) {
	json, _ := getJson(fmt.Sprintf("%s/getUpdates", getUrl()))
	return len(gjson.Get(json, "result").Array())
}

func Mailing(text string) {
	for _, chat := range chats {
		SendMessage(text, chat)
	}
}

func SendMessage(text string, chatId string) (bool, error) {
	url := fmt.Sprintf("%s/sendMessage", getUrl())
	body, _ := json.Marshal(map[string]string{
		"chat_id": chatId,
		"text":    appName + text,
	})
	response, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()
	return true, nil
}
