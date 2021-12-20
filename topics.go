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
