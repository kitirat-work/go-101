package main

func main() {
	// Test prefixer
	prefix := prefixer("Hello, ")
	result := prefix("World!")
	println(result)

	// Chain prefixers
	result = prefixer("Hello, ")("World!")
	println(result)

}

func prefixer(prefix string) func(string) string {
	return func(s string) string {
		return prefix + s
	}
}
