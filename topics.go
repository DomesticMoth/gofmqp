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

import (
	"strings"
)

const (
	delimiter byte = 47
	uniLevel byte = 42
)

type MsgTopic string

func tirmTopic(t MsgTopic) []string {
	tt := strings.Split(string(t), "/")
	start := 0
	end := len(tt)
	for start < end {
		if tt[start] != "" {
			break
		}
		start += 1
	}
	if end - start < 1 {
		return tt[start:end]
	}
	for start < end {
		if tt[end-1] != "" {
			break
		}
		end -= 1
	}
	return tt[start:end]
}

func ClearTopic(t MsgTopic) MsgTopic{
	trimed := tirmTopic(t)
	out := ""
	for i, level := range trimed {
		out += level
		if i < len(trimed)-1 {
			out += "/"
		}
	}
	return MsgTopic(out)
}

func CompareTopics(a, b MsgTopic) bool {
	aa := tirmTopic(a)
	bb := tirmTopic(b)
	if len(aa) != len(bb) {
		return false
	}
	for i := 0; i < len(aa); i++ {
		if !(aa[i] == bb[i] || aa[i] == "*" || bb[i] == "*") {
			return false
		}
	}
	return true
}

/*
	TODO Write logic for comparing topics in binary format, without conversion to strings
*/
