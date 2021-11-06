package cmd

import (
	"testing"
)

func TestBubbleSort(t *testing.T) {
	t.Run("simple case", func(t *testing.T) {
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
			{1, 2, 2, 4, 5, 6, 9, 10}, // i = 6, j = 1, end
		}

		current := make([]int, len(snapshots[0]))
		copy(current, snapshots[0])

		ap := algoProcess{
			current:  current,
			stepSort: bubbleSort(snapshots[0]),
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

func TestHeapSort(t *testing.T) {
	t.Run("simple case", func(t *testing.T) {
		snapshots := [][]int{
			{3, 5, 1, 4, 7, 6, 2, 8}, // initial array
			{3, 5, 1, 8, 7, 6, 2, 4}, // i = 3, heapify start
			{3, 5, 6, 8, 7, 1, 2, 4}, // i = 2
			{3, 8, 6, 5, 7, 1, 2, 4}, // i = 1
			{8, 7, 6, 5, 3, 1, 2, 4}, // i = 0
			{7, 5, 6, 4, 3, 1, 2, 8}, // i = 6, sorting start
			{6, 5, 2, 4, 3, 1, 7, 8}, // i = 5
			{5, 4, 2, 1, 3, 6, 7, 8}, // i = 4
			{4, 3, 2, 1, 5, 6, 7, 8}, // i = 3
			{3, 1, 2, 4, 5, 6, 7, 8}, // i = 2
			{2, 1, 3, 4, 5, 6, 7, 8}, // i = 1
			{1, 2, 3, 4, 5, 6, 7, 8}, // i = 0, end
		}

		current := make([]int, len(snapshots[0]))
		copy(current, snapshots[0])

		ap := algoProcess{
			current:  current,
			stepSort: heapSort(snapshots[0]),
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

func TestMergeSort(t *testing.T) {
	t.Run("simple case", func(t *testing.T) {
		snapshots := [][]int{
			{5, 3, 4, 1, 7, 6, 8, 2}, // initial array
			{3, 5, 4, 1, 7, 6, 8, 2}, // 3, 5
			{3, 5, 1, 4, 7, 6, 8, 2}, // 1, 4
			{1, 3, 4, 5, 7, 6, 8, 2}, // 1, 3, 4, 5
			{1, 3, 4, 5, 6, 7, 8, 2}, // 6, 7
			{1, 3, 4, 5, 6, 7, 2, 8}, // 2, 8
			{1, 3, 4, 5, 2, 6, 7, 8}, // 2, 6, 7, 8
			{1, 2, 3, 4, 5, 6, 7, 8}, // end
		}

		current := make([]int, len(snapshots[0]))
		copy(current, snapshots[0])

		ap := algoProcess{
			current:  current,
			stepSort: mergeSort(snapshots[0]),
		}
		i := 1

		for !ap.done {
			if guessed := ap.next(snapshots[i]); !guessed {
				t.Fatalf("failed at step %d, the snapshot was %v", i, snapshots[i])
			}
			i++
		}
	})

	t.Run("complex case", func(t *testing.T) {
		snapshots := [][]int{
			{5, 3, 4, 1, 7, 6, 8, 2, 14, 12, 13, 10, 11, 9}, // initial array
			{3, 4, 5, 1, 7, 6, 8, 2, 14, 12, 13, 10, 11, 9}, // 3, 4, 5
			{3, 4, 5, 1, 6, 7, 8, 2, 14, 12, 13, 10, 11, 9}, // 1, 6, 7, 8
			{1, 3, 4, 5, 6, 7, 8, 2, 14, 12, 13, 10, 11, 9}, // 1, 3, 4, 5, 6, 7, 8
			{1, 3, 4, 5, 6, 7, 8, 2, 12, 14, 13, 10, 11, 9}, // 12, 14
			{1, 3, 4, 5, 6, 7, 8, 2, 12, 14, 10, 13, 11, 9}, // 10, 13
			{1, 3, 4, 5, 6, 7, 8, 2, 12, 14, 10, 13, 9, 11}, // 9, 11
			{1, 3, 4, 5, 6, 7, 8, 2, 12, 14, 9, 10, 11, 13}, // 9, 10, 11, 13
			{1, 3, 4, 5, 6, 7, 8, 2, 9, 10, 11, 12, 13, 14}, // 2, 9, 10, 11, 12, 13, 14
			{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}, // end
		}

		current := make([]int, len(snapshots[0]))
		copy(current, snapshots[0])

		ap := algoProcess{
			current:  current,
			stepSort: mergeSort(snapshots[0]),
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

func TestQuickSort(t *testing.T) {
	t.Run("simple case", func(t *testing.T) {
		snapshots := [][]int{
			{5, 3, 2, 1, 7, 10, 6}, // initial array
			{5, 3, 2, 1, 6, 10, 7},
			{1, 3, 2, 5, 6, 10, 7},
			{1, 2, 3, 5, 6, 10, 7},
			{1, 2, 3, 5, 6, 7, 10}, // already sorted at this step
			{1, 2, 3, 5, 6, 7, 10}, // second time is required since the algorithm needs to end
		}

		current := make([]int, len(snapshots[0]))
		copy(current, snapshots[0])

		ap := algoProcess{
			current:  current,
			stepSort: quickSort(snapshots[0]),
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
