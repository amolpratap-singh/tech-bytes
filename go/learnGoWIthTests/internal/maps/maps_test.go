package maps

import "testing"

func TestSearch(t *testing.T) {
	t.Run("Test Case for Search function", func(t *testing.T) {
		dictionary := map[string]string{"sentinel": "lookout, watch, watcher"}
		got := Search(dictionary, "sentinel")
		want := "lookout, watch, watcher"

		assertStrings(t, got, want)
	})
}

func TestCustomSearch(t *testing.T) {
	dictionary := Dictionary{"idiomatic": "using, containing, or denoting expressions that are natural to a native speaker."}
	t.Run("Test case for custom dictionary with known word", func(t *testing.T) {
		got, _ := dictionary.CustomSearch("idiomatic")
		want := "using, containing, or denoting expressions that are natural to a native speaker."

		assertStrings(t, got, want)
	})

	t.Run("Test case for custom dicitioary with unknown word", func(t *testing.T) {
		_, got := dictionary.CustomSearch("idioms")

		if got == nil {
			t.Fatal("expected to get an error.")
		}
		assertError(t, got, ErrNotFound)
	})
}

func TestAdd(t *testing.T) {
	t.Run("Test case to add key in dicitionary", func(t *testing.T) {
		dictionary := Dictionary{}
		dictionary.Add("idioms", "a common code pattern or practice specific to a programming language")

		got, err := dictionary.CustomSearch("idioms")
		want := "a common code pattern or practice specific to a programming language"

		if err != nil {
			t.Fatal("Should find added word:", err)
		}

		assertStrings(t, got, want)

	})
}

func assertStrings(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}

}

func assertError(t testing.TB, got, want error) {
	t.Helper()
	if got != want {
		t.Errorf("got error %q want %q", got, want)
	}
}