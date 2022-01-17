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
    "reflect"
    "sort"
)

func setToSlice(set map[Id]bool) []Id {
	ret:= make([]Id, 0, len(set))
    for k := range set {
        ret = append(ret, k)
    }
    return ret
}

func arraySortedEqual(a, b []Id) bool {
    if len(a) != len(b) {return false }

    a_copy := make([]Id, len(a))
    b_copy := make([]Id, len(b))

    copy(a_copy, a)
    copy(b_copy, b)

    sort.Slice(a_copy, func(i, j int) bool { return a_copy[i] < a_copy[j] })
    sort.Slice(b_copy, func(i, j int) bool { return b_copy[i] < b_copy[j] })

    return reflect.DeepEqual(a_copy, b_copy)
}

func TestSubscriptions(t *testing.T) {
	var subs []Id
	var expected []Id
    var cases = []struct{
    	req MsgTopic
    	resp []Id
    }{
    	{"1/1/1", []Id{2, 0}},
    	{"1/1/2", []Id{0}},
    	{"1/1/*", []Id{2, 0}},
    	{"1/2/1", []Id{2, 1}},
    	{"1/2/2", []Id{1}},
    	{"1/2/*", []Id{2, 1}},
    }
	router := NewRouterFrq(NO_COLLECT)
	router.Sub(0, "1/1/1")
	router.Sub(0, "1/1/2")
	router.Sub(1, "1/2/1")
	router.Sub(1, "1/2/2")
	router.Sub(2, "1/1/1")
	router.Sub(2, "1/2/1")

	for _, cas := range cases{
		subs = setToSlice(router.Route(cas.req))
		expected = cas.resp
		if !arraySortedEqual(subs, expected) {
			t.Error(
				"\nExpected:", expected,
				"\nReceived:", subs,
				"\nRouter:", router,
			)
		}
	}
}

func TestUnsubscriptions(t *testing.T) {
	var subs []Id
	var expected []Id
    var cases = []struct{
    	req MsgTopic
    	resp []Id
    }{
    	{"1/1/1", []Id{2}},
    	{"1/1/2", []Id{}},
    	{"1/1/*", []Id{2}},
    	{"1/2/1", []Id{2}},
    	{"1/2/2", []Id{}},
    	{"1/2/*", []Id{2}},
    }
	router := NewRouterFrq(NO_COLLECT)
	router.Sub(0, "1/1/1")
	router.Sub(0, "1/1/2")
	router.Sub(1, "1/2/1")
	router.Sub(1, "1/2/2")
	router.Sub(2, "1/1/1")
	router.Sub(2, "1/2/1")
	router.Route("*/*/*")
	router.UnsubAll(0)
	router.UnsubAll(1)
	router.UnsubAll(2)
	router.Collect()
	router.Sub(2, "1/1/1")
	router.Sub(2, "1/2/1")

	for _, cas := range cases{
		subs = setToSlice(router.Route(cas.req))
		expected = cas.resp
		if !arraySortedEqual(subs, expected) {
			t.Error(
				"\nExpected:", expected,
				"\nReceived:", subs,
				"\nRouter:", router,
			)
		}
	}
}

func BenchmarkRouting(b *testing.B) {
	router := NewRouterFrq(NO_COLLECT)
	router.Sub(10, "1/1/1")
	router.Sub(10, "1/1/2")
	router.Sub(20, "1/1/1")
	router.Sub(20, "1/1/2")
	router.Sub(30, "1/1/1")
	router.Sub(30, "1/1/2")
	router.Sub(40, "1/1/1")
	router.Sub(40, "1/1/2")

	router.Sub(11, "1/2/1")
	router.Sub(11, "1/2/2")
	router.Sub(21, "1/2/1")
	router.Sub(21, "1/2/2")
	router.Sub(31, "1/2/1")
	router.Sub(31, "1/2/2")
	router.Sub(41, "1/2/1")
	router.Sub(51, "1/2/2")

	router.Sub(12, "1/1/1")
	router.Sub(12, "1/2/1")
	router.Sub(22, "1/1/1")
	router.Sub(22, "1/2/1")
	router.Sub(32, "1/1/1")
	router.Sub(32, "1/2/1")
	router.Sub(42, "1/1/1")
	router.Sub(42, "1/2/1")

	for i := 0; i < b.N; i++ {
		router.Route("*/*/*")
	}
}
