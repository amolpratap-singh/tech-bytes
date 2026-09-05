package iteration

import "testing"

func TestRepeat(t *testing.T) {

	t.Run("Default repeat", func(t *testing.T) {
		got := Repeat("a")
		want := "aaa"

		assertCorrectMessage(t, got, want)
	})

	t.Run("repeat 0 times", func(t *testing.T) {
		got := Repeat("a", 0)
		want := ""
		
		assertCorrectMessage(t, got, want)
	})

	t.Run("repeat 5 times", func(t *testing.T) {
		got := Repeat("a", 5)
		want := "aaaaa"

		assertCorrectMessage(t, got, want)
	})
	
}

func TestCharacterCount(t *testing.T) {
	
	t.Run("Character 'a' count in word malayalam ", func(t *testing.T) {
		got := CharacterCount("malayalam", "a")
		want := 4

		if got != want {
			t.Errorf("Test Case for character count where got %q but want %q", got, want)
		}
	})

	t.Run("Empty character in word", func(t *testing.T) {
		got := CharacterCount("angular", "")
		want := 0

		if got != want {
			t.Errorf("Test Case for empty character count where %q but want %q", got, want)
		}

	})
}

func BenchmarkRepeat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Repeat("a", 1)
	}
}

func assertCorrectMessage(t testing.TB, got, want string) {
	/***
	 when it fails, the line number reported will be in our function call 
	 rather than inside our test helper. This will help other developers track down problems 
	 more easily.
	***/
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
