package cmd

import (
	"testing"
)

func TestBubbleSort(t *testing.T) {
	t.Run("test bubble sort", func(t *testing.T) {
		snapshots := [][]int{
			{5, 1, 4, 2, 6, 9, 10, 2}, // initial array
			{1, 5, 4, 2, 6, 9, 10, 2}, // i = 0, j = 1, switch
			{1, 4, 5, 2, 6, 9, 10, 2}, // i = 0, j = 2, switch
			{1, 4, 2, 5, 6, 9, 10, 2}, // i = 0, j = 3, switch
			{1, 4, 2, 5, 6, 9, 2, 10}, // i = 0, j = 7, switch
			{1, 2, 4, 5, 6, 9, 2, 10}, // i = 1, j = 2, switch
			{1, 2, 4, 5, 6, 2, 9, 10}, // i = 1, j = 6, switch
			{1, 2, 4, 5, 2, 6, 9, 10}, // i = 2, j = 5, switch
			{1, 2, 4, 2, 5, 6, 9, 10}, // i = 3, j = 4, switch
			{1, 2, 2, 4, 5, 6, 9, 10}, // i = 4, j = 3, switch
			{1, 2, 2, 4, 5, 6, 9, 10}, // i = 4, j = 4, iterate i
		}

		current := make([]int, len(snapshots[0]))
		copy(current, snapshots[0])

		ap := algoProcess{
			current:  current,
			stepSort: bubbleSort(),
		}
		i := 1

		for !ap.done {
			if guessed := ap.next(snapshots[i]); !guessed {
				t.Fatalf("failed at step %d, the snapshot was %v", i, snapshots[i])
			}
			i++
		}
	})

}
