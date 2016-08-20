package quickfault_test

import (
	"testing"
	"math/rand"
	"fmt"
	"github.com/brentdrich/quickfault"
)

type Getter interface {
	Get(string) (string, error)
}

type Table struct {
	Fail bool
	Data map[string]string
}

// Get is part of a contrived key value store that has the
// following properties:
//  1. When there is a network failure, it errors out.
//  2. When a value can't be found for a key, it returns "".
//  3. When a value is found for a key, it returns the value.
func (t *Table) Get(k string) (string, error) {
	if t.Fail {
		return "", fmt.Errorf("failed to connect.")
	}
	if v, ok := t.Data[k]; ok {
		return v, nil
	} else {
		return "", nil
	}
}

func GetFooValueFromTable(c Getter) (string, error) {
	return c.Get("foo")
}

func GetBarValueFromTable(c Getter) (string, error) {
	return c.Get("bar")
}

func example(t Getter) (string, error) {
	foo, err := GetFooValueFromTable(t)
	if err != nil {
		return "", err
	}
	if foo == "" {
		return "", fmt.Errorf("couldn't access foo")
	}

	bar, err := GetBarValueFromTable(t)
	if err != nil {
		return "", err
	}
	if bar == "" {
		return "", fmt.Errorf("couldn't access bar")
	}

	return foo+bar, nil
}

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

	quickfault.Check(func() error {
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