package processor

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	ctx := context.Background()

	var paramsArray [1000]int
	for i := 0; i < 1000; i++ {
		paramsArray[i] = i
	}

	// paramsArray := make([]int, 0)

	// for i := 1; i <= 1000; i++ {
	// 	paramsArray = append(paramsArray, i)
	// }

	fmt.Printf("paramsArray: %d", len(paramsArray))

	res, err := NewBatchProcessor(ctx, paramsArray, func(ctx context.Context, task any) (any, error) {
		index := task.(int)

		fmt.Println(index)

		if index%10 == 0 {
			return nil, errors.New("taskIndex: " + string(rune(index)))
		}
		return index, nil

	}).Run()

	fmt.Println(res)

	fmt.Println(err)
}
