package quickfault

import (
	"testing"
	"math/rand"
	"fmt"
)

type FaultyTable struct {
	Risk float32
	Fault bool
	Table Table
}

type Faulter interface {
	FaultOccured() bool
}

func (t *FaultyTable) Get(k string) (string, error) {
	if rand.Float32() < t.Risk {
		t.Fault = true
		switch rand.Intn(2) {
		case 0:
			return "", nil
		case 1:
			return "", fmt.Errorf("forced fault")
		}
	}
	return t.Table.Get(k)
}

func (f FaultyTable) FaultOccured() bool {
	return f.Fault
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
	table := FaultyTable{
		Risk: 0.2,
		Table: Table{
			Data: map[string]string{
				"foo": "fizz",
				"bar": "buzz",
			},
		},
	}

	max := 100
	min := 10
	j := 0
	for i := 0; i < max && j < min; i++ {
		actual, err := example(&table)
		if table.FaultOccured() {
			if actual != "" {
				t.Logf("expected '' on fault not '%s'", actual)
				t.Fail()
			}
			if err == nil {
				t.Logf("expected err on fault, got nil")
				t.Fail()
			}
			j++
		}
		table.Fault = false
	}
	if j < min {
		t.Logf("Failed to generate enough faults, wanted %d but got %d", min, j)
		t.Fail()
	}
}