# Labeling Your for Statements

By default, the break and continue keywords apply to the for loop that directly

contains them. What if you have nested for loops and want to exit or skip over

an iterator of an outer loop? Let’s look at examples.

**Example 1: Using a Labeled break**

In this example, we use a labeled break to stop iterating through the outer loop once we encounter the letter “l”:

```go
package main

import (
 "fmt"
)

func main() {
 strings := []string{"hello", "world", "golang"}
 
OuterLoop:
 for _, str := range strings {
  for _, ch := range str {
   fmt.Printf("Checking character: %c\n", ch)
   if ch == 'l' {
    fmt.Println("Found 'l', breaking outer loop.")
    break OuterLoop
   }
  }
 }
 fmt.Println("Exited loops.")
}
```

**Explanation:**

• We define an outer loop labeled OuterLoop.

• Inside the outer loop, we iterate over the characters of each string in the strings slice.

• When the character ‘l’ is found, the break OuterLoop exits both loops.

**Example 2: Using a Labeled continue**

In this example, we use a labeled continue to skip to the next iteration of the outer loop when we encounter the letter “l”:

```go
package main

import (
 "fmt"
)

func main() {
 strings := []string{"hello", "world", "golang"}
 
OuterLoop:
 for _, str := range strings {
  for _, ch := range str {
   if ch == 'l' {
    fmt.Println("Found 'l', skipping to the next string.")
    continue OuterLoop
   }
   fmt.Printf("Character: %c\n", ch)
  }
 }
 fmt.Println("Finished iterating over strings.")
}
```

**Explanation:**

• The outer loop is labeled OuterLoop.

• Inside the inner loop, when the character ‘l’ is encountered, the continue OuterLoop skips the current string and proceeds to the next string in the outer loop.

**Use Case Scenarios:**

1. **Early Exit in Nested Loops**: If you’re iterating over two or more levels of loops and need to exit from both loops upon a certain condition, a labeled break helps to jump out of all nested loops at once.

2. **Skipping to the Next Iteration in Outer Loops**: If you want to continue the outer loop while breaking only out of the inner loop, a labeled continue lets you skip over the outer loop’s remaining iterations when a certain condition in the inner loop is met.

These examples demonstrate how labeled for statements allow you to control the flow of nested loops in Go, which is especially useful in complex iteration scenarios.

## How can I do the Early Exit in Nested Loops without Using a Labeled break ?

You can achieve an early exit from nested loops in Go without using a labeled break by using one of these common alternatives:

**1. Using a Boolean Flag**

You can use a boolean flag to signal when an early exit condition is met. Once the flag is set, break from both loops.

Here’s an example:

```go
package main

import (
 "fmt"
)

func main() {
 matrix := [][]int{
  {1, 2, 3},
  {4, 5, 6},
  {7, 8, 9},
 }
 threshold := 5
 found := false

 for i := 0; i < len(matrix); i++ {
  for j := 0; j < len(matrix[i]); j++ {
   fmt.Printf("Processing matrix[%d][%d] = %d\n", i, j, matrix[i][j])
   if matrix[i][j] > threshold {
    fmt.Printf("Found a value greater than %d at matrix[%d][%d], exiting.\n", threshold, i, j)
    found = true
    break // Break inner loop
   }
  }
  if found {
   break // Break outer loop based on flag
  }
 }
 fmt.Println("Matrix processing completed.")
}
```

**Explanation:**

• A boolean flag found is used to track when the condition is met.

• After breaking from the inner loop, the outer loop checks the flag and breaks if the condition was satisfied.

**2. Using return to Exit the Function**

If the loop is part of a function and you want to exit both loops as soon as the condition is met, you can use return to terminate the function immediately.

Example:

```go
package main

import (
 "fmt"
)

func main() {
 matrix := [][]int{
  {1, 2, 3},
  {4, 5, 6},
  {7, 8, 9},
 }
 threshold := 5
 
 processMatrix(matrix, threshold)
 fmt.Println("Matrix processing completed.")
}

func processMatrix(matrix [][]int, threshold int) {
 for i := 0; i < len(matrix); i++ {
  for j := 0; j < len(matrix[i]); j++ {
   fmt.Printf("Processing matrix[%d][%d] = %d\n", i, j, matrix[i][j])
   if matrix[i][j] > threshold {
    fmt.Printf("Found a value greater than %d at matrix[%d][%d], exiting.\n", threshold, i, j)
    return // Exit the function
   }
  }
 }
}
```

**Explanation:**

• The function processMatrix returns immediately when the condition is met, stopping both loops without needing a labeled break.

**3. Using a goto Statement**

Though not commonly recommended, you can use a goto statement to jump to a label outside the loops when a condition is met.

Example:

```go
package main

import (
 "fmt"
)

func main() {
 matrix := [][]int{
  {1, 2, 3},
  {4, 5, 6},
  {7, 8, 9},
 }
 threshold := 5
 
 for i := 0; i < len(matrix); i++ {
  for j := 0; j < len(matrix[i]); j++ {
   fmt.Printf("Processing matrix[%d][%d] = %d\n", i, j, matrix[i][j])
   if matrix[i][j] > threshold {
    fmt.Printf("Found a value greater than %d at matrix[%d][%d], exiting.\n", threshold, i, j)
    goto Exit // Jump to exit
   }
  }
 }

Exit: // Label to jump to
 fmt.Println("Matrix processing completed.")
}
```

**Explanation:**

• When the condition is met, goto Exit immediately transfers control to the Exit label outside the loops.

**Comparison of Techniques:**

• **Boolean Flag**: Most common and readable approach, especially when breaking from deeply nested loops.

• return: Great when the loops are inside a function and you want to completely exit the function.

• goto: Should be used sparingly as it can make code harder to understand and maintain, but it’s a direct alternative to labeled break.

Each method provides a way to handle early exits from nested loops, depending on the specific structure and needs of your program.

## How can I do Skipping to the Next Iteration in Outer Loops without Using a Labeled continue ?

Skipping to the next iteration in outer loops without using a labeled continue in Go can be achieved using several approaches. Here are a few alternatives to handle this scenario:

**1. Using a Boolean Flag**

You can use a boolean flag to skip the remaining inner loop iterations and continue with the next iteration of the outer loop.

**Example:**

```go
package main

import (
 "fmt"
)

func main() {
 strings := []string{"hello", "world", "golang"}
 
 for _, str := range strings {
  skip := false
  for _, ch := range str {
   if ch == 'l' {
    fmt.Println("Found 'l', skipping to the next string.")
    skip = true
    break // Break inner loop
   }
   fmt.Printf("Character: %c\n", ch)
  }
  if skip {
   continue // Skip the rest of the outer loop iteration
  }
  fmt.Println("Finished processing the string.")
 }
 fmt.Println("Done with all strings.")
}
```

**Explanation:**

• The flag skip is used to signal when to skip the rest of the outer loop iteration.

• When the condition (ch == 'l') is met, the inner loop breaks, and the continue statement in the outer loop is triggered by checking the skip flag.

**2. Using return in a Function**

If the logic is inside a function and you want to continue processing in the outer loop, you can encapsulate the inner loop logic inside a separate function. If a condition is met, you can return early and move on to the next iteration.

**Example:**

```go
package main

import (
 "fmt"
)

func main() {
 strings := []string{"hello", "world", "golang"}
 
 for _, str := range strings {
  if processString(str) {
   fmt.Println("Skipping to the next string.")
   continue
  }
  fmt.Println("Finished processing the string.")
 }
 fmt.Println("Done with all strings.")
}

func processString(str string) bool {
 for _, ch := range str {
  if ch == 'l' {
   fmt.Println("Found 'l'.")
   return true // Return true to indicate skipping the outer loop
  }
  fmt.Printf("Character: %c\n", ch)
 }
 return false // Continue with the outer loop
}
```

**Explanation:**

• The processString function returns true if the condition is met, signaling that the outer loop should skip to the next iteration.

• If false is returned, the outer loop continues processing.

**3. Using goto to Skip to the Next Outer Loop Iteration**

Another alternative, though less common, is using goto to directly jump to a label that continues the outer loop.

**Example:**

```go
package main

import (
 "fmt"
)

func main() {
 strings := []string{"hello", "world", "golang"}

 for _, str := range strings {
  for _, ch := range str {
   if ch == 'l' {
    fmt.Println("Found 'l', skipping to the next string.")
    goto NextString // Jump to the label
   }
   fmt.Printf("Character: %c\n", ch)
  }
  fmt.Println("Finished processing the string.")
 NextString:
 }
 fmt.Println("Done with all strings.")
}
```

**Explanation:**

• The goto statement is used to jump to the label NextString, which skips the remaining iterations of the inner loop and continues with the next string in the outer loop.

**Comparison of Techniques:**

• **Boolean Flag**: Clean and readable, avoids jumps and works well for simple cases.

• **Using** return: Helpful when you want to structure your logic into functions and make code modular.

• **goto**: Should be used sparingly; it can be useful for skipping to specific parts of the code but may reduce code readability.

Each of these methods helps skip to the next iteration of the outer loop without needing a labeled continue, depending on your use case and coding style preferences.
