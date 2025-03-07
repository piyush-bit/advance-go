package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main3() {
	// compare time for single threaded and multithreaded merge sort have array 100 , 1000 , 10000 , 100000 and print the time taken by both single threaded and multithreaded merge sort
	for size := 10000; size <= 1000000000; size = size * 10 {
		arr := make([]int, size)
		for i := 0; i < size; i++ {
			arr[i] = rand.Intn(size)
		}
		fmt.Printf("%d,",size)

		start := time.Now()
		mergesort(arr)
		end := time.Now()
		fmt.Printf("%s," ,end.Sub(start))

		start = time.Now()
		sorted := make(chan []int)
		go MergesortMultiLimited(arr,sorted,10)
		<-sorted
		end = time.Now()
		fmt.Printf("%s\n" ,end.Sub(start))
	}

}

func Merge(arr1 []int, arr2 []int) []int {
	arr3 := make([]int, len(arr1)+len(arr2))
	i, j, k := 0, 0, 0
	for i < len(arr1) && j < len(arr2) {
		if arr1[i] < arr2[j] {
			arr3[k] = arr1[i]
			i++
			k++
		} else {
			arr3[k] = arr2[j]
			j++
			k++
		}
	}
	for i < len(arr1) {
		arr3[k] = arr1[i]
		i++
		k++
	}
	for j < len(arr2) {
		arr3[k] = arr2[j]
		j++
		k++
	}

	return arr3
}

func mergesort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	mid := len(arr) / 2
	left := mergesort(arr[:mid])
	right := mergesort(arr[mid:])
	return Merge(left, right)
}


func MergesortMultiLimited(arr []int, mergeE chan []int, depth int) {
	if len(arr) <= 1 {
		mergeE <- arr
		return
	}

	if depth <= 0 {
		// Use sequential sort when depth limit is reached
		mergeE <- mergesort(arr) // Replace with any efficient sequential sort
		return
	}

	mid := len(arr) / 2
	tempChan := make(chan []int, 2)

	go func() {
		leftChan := make(chan []int)
		go MergesortMultiLimited(arr[:mid], leftChan, depth-1)
		tempChan <- <-leftChan
		close(leftChan)
	}()

	go func() {
		rightChan := make(chan []int)
		go MergesortMultiLimited(arr[mid:], rightChan, depth-1)
		tempChan <- <-rightChan
		close(rightChan)
	}()

	left := <-tempChan
	right := <-tempChan

	mergeE <- Merge(left, right)
	close(tempChan)
}
