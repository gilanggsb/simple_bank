package db

import (
	"database/sql"
	"gilanggsb/simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T) Accounts {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(t.Context(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	createdAccount := CreateRandomAccount(t)
	account, err := testQueries.GetAccount(t.Context(), createdAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, account.ID, createdAccount.ID)
	require.Equal(t, account.Owner, createdAccount.Owner)
	require.Equal(t, account.Balance, createdAccount.Balance)
	require.Equal(t, account.Currency, createdAccount.Currency)

	require.WithinDuration(t, createdAccount.CreatedAt, account.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	createdAccount := CreateRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      createdAccount.ID,
		Balance: util.RandomMoney(),
	}

	updatedAccount, err := testQueries.UpdateAccount(t.Context(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)

	require.Equal(t, createdAccount.ID, updatedAccount.ID)
	require.Equal(t, createdAccount.Owner, updatedAccount.Owner)
	require.Equal(t, arg.Balance, updatedAccount.Balance)
	require.Equal(t, createdAccount.Currency, updatedAccount.Currency)

	require.WithinDuration(t, createdAccount.CreatedAt, updatedAccount.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	createdAccount := CreateRandomAccount(t)
	err := testQueries.DeleteAccount(t.Context(), createdAccount.ID)
	require.NoError(t, err)

	account, err := testQueries.GetAccount(t.Context(), createdAccount.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(t.Context(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
