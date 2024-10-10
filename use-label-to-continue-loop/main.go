package main

import "fmt"

func main() {
	usingALabeledBreak()
	usingALabeledContinue()

}

/*
1. **Early Exit in Nested Loops**: If you’re iterating over two or more levels of loops and need to exit from both loops upon a certain condition, a labeled break helps to jump out of all nested loops at once.
Output:
Using a labeled break to exit from nested loops.
Checking character: h
Checking character: e
Checking character: l
Found 'l', breaking outer loop.
Exited loops.
*/
func usingALabeledBreak() {
	fmt.Printf("Using a labeled break to exit from nested loops.\n")
	strings := []string{"hello", "world", "golang"}

OuterLoop: // Label for the outer loop
	for _, str := range strings {
		for _, ch := range str {
			fmt.Printf("Checking character: %c\n", ch)
			if ch == 'l' {
				fmt.Println("Found 'l', breaking outer loop.")
				break OuterLoop
			}
		}
	}
	fmt.Printf("Exited loops.\n\n\n\n")
}

/*
2. **Skipping to the Next Iteration in Nested Loops**: If you’re iterating over two or more levels of loops and need to skip to the next iteration of the outer loop upon a certain condition, a labeled continue helps to jump to the next iteration of the outer loop.
Output:
Using a labeled continue to skip to the next iteration of the outer loop.
Character: h
Character: e
Found 'l', skipping to the next string.
Character: w
Character: o
Character: r
Found 'l', skipping to the next string.
Character: g
Character: o
Found 'l', skipping to the next string.
Finished iterating over strings.
*/
func usingALabeledContinue() {
	fmt.Printf("Using a labeled continue to skip to the next iteration of the outer loop.\n")
	strings := []string{"hello", "world", "golang"}

OuterLoop: // Label for the outer loop
	for _, str := range strings {
		for _, ch := range str {
			if ch == 'l' {
				fmt.Println("Found 'l', skipping to the next string.")
				continue OuterLoop
			}
			fmt.Printf("Character: %c\n", ch)
		}
	}
	fmt.Printf("Finished iterating over strings.\n\n\n\n")
}
