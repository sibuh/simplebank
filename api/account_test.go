package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	mockdb "exercise/simplebank/db/mock"
	db "exercise/simplebank/db/sqlc"
	"exercise/simplebank/token"
	"exercise/simplebank/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountApi(t *testing.T) {
	user := randomUser()
	account := randomAccount(user.Username)
	testCases := []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{name: "ok",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, time.Minute, authType, user.Username)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{name: "not found",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, time.Minute, authType, user.Username)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)

			},
		},
		{name: "internal error",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, time.Minute, authType, user.Username)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},
		{name: "invalidID",
			accountID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, time.Minute, authType, user.Username)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},
	}
	for _, tc := range testCases {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		store := mockdb.NewMockStore(ctrl)
		tc.buildStubs(store)
		server := newTestServer(t, store)
		recorder := httptest.NewRecorder()
		url := fmt.Sprintf("/accounts/%d", tc.accountID)
		request, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)
		tc.setupAuth(t, request, server.tokenMaker)
		server.router.ServeHTTP(recorder, request)
		tc.checkResponse(t, recorder)
	}

}
func randomAccount(owner string) db.Account {
	return db.Account{
		ID:       util.RandInt(1000, 1),
		Owner:    owner,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
func randomUser() db.User {
	return db.User{
		Username: util.RandString(6),
		FullName: util.RandomOwner(),
	}
}
func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}

func TestCreateAccount(t *testing.T) {
	user := randomUser()
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		currency      string
		setup         func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "ok",
			currency: account.Currency,
			setup: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, time.Minute, authType, user.Username)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{Owner: account.Owner, Balance: 0, Currency: account.Currency}
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:     "internal server error ",
			currency: account.Currency,
			setup: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, time.Minute, authType, user.Username)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{Owner: account.Owner, Balance: 0, Currency: account.Currency}
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:     "invalid currency",
			currency: "asdfgh",
			setup: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, time.Minute, authType, user.Username)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			server := newTestServer(t, store)
			arg := createAccountRequest{Currency: tc.currency}
			byteCurrency, err := json.Marshal(arg)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(byteCurrency))
			require.NoError(t, err)
			tc.setup(t, request, server.tokenMaker)
			recorder := httptest.NewRecorder()
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}
func listBodyMatchAccount(t *testing.T, body *bytes.Buffer, account []db.Account) {

	var gotAccount []db.Account
	err := json.NewDecoder(body).Decode(&gotAccount)

	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
func TestListAccount(t *testing.T) {
	randOffset := int32(1)
	randLimit := int32(util.RandInt(10, 5))
	user := randomUser()
	account := randomAccount(user.Username)
	testCase := []struct {
		name          string
		offset        int32
		limit         int32
		setup         func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "ok",
			offset: randOffset,
			limit:  randLimit,
			setup: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, time.Minute, authType, user.Username)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListAccountsParams{
					Owner:  user.Username,
					Limit:  int32(randLimit),
					Offset: int32(randOffset - 1),
				}
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).Times(1).
					Return([]db.Account{account}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				listBodyMatchAccount(t, recorder.Body, []db.Account{account})
			},
		},
	}
	for _, tc := range testCase {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		store := mockdb.NewMockStore(ctrl)
		tc.buildStubs(store)
		server := newTestServer(t, store)
		queryString := fmt.Sprintf("/accounts/?page_id=%d&page_size=%d", tc.offset, tc.limit)
		request, err := http.NewRequest(http.MethodGet, queryString, nil)
		tc.setup(t, request, server.tokenMaker)
		require.NoError(t, err)
		recorder := httptest.NewRecorder()
		server.router.ServeHTTP(recorder, request)
		tc.checkResponse(t, recorder)
	}

}
