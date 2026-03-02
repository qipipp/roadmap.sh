package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

func getInt() (int, error) {
	var s string
	if _, err := fmt.Scan(&s); err != nil {
		return 0, err
	}
	res, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func main() {

	fmt.Printf(`Welcome to the Number Guessing Game!
I'm thinking of a number between 1 and 100.
You have 5 chances to guess the correct number.

Please select the difficulty level:
1. Easy (10 chances)
2. Medium (5 chances)
3. Hard (3 chances)

`)
	var limit int
	for {
		fmt.Printf("Enter your choice:")
		difficulty, err := getInt()
		if err != nil || difficulty < 1 || difficulty > 3 {
			fmt.Printf("Please enter a number (1~3).\n")
		} else {
			var s string
			switch difficulty {
			case 1:
				limit = 10
				s = "Easy"
			case 2:
				limit = 5
				s = "Medium"
			case 3:
				limit = 3
				s = "Hard"
			default:
				fmt.Printf("Please enter a number (1~3).\n")
				continue
			}
			fmt.Printf("Great! You have selected the %v difficulty level.\n", s)
			break
		}
	}
	ans := rand.Intn(100) + 1
	cnt := 0
	for {
		fmt.Printf("Enter your guess:")
		x, err := getInt()
		//fmt.Printf("your input:%v %v\n", x, err)
		if err != nil || x < 1 || x > 100 {
			fmt.Printf("Please enter a number (1~100).\n")
			continue
		}
		cnt += 1
		if ans < x {
			fmt.Printf("Incorrect! The number is less than %v.\n", x)
		} else if ans > x {
			fmt.Printf("Incorrect! The number is greater than %v.\n", x)
		} else {
			fmt.Printf("Congratulations! You guessed the correct number in %v attempts\n.", cnt)
			break
		}

		if cnt == limit {
			fmt.Printf("You lost!")
			break
		}
	}
}
