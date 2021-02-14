package kudos

type Message struct {
	Id      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type responseRtmStart struct {
	Ok    bool         `json:"ok"`
	Error string       `json:"error"`
	Url   string       `json:"url"`
	Self  responseSelf `json:"self"`
}

type responseSelf struct {
	Id string `json:"id"`
}
