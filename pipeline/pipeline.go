package pipeline

import (
	"encoding/json"
	"sync"
)

// Parse takes a JSON string and unmarshals it into a slice of Block structs.
func Parse(input string) ([]Block, error) {
	var blocks []Block
	if err := json.Unmarshal([]byte(input), &blocks); err != nil {
		return nil, err
	}
	return blocks, nil
}

// Run executes the pipeline blocks in parallel, resolving dependencies and executing functions.
func Run(blocks []Block, funcs map[string]Function) {
	blockMap := make(map[string]Block)
	results := &sync.Map{}

	// Create a map of blocks
	for _, b := range blocks {
		blockMap[b.Name] = b
	}

	executed := &sync.Map{}

	// Find start blocks with no dependencies
	var startBlocks []string
	for _, b := range blocks {
		if hasNoDependencies(b) {
			startBlocks = append(startBlocks, b.Name)
		}
	}

	// Execute blocks recursively
	var wg sync.WaitGroup
	for _, blockName := range startBlocks {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			executeBlock(name, blockMap, funcs, results, executed)
		}(blockName)
	}
	wg.Wait()
}

// executeBlock executes a single block and its outputs recursively.
func executeBlock(blockName string, blockMap map[string]Block, funcs map[string]Function, results *sync.Map, executed *sync.Map) {
	if _, exists := executed.Load(blockName); exists {
		return
	}

	block := blockMap[blockName]

	// Resolve inputs
	input := resolveInputs(block.Input, results)

	// Execute the function
	result := funcs[block.Function](input)

	// Save the result
	results.Store(block.Name, result)
	executed.Store(block.Name, true)

	// Process conditional outputs
	if block.Function == "conditional" {
		// For conditional blocks, choose the execution branch
		if boolResult, ok := result.(bool); ok {
			if boolResult && len(block.Output) > 0 {
				// true - execute first output
				executeBlock(block.Output[0], blockMap, funcs, results, executed)
			} else if !boolResult && len(block.Output) > 1 {
				// false - execute second output
				executeBlock(block.Output[1], blockMap, funcs, results, executed)
			}
		}
	} else {
		// For regular blocks, execute all outputs asynchronously
		var wg sync.WaitGroup
		for _, outputBlock := range block.Output {
			wg.Add(1)
			go func(outBlock string) {
				defer wg.Done()
				executeBlock(outBlock, blockMap, funcs, results, executed)
			}(outputBlock)
		}
		wg.Wait()
	}
}

// hasNoDependencies checks if a block has no dependencies.
func hasNoDependencies(block Block) bool {
	for _, v := range block.Input {
		if _, ok := v.(string); ok {
			// If there is a string value, it may be a dependency
			// For simplicity, we consider there are no dependencies only if all values are not strings
			// or strings are not block names
			return false
		}
	}
	return true
}

// resolveInputs resolves inputs without blocking
func resolveInputs(input map[string]any, results *sync.Map) map[string]any {
	resolved := make(map[string]any)
	for k, v := range input {
		if strVal, ok := v.(string); ok {
			if val, exists := results.Load(strVal); exists {
				resolved[k] = val
			} else {
				// If dependency is not found, use the string as is
				resolved[k] = v
			}
		} else {
			resolved[k] = v
		}
	}
	return resolved
}
