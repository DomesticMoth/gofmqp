package gofmqp

import (
	"io"
	"encoding/binary"
	"errors"
)

type PackageType bool

const (
	PackageTypePub PackageType = false
	PackageTypeSub = true
)

type SubscribeType bool

const (
	SubscribeTypeSub SubscribeType = false
	SubscribeTypeUnsub = true
)

type TopicSystem bool

const (
	TopicSystemRegular TopicSystem = false
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
	hbuf []byte
}

func NewMsgReader(input io.Reader) MsgReader{
	hbuf := make([]byte, 6)
	return MsgReader{input, hbuf}
}


func (r *MsgReader) NextRaw() (msg RawMessage, err error) {
	_, err = io.ReadFull(r.input, r.hbuf)
	if err != nil { return }
	msg.Type      = PackageType((r.hbuf[0]   & 128) > 0)
	msg.SubType   = SubscribeType((r.hbuf[0] & 64)  > 0)
	msg.TopicType = TopicSystem((r.hbuf[0]   & 32)  > 0)
	msg.Debug     = (r.hbuf[0]               & 16)  > 0
	msg.LastWill  = (r.hbuf[0]               & 8)   > 0
	msg.Cache     = (r.hbuf[0]               & 4)   > 0
	topicLen := uint8(r.hbuf[1])
	bodyLen := binary.BigEndian.Uint32(r.hbuf[2:])
	data := make([]byte, int(topicLen)+int(bodyLen))
	_, err = io.ReadFull(r.input, data)
	if err != nil { return }
	topic := data[:topicLen]
	body := data[topicLen:]
	msg.Topic = topic
	msg.Body = body
	return
}

func (r *MsgReader) NextUnchecked() (msg Message, err error) {
	raw, err := r.NextRaw()
	if err != nil { return }
	topic := MsgTopic(string(raw.Topic))
	msg = Message{
		raw.Type,
		raw.SubType,
		raw.TopicType,
		raw.Debug,
		raw.LastWill,
		raw.Cache,
		&topic,
		raw.Body,
	}
	return
}

func (r *MsgReader) Next() (msg Message, err error) {
	msg, err = r.NextUnchecked()
	if err != nil { return }
	topic := ClearTopic(*msg.Topic)
	msg.Topic = &topic
	return
}


type MsgWriter struct{
	out io.Writer
	hbuf []byte
}

func NewMsgWriter(out io.Writer) MsgWriter{
	hbuf := make([]byte, 6)
	return MsgWriter{out, hbuf}
}

func booltobyte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func (w *MsgWriter) SendRaw(msg *RawMessage) (err error) {
	tlen := len(msg.Topic)
	if tlen > 255 {
		err = errors.New("Too long topic")
	}
	w.hbuf[0] = 0
	w.hbuf[0] = w.hbuf[0] | (booltobyte(bool(msg.Type))      << 7)
	w.hbuf[0] = w.hbuf[0] | (booltobyte(bool(msg.SubType))   << 6)
	w.hbuf[0] = w.hbuf[0] | (booltobyte(bool(msg.TopicType)) << 5)
	w.hbuf[0] = w.hbuf[0] | (booltobyte(bool(msg.Debug))     << 4)
	w.hbuf[0] = w.hbuf[0] | (booltobyte(bool(msg.LastWill))  << 3)
	w.hbuf[0] = w.hbuf[0] | (booltobyte(bool(msg.Cache))     << 2)
	w.hbuf[1] = uint8(tlen)
	ln := uint32(len(msg.Body))
	binary.BigEndian.PutUint32(w.hbuf[2:], ln)
	_, err = w.out.Write(w.hbuf)
	if err != nil { return }
	_, err = w.out.Write(msg.Topic)
	if err != nil { return }
	_, err = w.out.Write(msg.Body)
	return
}

func (w *MsgWriter) Send(msg *Message) (err error) {
	topic := msg.Topic
	raw := RawMessage{
		msg.Type,
		msg.SubType,
		msg.TopicType,
		msg.Debug,
		msg.LastWill,
		msg.Cache,
		[]byte(string(*topic)),
		msg.Body,
	}
	return w.SendRaw(&raw)
}


/*

TODO Add unix socket, tcp and tls realisations

*/

