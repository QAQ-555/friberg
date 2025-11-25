package isocket

type Message struct {
	MsgType int         `json:"msgType"`
	Data    interface{} `json:"data"`
}
