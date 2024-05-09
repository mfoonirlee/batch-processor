package processor

import (
	"reflect"
)

func SplitByLength[T any](sourceArray []T, length int) (targetArray [][]T) {
	if len(sourceArray) == 0 || length == 0 {
		return nil
	}

	// split by length
	for i := 0; i < len(sourceArray); i += length {
		newIndex := i + length
		if newIndex > len(sourceArray) {
			newIndex = len(sourceArray)
		}
		targetArray = append(targetArray, sourceArray[i:newIndex])
	}

	return
}

// transform interface array like type to any slice
func transformInterfaceToAnySlice(input interface{}) (outputSlice []any) {
	inputSlice := reflect.ValueOf(input)

	// type check
	if inputSlice.Kind() != reflect.Slice && inputSlice.Kind() != reflect.Array {
		return nil
	}

	outputSlice = make([]interface{}, inputSlice.Len())
	for i := 0; i < inputSlice.Len(); i++ {
		outputSlice[i] = inputSlice.Index(i).Interface()
	}

	return
}
