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
	Long:  `Start the interactive process of exploring sorting algorithms`,
	Run: func(cmd *cobra.Command, args []string) {
		startIllustration()
	},
}

type promptContent struct {
	errorMsg string
	label    string
}

type algoProcess struct {
	original []int
	current  []int
	guesses  int
	done     bool
	stepSort func(arr []int) (bool, bool)
}

func (ap *algoProcess) next(guess []int) bool {
	changed := false

	// Sort step by step until a change's happended in the array
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

func startIllustration() {
	// Prompt: sorting algorithm
	arrLen := 5
	items := []string{
		"Bubblesort",
		"Heapsort",
		"Mergesort",
		"Quicksort",
	}

	algorithms := map[string]func() func([]int) (bool, bool){
		"Bubblesort": bubbleSort,
		"Heapsort":   heapSort,
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
		stepSort: algorithms[chosenItem](),
	}

	for !ap.done {
		validationMsg := `Please enter how would the array look at the next step: 
						In the format of the JSON array: [7, 2, 5, 4, ...]`
		guessPC := promptContent{
			validationMsg,
			fmt.Sprintf("How would the array %v look at the next step? ", ap.current),
		}
		guessStr := promptGetInput(guessPC)

		guess, err := atoiSlice(guessStr)
		if err != nil {
			fmt.Println(validationMsg)
			continue
		}

		ap.next(guess)
	}
}

func generateRandomArray(len int) []int {
	// Provide seed
	rand.Seed(time.Now().Unix())

	return rand.Perm(len)
}

func atoiSlice(input string) ([]int, error) {
	var guess []int
	err := json.Unmarshal([]byte(input), &guess)
	if err != nil {
		return nil, errors.New("validation error")
	}

	return guess, nil
}

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

func bubbleSort() func(arr []int) (bool, bool) {
	i, j := 0, 1

	return func(arr []int) (finished bool, changed bool) {
		changed = false

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

func heapSort() func(arr []int) (bool, bool) {
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
