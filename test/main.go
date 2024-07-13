package main

import (
	"flim"
	"flim/common"
	"fmt"
)

func main() {
	expr, err := flim.ParseFile("./test.flim")

	if err != nil {
		panic(err)
	}

	exportedValues := map[string]interface{}{
		"hostname":  "100.100.100.102",
		"base_port": 8000,
	}

	handlers := map[string]common.HandlerFunc{
		// We can use a handler to enforce required keys
		"inventory": func(data interface{}) (interface{}, error) {
			requiredKeys := []string{"items"}

			for _, key := range requiredKeys {
				if _, ok := data.(map[string]interface{})[key]; !ok {
					return nil, fmt.Errorf("inventory missing required key `%s'", key)
				}
			}

			return data, nil
		},

		// We can use a handler to provide default values
		"item": func(data interface{}) (interface{}, error) {
			out := map[string]interface{}{
				"timeout": 30,
			}

			for key, val := range data.(map[string]interface{}) {
				out[key] = val
			}

			return out, nil
		},

		// We can use a handler to dynamically replace values, a la variables
		"from": func(data interface{}) (interface{}, error) {
			val, exists := exportedValues[data.(string)]

			if !exists {
				return nil, fmt.Errorf("unknown variable `%s'", data.(string))
			}

			return val, nil
		},

		// We can do other weird things
		"add": func(data interface{}) (interface{}, error) {
			inputs := data.([]interface{})
			a := inputs[0].(int)
			b := inputs[1].(int64)

			return a + int(b), nil
		},

		"square": func(data interface{}) (interface{}, error) {
			input := data.(int64)

			return input * input, nil
		},
	}

	fmt.Println(expr.ToString())

	output, err := expr.Evaluate(handlers)

	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println(output)

	serialized, err := common.Serialize(expr, false, 4)

	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println(serialized)

	minified, err := common.Minify(expr)

	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println(minified)
}
