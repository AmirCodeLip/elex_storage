package message_broker

type BaseEvent struct {
	Topic string
	Data  []byte
}
