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

	t.Run("new word", func(t *testing.T) {
		dictionary := Dictionary{}
		word := "test"
		definition := "this is just a test"
		err := dictionary.Add(word, definition)

		assertError(t, err, nil)
		assertDefinition(t, dictionary, word, definition)
	})

	t.Run("existing word", func(t *testing.T) {
		word := "test"
		definition := "this is just a test"
		dictionary := Dictionary{word: definition}
		err := dictionary.Add(word, "new test")

		assertError(t, err, ErrWordExists)
		assertDefinition(t, dictionary, word, definition)
	})
}

func assertDefinition(t testing.TB, dictionary Dictionary, word, definition string) {
	t.Helper()

	got, err := dictionary.CustomSearch(word)
	if err != nil {
		t.Fatal("should find added word:", err)
	}
	assertStrings(t, got, definition)
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