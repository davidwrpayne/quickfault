package quickfault

import (
	"testing"
	"math/rand"
	"fmt"
)

type FaultyTable struct {
	Risk float32
	fault bool
	Table Table
}

func (t *FaultyTable) Get(k string) (string, error) {
	if rand.Float32() < t.Risk {
		t.fault = true
		switch rand.Intn(2) {
		case 0:
			return "", nil
		case 1:
			return "", fmt.Errorf("forced fault")
		}
	}
	return t.Table.Get(k)
}

func (f FaultyTable) Fault() bool {
	return f.fault
}

func (f *FaultyTable) Reset() {
	f.fault = false
}

func TestExample_Positive(t *testing.T) {
	table := Table{
		Data: map[string]string{
			"foo": "fizz",
			"bar": "buzz",
		},
	}
	actual, _ := example(&table)
	if actual != "fizzbuzz" {
		t.Logf("wanted 'fizzbuzz' got '%s'", actual)
		t.Fail()
	}
}

func TestExample_Negative(t *testing.T) {
	// wrap the Table in a FaultyTable
	table := &FaultyTable{
		Risk: 0.2,
		Table: Table{
			Data: map[string]string{
				"foo": "fizz",
				"bar": "buzz",
			},
		},
	}

	Check(func() error {
		actual, err := example(table)

		// property 1 - on failure, expect empty output
		if actual != "" {
			return fmt.Errorf("fault should result in empty output")
		}

		// property 2 - on failure, expect an err
		if err == nil {
			return fmt.Errorf("fault should result in err")
		}

		return nil
	}, table, 10, 100)
}

type Q interface {
	Reset()
	Fault() bool
}

func Check(h func() error, q Q, min int, max int) error {
	j := 0
	for i := 0; i < max && j < min; i++ {
		q.Reset()
		err := h()
		if q.Fault() {
			if err != nil {
				return fmt.Errorf("test failure: %s", err)
			}
			j++
		}
	}
	if j < min {
		return fmt.Errorf("failed to generate enough test cases, needed %d got %d", min, j)
	}
	return nil
}