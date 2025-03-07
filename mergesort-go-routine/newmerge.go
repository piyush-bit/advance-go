package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const MinSize = 10000 // Threshold for switching to sequential sort

func main2() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // Use all available CPU cores
	
	for size := 100; size <= 10000000; size *= 10 {
		// Generate random array
		arr := make([]int, size)
		for i := 0; i < size; i++ {
			arr[i] = rand.Intn(size)
		}
		
		// Make copies for different sorts
		arr1 := make([]int, len(arr))
		arr2 := make([]int, len(arr))
		copy(arr1, arr)
		copy(arr2, arr)

		fmt.Printf("Size: %d\n", size)
		
		// Sequential merge sort
		start := time.Now()
		mergeSort(arr1)
		fmt.Printf("Sequential: %v\n", time.Since(start))

		// Parallel merge sort
		start = time.Now()
		ParallelMergeSort(arr2)
		fmt.Printf("Parallel: %v\n\n", time.Since(start))
	}
}

// Sequential merge sort
func mergeSort(arr []int) {
	if len(arr) <= 1 {
		return
	}
	
	mid := len(arr) / 2
	mergeSort(arr[:mid])
	mergeSort(arr[mid:])
	merge(arr, mid)
}

// Parallel merge sort
func ParallelMergeSort(arr []int) {
	parallelMergeSortWorker(arr, runtime.NumCPU())
}

func parallelMergeSortWorker(arr []int, maxGoroutines int) {
	// Use sequential sort for small arrays or when we've used up our goroutine budget
	if len(arr) <= MinSize || maxGoroutines <= 1 {
		mergeSort(arr)
		return
	}

	mid := len(arr) / 2
	
	// Create WaitGroup for synchronization
	var wg sync.WaitGroup
	wg.Add(1)
	
	// Sort left half in a new goroutine
	go func() {
		defer wg.Done()
		parallelMergeSortWorker(arr[:mid], maxGoroutines/2)
	}()
	
	// Sort right half in current goroutine
	parallelMergeSortWorker(arr[mid:], maxGoroutines/2)
	
	// Wait for both halves to complete
	wg.Wait()
	
	// Merge the sorted halves
	merge(arr, mid)
}

// In-place merge function
func merge(arr []int, mid int) {
	if len(arr) <= 1 {
		return
	}

	// Create temporary slices
	left := make([]int, mid)
	right := make([]int, len(arr)-mid)
	copy(left, arr[:mid])
	copy(right, arr[mid:])

	// Merge back into original array
	i, j, k := 0, 0, 0
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			arr[k] = left[i]
			i++
		} else {
			arr[k] = right[j]
			j++
		}
		k++
	}

	// Copy remaining elements
	for i < len(left) {
		arr[k] = left[i]
		i++
		k++
	}
	for j < len(right) {
		arr[k] = right[j]
		j++
		k++
	}
}