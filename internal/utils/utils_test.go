package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSplitNotFullSuccess(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expected := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10}}

	actual, err := Split(items, 3)

	require.NoError(t, err)

	require.Equal(t, actual, expected)
}

func TestSplitFullSuccess(t *testing.T) {
	items := []int{1, 2, 3, 4}
	expected := [][]int{{1, 2}, {3, 4}}

	actual, err := Split(items, 2)

	require.NoError(t, err)

	require.Equal(t, actual, expected)
}

func TestSplitFailure(t *testing.T) {
	items := []int{1, 2, 3, 4}

	val, err := Split(items, 0)

	require.Nil(t, val)

	require.Error(t, err)
}

func TestSplitEmpty(t *testing.T) {
	items := []int{}

	actual, err := Split(items, 5)

	require.NoError(t, err)

	require.Len(t, actual, 0)
}

func TestFilterSuccess(t *testing.T) {
	items := []int{1, 3, 5, 7, 8}
	expected := []int{1, 5, 8}

	actual := Filter(items, []int{3, 7})

	require.Equal(t, actual, expected)
}

func TestFilterWithoutExcludeSuccess(t *testing.T) {
	items := []int{1, 2, 4}
	expected := []int{1, 2, 4}

	actual := Filter(items, nil)

	require.Equal(t, actual, expected)
}

func TestReverseSuccess(t *testing.T) {
	items := map[string]string{"first": "second", "1": "2", "test-1": "test-3"}
	expected := map[string]string{"second": "first", "2": "1", "test-3": "test-1"}

	actual, err := Reverse(items)
	require.NoError(t, err)
	require.Equal(t, actual, expected)
}

func TestReverseFailure(t *testing.T) {
	items := map[string]string{"first": "second", "third": "second"}

	val, err := Reverse(items)

	require.Nil(t, val)

	require.Error(t, err)
}
