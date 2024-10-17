# Exercises

These exercises test your knowledge of functions in Go and how to use them.
Solutions are available in the exercise_solutions directory in the Chapter 5 repository.

1. The simple calculator program doesn’t handle one error case: division by zero.
Change the function signature for the math operations to return both an int and
an error. In the div function, if the divisor is 0, return errors.New("division
by zero") for the error. In all other cases, return nil. Adjust the main function
to check for this error.
2. Write a function called fileLen that has an input parameter of type string and
returns an int and an error. The function takes in a filename and returns the
number of bytes in the file. If there is an error reading the file, return the error.
Use defer to make sure the file is closed properly.
3. Write a function called prefixer that has an input parameter of type string
and returns a function that has an input parameter of type string and returns a
string. The returned function should prefix its input with the string passed into
prefixer. Use the following main function to test prefixer: