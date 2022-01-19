/*
    This file is part of gofmqp.

    gofmqp is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    gofmqp is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with gofmqp.  If not, see <https://www.gnu.org/licenses/>.
*/
package gofmqp

type Cache struct{
	messages map[MsgTopic]Message
}

func NewCache() Cache {
	return Cache{ make(map[MsgTopic]Message) }
}

func (cache *Cache) Add(msg Message){
	if len(msg.Body) > 0 {
		cache.messages[*msg.Topic] = msg
	} else {
		delete(cache.messages, *msg.Topic)
	}
}

func (cache *Cache) Get(topic MsgTopic) []Message{	
	ret := []Message{}
	for ptopic, msg := range cache.messages {
		if CompareTopics(topic, ptopic) {
			ret = append(ret, msg)
		}
	}
	return ret
}
