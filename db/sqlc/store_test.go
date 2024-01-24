package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// Check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// Check from_entry
		from_entry := result.FromEntry
		require.NotEmpty(t, from_entry)
		require.Equal(t, account1.ID, from_entry.AccountID)
		require.Equal(t, (-1)*amount, from_entry.Amount)
		require.NotZero(t, from_entry.ID)
		require.NotZero(t, from_entry.CreatedAt)

		_, err = store.GetEntry(context.Background(), from_entry.ID)

		// Check to entry
		to_entry := result.ToEntry
		require.NotEmpty(t, to_entry)
		require.Equal(t, account2.ID, to_entry.AccountID)
		require.Equal(t, amount, to_entry.Amount)
		require.NotZero(t, from_entry.ID)
		require.NotZero(t, from_entry.CreatedAt)

		_, err = store.GetEntry(context.Background(), to_entry.ID)
	}
}
