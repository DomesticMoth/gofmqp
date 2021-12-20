package gofmqp

import (
    "testing"
)

type testpair struct {
	before MsgTopic
	after MsgTopic
}

func TestClearTopic(t *testing.T) {
	cases := []testpair{
		{"a", "a"},
		{"/a", "a"},
		{"a/", "a"},
		{"/a/", "a"},
		{"a/b/c/d", "a/b/c/d"},
		{"/a/b/c/d", "a/b/c/d"},
		{"a/b/c/d/", "a/b/c/d"},
		{"/a/b/c/d/", "a/b/c/d"},
		{"//////a/b/c/d///////", "a/b/c/d"},
		{"/", ""},
		{"///////", ""},
	}
	for _, c := range cases{
		after := ClearTopic(c.before)
		if after != c.after {
			t.Error(
				"Before:", c.before,
				"After:", after,
				"Expected", c.after,
			)
		}
	}
}

func TestCompareTopics(t *testing.T) {
	topics := []MsgTopic {
		"home/living-space/living-room1/temperature",
		"home/living-space/living-room1/*",
		"home/living-space/*/temperature",
		"home/*/living-room1/temperature",
		"*/living-space/living-room1/temperature",
		"*/*/*/*",
	}
	for _, a := range topics {
		for _, b := range topics {
			if !CompareTopics(a, b) {
				t.Error(a, b)
			}
		}
	}
}
