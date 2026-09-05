package pointer

import (
	"fmt"
	"testing"
)

func TestWallet(t *testing.T) {

	t.Run("Test Case for Deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		fmt.Printf("address of balance in test is %p \n", &wallet.balance)

		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("Test Case for Withdraw", func(t *testing.T) {
		wallet := Wallet{balance: 30}
		err := wallet.Withdraw(Bitcoin(10))

		assertNoError(t, err)
		assertBalance(t, wallet, Bitcoin(20))
	})

	t.Run("Test case for withdraw insufficient funds", func(t *testing.T) {
		startingBalance := Bitcoin(20)
		wallet := Wallet{startingBalance}

		err := wallet.Withdraw(Bitcoin(100))

		assertError(t, err, ErrInsufficientFunds)
		assertBalance(t, wallet, startingBalance)

	})

}

func assertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatal("got an error but didn't want one")
	}
}

func assertError(t testing.TB, got error, want error) {
	t.Helper()
	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}
	
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertBalance(t testing.TB, wallet Wallet, want Bitcoin) {
	got := wallet.Balance()

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}