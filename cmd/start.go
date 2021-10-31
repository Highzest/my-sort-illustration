/*
Copyright © 2021 Yermek Menzhessarov <epmek96@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the interactive process of exploring sorting algorithms",
	Long: `Start the interactive process of exploring sorting algorithms. 
			Choose a sorting algorithm and then guess the next iteration of it.`,
	Run: func(cmd *cobra.Command, args []string) {
		startIllustration()
	},
}

type promptContent struct {
	errorMsg string
	label    string
}

// algoProcess is the main struct that oversees the whole process of sorting step by step.
type algoProcess struct {
	original []int
	current  []int
	guesses  int
	done     bool
	stepSort func(arr []int) (bool, bool)
}

// next sorts until there is a change in the slice,
// then checks and returns whether the provided guess is equal to the current iteration of the slice.
func (ap *algoProcess) next(guess []int) bool {
	changed := false

	// Sort step by step until a change's happened in the array
	for !ap.done && !changed {
		ap.done, changed = ap.stepSort(ap.current)
	}

	correct := isEqual(ap.current, guess)

	if correct {
		ap.guesses++
	} else {
		fmt.Printf("Correct one: %v\n", ap.current)
	}

	return correct
}

func init() {
	rootCmd.AddCommand(startCmd)
}

// promptGetInput provides a user with a prompt and returns the typed answer of the user.
func promptGetInput(pc promptContent) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.errorMsg)
		} else if _, err := atoiSlice(input); err != nil {
			return errors.New(pc.errorMsg)
		}

		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }}",
		Valid:   "{{ . | green }}",
		Invalid: "{{ . | red }}",
		Success: "{{ . | bold }}",
	}

	prompt := promptui.Prompt{
		Label:     pc.label,
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return result
}

// promptGetSelect provides a user with choices and returns the chosen item.
func promptGetSelect(items []string, pc promptContent) string {
	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.Select{
			Label: pc.label,
			Items: items,
		}

		index, result, err = prompt.Run()
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return result
}

// startIllustration is the main function that the user interacts with.
func startIllustration() {
	// Prompt: sorting algorithm
	arrLen := 5
	items := []string{
		"Bubblesort",
		"Heapsort",
		"Mergesort",
		//"Quicksort",
	}

	algorithms := map[string]func(original []int) func([]int) (bool, bool){
		"Bubblesort": bubbleSort,
		"Heapsort":   heapSort,
		"Mergesort":  mergeSort,
	}

	algPC := promptContent{
		"Please choose the sorting algorithm. ",
		"Which algorithm would you choose? ",
	}
	chosenItem := promptGetSelect(items, algPC)

	// Generate array of ints
	arr := generateRandomArray(arrLen)

	original := make([]int, arrLen)
	copy(original, arr)

	// Start the main loop for illustrating the sorting algorithm
	ap := algoProcess{
		original: original,
		current:  arr,
		stepSort: algorithms[chosenItem](original),
	}

	for !ap.done {
		validationMsg := `Please enter how would the array look at the next step when it changes: 
						In the format of the JSON array: [7, 2, 5, 4, ...]`
		promptMsg, err := json.Marshal(ap.current)
		if err != nil {
			fmt.Println(validationMsg)
			continue
		}

		guessPC := promptContent{
			validationMsg,
			fmt.Sprintf("How would the array %s change at the next step? ", string(promptMsg)),
		}
		guessStr := promptGetInput(guessPC)

		guess, err := atoiSlice(guessStr)
		if err != nil {
			fmt.Println(validationMsg)
			continue
		}

		ap.next(guess)
	}

	fmt.Println("Congrats! You won!")
}

// generateRandomArray generates random slice of integers with len as the parameter for the slice's length.
func generateRandomArray(len int) []int {
	// Provide seed
	rand.Seed(time.Now().Unix())

	return rand.Perm(len)
}

// atoiSlice parses a JSON array of integers into a slice of integers.
func atoiSlice(input string) ([]int, error) {
	var guess []int
	err := json.Unmarshal([]byte(input), &guess)
	if err != nil {
		return nil, errors.New("validation error")
	}

	return guess, nil
}

// isEqual checks whether the two integer slices are equal.
func isEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// bubbleSort is an "iterative" implementation of the Bubble Sort sorting algorithm.
func bubbleSort(original []int) func(arr []int) (bool, bool) {
	i, j := 0, 1

	return func(arr []int) (finished bool, changed bool) {
		if arr[j] < arr[j-1] {
			arr[j], arr[j-1] = arr[j-1], arr[j]
			changed = true
		}

		j++

		if !(j < len(arr)-i) {
			i++

			if !(i < len(arr)) {
				return true, changed
			}

			j = 1
		}

		return false, changed
	}
}

// heapSort is an "iterative" implementation of the Heap Sort sorting algorithm.
func heapSort(original []int) func(arr []int) (bool, bool) {
	var i int
	var step int
	return func(arr []int) (finished bool, changed bool) {
		// Step 0: initialization
		if step == 0 {
			// Parent of the last element in the array (which is a leaf)
			i = ((len(arr) - 1) - 1) / 2
			step = 1
		}

		// Step 1: heapify
		if step == 1 {
			if !(i >= 0) {
				// reset values for step 2
				step = 2
				i = len(arr) - 1
				return false, false
			}

			siftDown(arr, i, len(arr)-1)

			i--
			return false, true
		}

		// Step 2: sorting process
		if step == 2 {
			arr[0], arr[i] = arr[i], arr[0]
			i--
			siftDown(arr, 0, i)

			if !(i > 0) {
				return true, true
			}

			return false, true
		}

		return false, false
	}
}

// siftDown pushes the root of the heap down until its position is found.
func siftDown(heap []int, lo, hi int) {
	root := lo
	for {
		child := root*2 + 1
		if child > hi {
			break
		}
		if child+1 < hi && heap[child] < heap[child+1] {
			child++
		}
		if heap[root] < heap[child] {
			heap[root], heap[child] = heap[child], heap[root]
			root = child
		} else {
			break
		}

	}
}

type mergeSortStep struct {
	leftBeg, leftEnd, rightBeg, rightEnd int
}

type mergeSortPath struct {
	path     []mergeSortStep
	currStep int
}

// calculateMergeSortPath traverses indices just like the Merge Sort algorithm would
// in order to populate every mergeSortPath's mergeSortStep with all the merge function's argument values.
func (m *mergeSortPath) calculateMergeSortPath(indices []int) {
	if len(indices) <= 1 {
		return
	}
	middle := len(indices) / 2

	m.calculateMergeSortPath(indices[:middle])
	m.calculateMergeSortPath(indices[middle:])

	m.path = append(m.path, mergeSortStep{
		indices[0],
		indices[middle-1] + 1,
		indices[middle],
		indices[len(indices)-1] + 1,
	})
}

func (m *mergeSortPath) iterateStep() {
	m.currStep++
}

func (m *mergeSortPath) isEnd() bool {
	return m.currStep == len(m.path)
}

func (m *mergeSortPath) getLeftBeg() int {
	return m.path[m.currStep].leftBeg
}

func (m *mergeSortPath) getLeftEnd() int {
	return m.path[m.currStep].leftEnd
}

func (m *mergeSortPath) getRightBeg() int {
	return m.path[m.currStep].rightBeg
}

func (m *mergeSortPath) getRightEnd() int {
	return m.path[m.currStep].rightEnd
}

// mergeSort is an "iterative" implementation of the Merge Sort sorting algorithm.
func mergeSort(original []int) func(arr []int) (bool, bool) {
	var seq []int
	m := mergeSortPath{[]mergeSortStep{}, 0}

	// generate sequence 0 to length-1
	for i := 0; i < len(original); i++ {
		seq = append(seq, i)
	}

	m.calculateMergeSortPath(seq)

	return func(arr []int) (finished bool, changed bool) {
		sorted := merge(
			arr[m.getLeftBeg():m.getLeftEnd()],
			arr[m.getRightBeg():m.getRightEnd()],
		)

		for i, j := m.getLeftBeg(), 0; i < m.getRightEnd(); i, j = i+1, j+1 {
			if arr[i] != sorted[j] {
				changed = true
			}
			arr[i] = sorted[j]
		}
		m.iterateStep()

		return m.isEnd(), changed
	}
}

// merge combines two slices in a way that at each iteration
// the smallest of the two first elements from each slice is inserted first to the new slice.
func merge(left, right []int) []int {
	result := make([]int, len(left)+len(right))

	for i := 0; len(left) > 0 || len(right) > 0; i++ {
		if len(left) > 0 && len(right) > 0 {
			if left[0] < right[0] {
				result[i] = left[0]
				left = left[1:]
			} else {
				result[i] = right[0]
				right = right[1:]
			}
		} else if len(left) > 0 {
			result[i] = left[0]
			left = left[1:]
		} else if len(right) > 0 {
			result[i] = right[0]
			right = right[1:]
		}
	}

	return result
}
