package logger

import "fmt"

func Log(params map[string]any) {
	for key, value := range params {
		fmt.Printf("%s: %v\n", key, value)
	}
}
