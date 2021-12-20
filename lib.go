package gofmqp

import (
	"io"
	//"encoding/binary"
)

type PackageType bool

const (
	PackageTypePub bool = false
	PackageTypeSub = true
)

type SubscribeType bool

const (
	SubscribeTypeSub bool = false
	SubscribeTypeUnsub = true
)

type TopicSystem bool

const (
	TopicSystemRegular bool = false
	TopicSystemFeedback = true
)

type RawMessage struct{
	Type PackageType
	SubType SubscribeType
	TopicType TopicSystem
	Debug bool
	LastWill bool
	Cache bool
	Topic []byte
	Body []byte
}

type Message struct{
	Type PackageType
	SubType SubscribeType
	TopicType TopicSystem
	Debug bool
	LastWill bool
	Cache bool
	Topic *MsgTopic
	Body []byte
}


type MsgReader struct{
	input io.Reader
	buf []byte
}

func NewMsgReader(input io.Reader, buf int) MsgReader{
	if buf < 1024 {
		buf = 1024
	}
	b := make([]byte, buf)
	return MsgReader{input, b}
}


/*func (r *MsgReader) NextRaw() (msg RawMessage, err error) {
	_, err = io.ReadAtLeast(r.input, r.buf, 6)
	if err != nil { return }
	msg.Type = PackageType((r.buf[0] & 128) > 0)
	msg.SubType = SubscribeType((r.buf[0] & 64) > 0)
	msg.TopicType = TopicSystem((r.buf[0] & 32) > 0)
	msg.Debug = (r.buf[0] & 16) > 0
	msg.LastWill = (r.buf[0] & 8) > 0
	msg.Cache = (r.buf[0] & 4) > 0
	topicLen := uint8(r.buf[1])
	bodyLen := binary.BigEndian.Uint32(r.buf[2:6])
	data := make([]byte, int(topicLen)+int(bodyLen))
	topic := data[:topicLen]
	body := data[topicLen+1:]
	copy(topic, r.buf[7:7+topicLen])
	//
	msg.Topic = topic
	msg.Body = body
	return
}*/

/*
	TODO Add message body reading function
*/
