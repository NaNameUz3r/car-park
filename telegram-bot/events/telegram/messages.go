package telegram

const msgHelp = `This is CarParkDev helper bot. In order to make requests to CarParkDev service you should login with
manager credentials by sending me command "/login <USERNAME> <PASSWORD>".

//TODO: Add help about reports fetch command when will be implemented.

Send command /logout if you need to login with another credentials pair.

Send /help to see this message again.

`

const msgHello = "Hello Friend. \n\n" + msgHelp

const (
	msgUnknownCommand = "I do not know this command. See /help"
	msgNotLoggedIn    = `You have not logged in right now. Use "/login <USERNAME> <PASSWORD>"`
)
