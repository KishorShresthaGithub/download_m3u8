package src

import (
	"fmt"
	"testing"
)

type CompareFunction[S any, E comparable] func(data S) (E, E)

func Runner[S any, T comparable](t *testing.T, dataGroup []S, f CompareFunction[S, T]) {

	for i := range dataGroup {
		t.Run(fmt.Sprintf("Test [%v]", i), func(t *testing.T) {
			actual, expected := f(dataGroup[i])

			if actual != expected {
				t.Logf("expected %v, actual %v", expected, actual)
				t.Fail()
			}

		})

	}

}

func TestCheckLink(t *testing.T) {

	type DataGroup struct {
		in       string
		expected bool
	}

	data := []DataGroup{{"https://asd.com/asd", true}, {"ttps://asd.com/asd", false}, {"aasd.comhttp:.//asd", false}, {"http://asd/asd", true}}

	Runner(t, data, func(data DataGroup) (bool, bool) {
		actual := CheckIfLineIsLink(data.in)
		expected := data.expected

		return actual, expected
	})

}

func TestBasename(t *testing.T) {
	type DataGroup struct {
		in       string
		expected string
	}

	data := []DataGroup{{"https://asd.com/asd", "asd"}, {"ttps://asd.com/asd1", "asd1"}, {"aasd.comhttp:.//path", "path"}, {"http://asd/123", "123"}}

	Runner(t, data, func(data DataGroup) (string, string) {
		return GetBasename(data.in), data.expected
	})

}

func TestHash(t *testing.T) {
	type DataGroup struct {
		in       string
		expected bool
	}

	data := []DataGroup{{"#asd", true}, {"#123", true}, {"#@#$", true}, {"asd#123", false}, {"123", false}, {"http://asd/123", false}}

	Runner(t, data, func(data DataGroup) (bool, bool) {
		return CheckIfLineStartsWithHash(data.in), data.expected
	})

}

func TestBase(t *testing.T) {
	type DataGroup struct {
		in       string
		expected string
	}

	data := []DataGroup{{"https://asd.com/1", "https://asd.com/"}, {"https://asd.com/1?q=2", "https://asd.com/"}, {"https://asd.com/1/2/3", "https://asd.com/1/2"}}

	Runner(t, data, func(data DataGroup) (string, string) {
		return GetBaseUrl(data.in), data.expected
	})

}
