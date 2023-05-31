package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

//TestTransferTx Testing transfer transaction
func TestTransferTx(t *testing.T) {
	store := NewStore(TestConn)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	amount := int64(10)

	//channel to store 5 go routine running
	n := 5
	transferCh := make(chan TransferTxResult, n)
	errorCh := make(chan error, n)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			//fill every go routine result to channel
			transferCh <- result
			errorCh <- err
		}()
	}

	// check results
	existed := make(map[int]bool)

	//check all chanel produced from go routine
	for i := 0; i < n; i++ {
		//check err first
		er := <-errorCh
		require.NoError(t, er)

		//check result
		tf := <-transferCh
		//check transfer
		require.NotEmpty(t, tf.Transfer)                         //transfer struct empty check
		require.NotZero(t, tf.Transfer.ID)                       //transfer id null check
		require.NotZero(t, tf.Transfer.CreatedAt)                //create at null check
		require.Equal(t, account1.ID, tf.Transfer.FromAccountID) //from account id eq check
		require.Equal(t, account2.ID, tf.Transfer.ToAccountID)   //to account id eq check
		require.Equal(t, amount, tf.Transfer.Amount)             //ammount eq check

		//check actual transfer data from DB
		//if Transfer record are found will be return a transfer struct and a nil error
		_, err := store.queires.GetTransfer(context.Background(), tf.Transfer.ID)
		require.NoError(t, err)

		//check FromEntry
		fromEntry := tf.FromEntry
		require.NotEmpty(t, fromEntry)                     //entry struct empty check
		require.NotZero(t, fromEntry.ID)                   //from id empty check
		require.NotZero(t, fromEntry.CreatedAt)            //createat null check
		require.Equal(t, -amount, fromEntry.Amount)        //-amount eq check
		require.Equal(t, account1.ID, fromEntry.AccountID) //account id eq check

		//check actual entry data from DB
		_, err = store.queires.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		//check toEntry
		toEntry := tf.ToEntry
		require.NotEmpty(t, toEntry)                     //entry struct empty check
		require.NotZero(t, toEntry.ID)                   //from id empty check
		require.NotZero(t, toEntry.CreatedAt)            //createat null check
		require.Equal(t, amount, toEntry.Amount)         //-amount eq check
		require.Equal(t, account2.ID, toEntry.AccountID) //account id eq check

		//check actual entry data from DB
		_, err = store.queires.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//check balance in every account
		fromAcc := tf.FromAccount
		require.NoError(t, err)
		require.Equal(t, account1.ID, fromAcc.ID)

		toAcc := tf.ToAccount
		require.NoError(t, err)
		require.Equal(t, account2.ID, toAcc.ID)

		// check every transfered balances
		fmt.Println(">> tx:", fromAcc.Balance, toAcc.Balance)

		diff1 := account1.Balance - fromAcc.Balance
		diff2 := toAcc.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedAccount1, err := store.queires.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.queires.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)

}

//TestTransferTxDeadlock for testing deadlock in database
func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(TestConn)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	n := 10
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balance
	updatedAccount1, err := store.queires.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.queires.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
