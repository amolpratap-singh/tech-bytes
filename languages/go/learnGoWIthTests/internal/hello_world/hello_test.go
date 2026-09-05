package hello

import "testing"

func TestHello(t *testing.T) {

	t.Run("say hello to people", func(t *testing.T) {
		got := Hello("Chris", "")
		want := "Hello, Chris"
		assertCorrectMessage(t, got, want)
	})

	t.Run("say 'Hello World!' when an empty string is supplied", func(t *testing.T) {
		got := Hello("", "")
		want := "Hello, World"
		assertCorrectMessage(t, got, want)
	})

	t.Run("in Spanish", func(t *testing.T) {
		got := Hello("Elodie", "Spanish")
		want := "Hola, Elodie"
		assertCorrectMessage(t, got, want)
	})

	t.Run("in French", func(t *testing.T) {
		got := Hello("Elodie", "French")
		want := "Bonjour, Elodie"
		assertCorrectMessage(t, got, want)
	})
}


/***
When we have more than one argument of the same type (in our case two strings) rather 
than having (got string, want string) we can shorten it to (got, want string).
**/
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