package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"module-package/calculator"
)

func main() {
	godotenv.Load()
	fmt.Println(os.Getenv("GO_ENV"))
	fmt.Println(calculator.Offset)
	fmt.Println(calculator.Sum(1, 2))
	fmt.Println(calculator.Multiply(1, 2))
}
