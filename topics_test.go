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
