package main

import "errors"

func main() {

	// Test add
	result, err := add(1, 2)
	if err != nil {
		println("Error:", err.Error())
	} else {
		println("Addition:", result)
	}

	// Test sub
	result, err = sub(1, 2)
	if err != nil {
		println("Error:", err.Error())
	} else {
		println("Subtraction:", result)
	}

	// Test mul
	result, err = mul(1, 2)
	if err != nil {
		println("Error:", err.Error())
	} else {
		println("Multiplication:", result)
	}

	// Test div
	result, err = div(1, 2)
	if err != nil {
		println("Error:", err.Error())
	} else {
		println("Division:", result)
	}

	// Test div by zero
	result, err = div(1, 0)
	if err != nil {
		println("Error:", err.Error())
	} else {
		println("Division:", result)
	}

}

func add(a, b int) (int, error) {
	return a + b, nil
}

func sub(a, b int) (int, error) {
	return a - b, nil
}

func mul(a, b int) (int, error) {
	return a * b, nil
}

func div(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}
