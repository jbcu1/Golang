package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"telegram/config"
	"telegram/models"
)

func main() {
	apiURL := config.BotAPI + config.BotToken
	offset := 0
	for {
		updates, err := getUpdates(apiURL, offset)
		if err != nil {
			log.Panicln("Something went wrong", err.Error())
		}
		fmt.Println(updates)
		for _, update := range updates {
			err = respons(apiURL, update)
			offset = update.UpdateId + 1
		}

	}
}

func getUpdates(apiURL string, offset int) ([]models.Update, error) {
	resp, err := http.Get(apiURL + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse models.RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}

	return restResponse.Result, nil

}

func respons(apiURL string, update models.Update) error {
	var botMessage models.BotMessage
	botMessage.ChatID = update.Message.Chat.ChatID
	botMessage.Text = update.Message.Text
	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(apiURL+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
}
