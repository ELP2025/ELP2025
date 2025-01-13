package main

import (
	"fmt"
	"math/rand"
)

func create_mat(rows int, cols int) [][]int {
	board := make([][]int, rows)
	for i := range board {
		board[i] = make([]int, cols)
	}
	return board
}

func random0_1_mat(matrix [][]int) [][]int {
	for i := range matrix {
		for j := range matrix[i] {
			if rand.Intn(10) < 2 {
				matrix[i][j] = 0
			} else {
				matrix[i][j] = 1
			}
		}
	}
	return matrix
}

func printMatrix(matrix [][]int) {
	for _, row := range matrix {
		for _, value := range row {
			fmt.Printf("%d ", value)
		}
		fmt.Println() // Nouvelle ligne après chaque ligne de la matrice
	}
}

func nb_voisin(matrix [][]int, x int, y int) int {
	//Faut tester les deux possibilités, une fonction qui se casse pas la tête qui teste si le pixel est en dehors, et une fctn qui sépare les cas de bords et les cas interieurs
	nb_voisin_1 := 0

	for i := x - 1; i < x+2; i++ {
		for j := y - 1; j < y+2; j++ {
			if i >= 0 && i <= (len(matrix)-1) && j >= 0 && j <= (len(matrix[0])-1) {
				if matrix[i][j] == 1 {
					nb_voisin_1 += 1
				}
			}
		}
	}
	if matrix[x][y] >= 1 {
		nb_voisin_1 -= 1
	}
	return nb_voisin_1
}

func main() {
	matrix := create_mat(10, 10)
	random0_1_mat(matrix)
	//matrix[9][9] = 5
	printMatrix(matrix)
	fmt.Println() // Nouvelle ligne après chaque ligne de la matrice
	fmt.Println(nb_voisin(matrix, 0, 4))
}
