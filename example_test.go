package quickfault

import "testing"

func TestExample_Positive(t *testing.T) {
	actual, _ := example()
	if actual != "fizzbuzz" {
		t.Logf("wanted 'fizzbuzz' got '%s'", actual)
		t.Fail()
	}
}