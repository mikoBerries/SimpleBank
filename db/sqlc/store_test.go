package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(TestConn)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	amount := int64(1000000)

	//channel to store 5 go routine running
	transferCh := make(chan TransferTxResult)
	errorCh := make(chan error)

	for i := 0; i < 5; i++ {
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

	//check all chaneel produced from go routine
	for i := 0; i < 5; i++ {
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
		require.Equal(t, account2, tf.Transfer.ToAccountID)      //to account id eq check
		require.Equal(t, amount, tf.Transfer.Amount)             //ammount eq check

		//check actual transfer data from DB
		//if Transfer record are found will be return a transfer struct and a nil error
		_, err := store.queires.GetTransfer(context.Background(), tf.Transfer.ID)
		require.NoError(t, err)

		//check FromEntry
		fromEntry := tf.FromEntry
		require.NotEmpty(t, fromEntry)              //entry struct empty check
		require.NotZero(t, fromEntry.ID)            //from id empty check
		require.NotZero(t, fromEntry.CreatedAt)     //createat null check
		require.Equal(t, -amount, fromEntry.Amount) //-amount eq check
		require.Equal(t, account1, fromEntry.ID)    //account id eq check

		//check actual entry data from DB
		_, err = store.queires.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		//check FromEntry
		toEntry := tf.ToEntry
		require.NotEmpty(t, toEntry)             //entry struct empty check
		require.NotZero(t, toEntry.ID)           //from id empty check
		require.NotZero(t, toEntry.CreatedAt)    //createat null check
		require.Equal(t, amount, toEntry.Amount) //-amount eq check
		require.Equal(t, account1, toEntry.ID)   //account id eq check

		//check actual entry data from DB
		_, err = store.queires.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)
	}

}
