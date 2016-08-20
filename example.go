package quickfault

import "fmt"

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
	foo, _ := GetFooValueFromTable(t)
	bar, _ := GetBarValueFromTable(t)
	return foo+bar, nil
}