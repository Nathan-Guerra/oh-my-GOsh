package autocomplete

import "testing"

func TestLargestPrefix(t *testing.T) {
	ac := GetCommandAutocompleter()
	ac.SetBuiltins([]string{"echo", "exitt", "echo2", "exitter", "exitting"})
	ac.EagerLoad()

	ac.Match("e")
	p := ac.LargestCommonPrefix("e")
	t.Logf("prefix of 'e': %s", p)
	if len(p) != 0 {
		t.Errorf("Expecting 0 byte prefix for the search 'e', found %d", len(p))
	}

	ac.Match("ec")
	p = ac.LargestCommonPrefix("ec")
	t.Logf("prefix of 'ec': %s", p)
	if len(p) != 2 {
		t.Errorf("Expecting 2 byte prefix for the search 'ec', found %d", len(p))
	}

	ac.Match("ex")
	p = ac.LargestCommonPrefix("ex")
	t.Logf("prefix of 'ex': %s", p)
	if len(p) != 3 {
		t.Errorf("Expecting 2 byte prefix for the search 'ex', found %d", len(p))
	}
}
