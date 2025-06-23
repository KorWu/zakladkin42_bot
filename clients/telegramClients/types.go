package telegramClients

type Update struct {
	UpdateID int              `json:"update_id"`
	Message  *IncomingMessage `json:"message"`
}

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From User   `json:"from"`
	Chat Chat   `json:"chat"`
}

type Chat struct {
	ChatID int `json:"id"`
}

type User struct {
	Username string `json:"username"`
}
