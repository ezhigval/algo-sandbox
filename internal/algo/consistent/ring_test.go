package consistent

import "testing"

func TestRing_LookupStable(t *testing.T) {
	r := NewRing([]string{"node-a", "node-b", "node-c"}, 50)
	n1 := r.Lookup("user-42")
	n2 := r.Lookup("user-42")
	if n1 != n2 {
		t.Fatalf("expected stable mapping, got %s vs %s", n1, n2)
	}
}
