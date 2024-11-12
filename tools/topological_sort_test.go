package tools_test

import (
	"testing"

	"github.com/arknights-w/go-workflow/tools"
)

func TestSuccess(t *testing.T) {
	edges := map[int][]int{
		1: {2, 3},
		2: {4, 5},
		3: {5},
		4: {6},
		5: {6},
		6: {},
	}
	sorted, cycle := tools.TopologicalSort(edges)
	t.Log(sorted)
	t.Log(cycle)
}

func TestCircular(t *testing.T) {
	edges := map[int][]int{
		0: {1},
		1: {2},
		2: {3},
		3: {1},
	}
	sorted, cycle := tools.TopologicalSort(edges)
	t.Log(sorted)
	t.Log(cycle)
}
