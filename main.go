package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/qquiqlerr/pipeline/pipeline"
)

func Add(input map[string]interface{}) interface{} {
	a := input["a"].(float64)
	b := input["b"].(float64)
	result := a + b
	fmt.Printf("ADD: %.0f + %.0f = %.0f\n", a, b, result)
	return result
}

func Multiply(input map[string]interface{}) interface{} {
	a := input["a"].(float64)
	b := input["b"].(float64)
	result := a * b
	fmt.Printf("MULTIPLY: %.0f * %.0f = %.0f\n", a, b, result)
	return result
}

func Concat(input map[string]interface{}) interface{} {
	str1 := input["str1"]
	str2 := input["str2"]
	result := fmt.Sprintf("%v%v", str1, str2)
	fmt.Printf("CONCAT: '%v' + '%v' = '%s'\n", str1, str2, result)
	return result
}

// Conditional checks if a value is less than 100 and returns true or false.
func Conditional(input map[string]interface{}) interface{} {
	value := input["value"].(float64)
	result := value < 100
	fmt.Printf("CONDITIONAL: %.0f < 100? %t\n", value, result)
	return result
}

// Counter increments a current value by a given increment.
func Counter(input map[string]interface{}) interface{} {
	current := input["current"].(float64)
	increment := input["increment"].(float64)
	result := current + increment
	fmt.Printf("COUNTER: %.0f + %.0f = %.0f\n", current, increment, result)
	return result
}

func readPipelineFromFile(filename string) ([]pipeline.Block, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return pipeline.Parse(string(data))
}

func runDemo(title, filename string, funcs map[string]pipeline.Function) {
	fmt.Printf("Demo: %s\n", title)
	blocks, err := readPipelineFromFile(filename)
	if err != nil {
		log.Printf("Failed to read %s: %v\n", filename, err)
		return
	}

	pipeline.Run(blocks, funcs)
	fmt.Println()
}

func main() {
	// Define functions for the pipeline
	funcs := map[string]pipeline.Function{
		"add":         Add,
		"multiply":    Multiply,
		"concatenate": Concat,
		"conditional": Conditional,
		"counter":     Counter,
	}

	fmt.Println("Demonstrating non-linear pipelines from JSON files")
	fmt.Println()

	runDemo("Simple conditional branching", "pipeline1.json", funcs)

	runDemo("Parallel branching", "pipeline2.json", funcs)

	fmt.Println("Demo: a loop with condition")
	fmt.Println("WARNING: This example demonstrates a cyclic pipeline")
	fmt.Println()

	// For cyclic example, run in a separate goroutine with timeout
	blocks, err := readPipelineFromFile("pipeline_cycle.json")
	if err != nil {
		log.Printf("Failed to read pipeline_cycle.json: %v\n", err)
	} else {
		// Run with timeout
		done := make(chan bool)
		go func() {
			pipeline.Run(blocks, funcs)
			done <- true
		}()

		// Wait for completion or timeout
		select {
		case <-done:
			fmt.Println("Loop completed naturally")
		case <-time.After(5 * time.Second):
			fmt.Println("Loop stopped due to timeout")
		}
	}

	fmt.Println("All demonstrations completed")
}
