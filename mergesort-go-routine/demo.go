package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}

func gridGame(grid [][]int) int64 {
    row1s:=0
	for _, row := range grid[0] {
		row1s += row
	}

	minimum := row1s- grid[0][0]
	rsum := row1s- grid[0][0]
	lsum := 0

	for i:=1;i<len(grid[0]);i++{
		rsum -= grid[0][i]
		lsum += grid[0][i-1]
		fmt.Println(rsum," ",lsum)
		minimum = min(minimum, max(rsum, lsum))
	}

	return int64(minimum)

}