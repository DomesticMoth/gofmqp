package gofmqp

import (
    "testing"
    "bytes"
)

func compare(a,b Message) bool {
	return a.Type == b.Type && 
			a.SubType == b.SubType &&
			a.TopicType == b.TopicType &&
			a.Debug == b.Debug &&
			a.LastWill == b.LastWill &&
			a.Cache == b.Cache &&
			*a.Topic == *b.Topic &&
			bytes.Equal(a.Body, b.Body)
}

func TestRawIO(t *testing.T) {
	topic := MsgTopic("12345/6789")
	msg := Message{
		PackageTypePub,
		SubscribeTypeUnsub,
		TopicSystemRegular,
		true,
		false,
		true,
		&topic,
		[]byte{0,1,2,3,4,5,6,7,8},
	}
	buf := new(bytes.Buffer)
	writer := NewMsgWriter(buf)
	reader := NewMsgReader(buf)
	err := writer.Send(&msg)
	if err != nil { t.Error(err) }
	m, err := reader.NextUnchecked()
	if err != nil { t.Error(err) }
	if !compare(m, msg) {
		t.Error("A:", msg, *msg.Topic, "\nB:", m, *m.Topic)
	}
}
