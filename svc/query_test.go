/*
Copyright (c) 2019 Ben Morrison (gbmor)

This file is part of Getwtxt.

Getwtxt is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Getwtxt is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Getwtxt.  If not, see <https://www.gnu.org/licenses/>.
*/

package svc // import "git.sr.ht/~gbmor/getwtxt/svc"

import (
	"net"
	"reflect"
	"strings"
	"testing"

	"git.sr.ht/~gbmor/getwtxt/registry"
)

func Test_dedupe(t *testing.T) {
	t.Run("Simple Deduplication Test", func(t *testing.T) {
		start := []string{
			"first",
			"second",
			"third",
			"third",
		}
		finish := dedupe(start)
		if reflect.DeepEqual(start, finish) {
			t.Errorf("Deduplication didn't occur\n")
		}
		if len(finish) != 3 {
			t.Errorf("Ending length not what was expected\n")
		}
	})
}

func Benchmark_dedupe(b *testing.B) {
	start := []string{
		"first",
		"second",
		"third",
		"third",
	}
	for i := 0; i < b.N; i++ {
		dedupe(start)
	}
}

func Test_parseQueryOut(t *testing.T) {
	initTestConf()

	urls := testTwtxtURL
	nick := "getwtxttest"

	out, _, err := registry.GetTwtxt(urls, nil)
	if err != nil {
		t.Errorf("Couldn't set up test: %v\n", err)
	}

	statusmap, err := registry.ParseUserTwtxt(out, nick, urls)
	if err != nil {
		t.Errorf("Couldn't set up test: %v\n", err)
	}

	twtxtCache.AddUser(nick, urls, net.ParseIP("127.0.0.1"), statusmap)

	t.Run("Parsing Status Query", func(t *testing.T) {
		data, err := twtxtCache.QueryAllStatuses()
		if err != nil {
			t.Errorf("%v\n", err)
		}

		out := parseQueryOut(data)

		conv := strings.Split(string(out), "\n")

		if !reflect.DeepEqual(data, conv) {
			t.Errorf("Pre- and Post- parseQueryOut data are inequal:\n%#v\n%#v\n", data, conv)
		}
	})
}

func Benchmark_parseQueryOut(b *testing.B) {
	initTestConf()

	urls := testTwtxtURL
	nick := "getwtxttest"

	out, _, err := registry.GetTwtxt(urls, nil)
	if err != nil {
		b.Errorf("Couldn't set up test: %v\n", err)
	}

	statusmap, err := registry.ParseUserTwtxt(out, nick, urls)
	if err != nil {
		b.Errorf("Couldn't set up test: %v\n", err)
	}

	twtxtCache.AddUser(nick, urls, net.ParseIP("127.0.0.1"), statusmap)

	data, err := twtxtCache.QueryAllStatuses()
	if err != nil {
		b.Errorf("%v\n", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		parseQueryOut(data)
	}
}

func Test_joinQueryOuts(t *testing.T) {
	first := []string{
		"one",
		"two",
		"three",
	}
	second := []string{
		"three",
		"four",
		"five",
		"six",
	}
	t.Run("Joining two string slices", func(t *testing.T) {
		third := joinQueryOuts(first, second)
		if len(third) != (len(first) + len(second) - 1) {
			t.Errorf("Was not combined or deduplicated properly\n")
		}
		fourth := make([]string, 6)
		for i := 0; i < len(first); i++ {
			fourth[i] = first[i]
		}
		for i := 1; i < len(second); i++ {
			fourth[2+i] = second[i]
		}
		if !reflect.DeepEqual(fourth, third) {
			t.Errorf("Output not deeply equal to manual construction\n")
		}
	})
}

func Benchmark_joinQueryOuts(b *testing.B) {
	first := []string{
		"one",
		"two",
		"three",
	}
	second := []string{
		"three",
		"four",
		"five",
		"six",
	}
	for i := 0; i < b.N; i++ {
		joinQueryOuts(first, second)
	}
}

func Test_compositeStatusQuery(t *testing.T) {
	initTestConf()
	mockRegistry()

	t.Run("Composite Query Test", func(t *testing.T) {
		out1, err := twtxtCache.QueryInStatus("sqlite")
		if err != nil {
			t.Errorf("%v\n", err)
		}
		out2, err := twtxtCache.QueryInStatus("Sqlite")
		if err != nil {
			t.Errorf("%v\n", err)
		}
		out3, err := twtxtCache.QueryInStatus("SQLITE")
		if err != nil {
			t.Errorf("%v\n", err)
		}

		outro := make([]string, 0)
		outro = append(outro, out1...)
		outro = append(outro, out2...)
		outro = append(outro, out3...)
		out := dedupe(outro)

		data := compositeStatusQuery("sqlite", nil)

		if !reflect.DeepEqual(out, data) {
			t.Errorf("Returning different data.\nManual: %v\nCompositeQuery: %v\n", out, data)
		}
	})
}

func Benchmark_compositeStatusQuery(b *testing.B) {
	initTestConf()
	statuses, _, _ := registry.GetTwtxt(testTwtxtURL, nil)
	parsed, _ := registry.ParseUserTwtxt(statuses, "getwtxttest", testTwtxtURL)
	_ = twtxtCache.AddUser("getwtxttest", testTwtxtURL, net.ParseIP("127.0.0.1"), parsed)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		compositeStatusQuery("sqlite", nil)
	}

}
