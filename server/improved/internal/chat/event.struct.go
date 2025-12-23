package chat

/*
Server → Client (new contract)
*/
type ChatEvent struct {
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`

	// message
	Sender  string `json:"sender,omitempty"`
	Content string `json:"content,omitempty"`

	// system
	Event string `json:"event,omitempty"`
	User  string `json:"user,omitempty"`
}
