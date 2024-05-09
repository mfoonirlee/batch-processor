package processor

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	ctx := context.Background()

	var paramsArray [1000]int
	for i := 0; i < 1000; i++ {
		paramsArray[i] = i
	}

	fmt.Printf("paramsArray: %d", len(paramsArray))

	res, err := NewBatchProcessor(ctx, paramsArray, func(ctx context.Context, task any) (any, error) {
		index := task.(int)

		fmt.Println(index)

		if index%10 == 0 {
			return nil, errors.New(fmt.Sprintf("%d", index))
		}
		return index, nil

	}).Run()

	fmt.Println(res)

	// test error nums
	expectErrs := 100
	realErrs := len(strings.Split(err.Error(), ";"))
	if realErrs != expectErrs {
		t.Errorf("expect %d errors, but got %d", expectErrs, realErrs)
	}

	// test res nums
	expectResCount := 900
	if len(res.([]interface{})) != expectResCount {
		t.Errorf("expect %d results, but got %d", expectResCount, len(res.([]interface{})))
	}

}
