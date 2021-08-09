package utils

import (
	"reflect"
	"testing"
)

func TestSplitNotFullSuccess(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expected := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10}}

	actual, err := Split(items, 3)

	if err != nil {
		t.Fatalf("Error %s not expected", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nActual:%v\nExpect:%v", actual, expected)
	}
}

func TestSplitFullSuccess(t *testing.T) {
	items := []int{1, 2, 3, 4}
	expected := [][]int{{1, 2}, {3, 4}}

	actual, err := Split(items, 2)

	if err != nil {
		t.Fatalf("Error %s not expected", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nActual:%v\nExpect:%v", actual, expected)
	}
}

func TestSplitFailure(t *testing.T) {
	items := []int{1, 2, 3, 4}

	_, err := Split(items, 0)

	if err == nil {
		t.Fatal("Error expected")
	}
}

func TestSplitEmpty(t *testing.T) {
	items := []int{}

	actual, err := Split(items, 5)

	if err != nil {
		t.Fatalf("Error %s not expected", err)
	}

	if len(actual) != 0 {
		t.Fatal("Expected empty slice")
	}
}

func TestFilterSuccess(t *testing.T) {
	items := []int{1, 3, 5, 7, 8}
	expected := []int{1, 5, 8}

	actual := Filter(items, []int{3, 7})

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nActual:%v\nExpect:%v", actual, expected)
	}
}

func TestFilterWithoutExcludeSuccess(t *testing.T) {
	items := []int{1, 2, 4}
	expected := []int{1, 2, 4}

	actual := Filter(items, nil)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nActual:%v\nExpect:%v", actual, expected)
	}
}

func TestReverseSuccess(t *testing.T) {
	items := map[string]string{"first": "second", "1": "2", "test-1": "test-3"}
	expected := map[string]string{"second": "first", "2": "1", "test-3": "test-1"}

	actual, err := Reverse(items)

	if err != nil {
		t.Fatalf("Error %s not expected", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nActual:%v\nExpect:%v", actual, expected)
	}
}

func TestReverseFailure(t *testing.T) {
	items := map[string]string{"first": "second", "third": "second"}

	_, err := Reverse(items)

	if err == nil {
		t.Fatal("Error expected")
	}
}
