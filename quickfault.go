package quickfault

import "fmt"

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