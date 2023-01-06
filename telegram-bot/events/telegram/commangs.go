package telegram

import (
	"car-park/telegram-bot/lib/er"
	"car-park/telegram-bot/storage"
	"log"
	"regexp"
	"strings"
)

const (
	LoginCmd = "/login"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) processCmd(text string, chatID int, userName string) error {
	text = strings.TrimSpace(text)

	log.Println("got new command '%s' from '%s'", text, userName)

	//TODO:
	// /login command, which save CarPark manager credentials for current user. This also should check if credentials are valid (by CarPark response code)
	// /logout commant, which removes saved credentials file of current user.
	// /help command, which provide information in chat about all available commands and their usage
	// /start command - autocommand when user add bot in telegram. Also must send help message.

	switch text {
	case LoginCmd:
	case HelpCmd:
	case StartCmd:
	default:

	}
}

func (p *Processor) saveCredentials(chatID int, text string, userName string) (err error) {
	r := regexp.MustCompile("[^\\s]+")
	textSlice := r.FindAllString(text, -1)

	// TODO: Add here a http request to check carpark:8888 with provided credentials. GET / with bacic auth should be enough, to get 200 or 401

	creds := &storage.Credentials{
		Username:        userName,
		CarParkLogin:    textSlice[1],
		CarParkPassword: textSlice[2],
	}

	isExist, err := p.storage.IsExistsCredentials(userName)
	if err != nil {
		return er.Wrap("can't check credentials existence while savind credentials", err)
	}

	if isExist {
		return p.tg.SendMessage(chatID, "")
	}
}
