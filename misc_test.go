package checkxml

import (
	// "fmt"
	"testing"
)

func TestHasTags(t *testing.T) {
	// fmt.Println("===================== TestHasTags ...")
	result := []string{
		"this.Is.a.test",
		"of",
		"some.Dummy.values",
	}

	check := []string{
		"this.Is.a.test",
		"of",
	}
	if ok, v := HasTags(result, check...); !ok {
		t.Fatalf("not true: %v", v)
	}

	check = []string{
		"this.Is",
		"some.values",
	}
	if ok, v := HasTags(result, check...); ok {
		t.Fatalf("true: %v", v)
	} else if len(v) != 2 {
		t.Fatalf("result has len %d: %v", len(v), v)
	}
}
