package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/bagashiz/Simple-Bank/db/mock"
	db "github.com/bagashiz/Simple-Bank/db/sqlc"
	"github.com/bagashiz/Simple-Bank/token"
	"github.com/bagashiz/Simple-Bank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

// TestGetAccountAPI tests the getAccount API using mock database.
func TestGetAccountAPI(t *testing.T) {
  user, _ := randomUser(t)
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		accountID     int64
    setupAuth func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		// case 1: get account successfully.
		{
			name:      "OK",
			accountID: account.ID,
      setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
        addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
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
		// case 2: account not found.
		{
			name:      "NotFound",
			accountID: account.ID,
      setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
        addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
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
		// case 3: internal server error.
		{
			name:      "InternalError",
			accountID: account.ID,
      setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
        addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
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
		// case 4: bad request with non-number account id.
		{
			name:      "InvalidID",
			accountID: 0,
      setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
        addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
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
		// case 5: unauthorized user.
		{
			name:      "UnauthorizedUser",
			accountID: account.ID,
      setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
        addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "unauthorized_user", time.Minute)
      },
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		// case 6: no authorization.
		{
			name:      "NoAuthorization",
			accountID: account.ID,
      setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
      },
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(),gomock.Any()).
          Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			// build stubs
			tc.buildStubs(store)

			// start test server and send request
			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

      tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

// TestCreateAccountAPI tests the createAccount API using mock database.
func TestCreateAccountAPI(t *testing.T) {
  user, _ := randomUser(t)
	account := randomAccount(user.Username)
  
	testCases := []struct {
		name          string
    body          gin.H
    setupAuth func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		// case 1: create account successfully.
    {
      name: "OK",
      body: gin.H{
        "currency": account.Currency,
      },
      setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
        addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
      },
      buildStubs: func(store *mockdb.MockStore) {
        arg := db.CreateAccountParams{
          Owner:      account.Owner,
          Currency:   account.Currency,
          Balance:    0,
        }

        store.EXPECT().
          CreateAccount(gomock.Any(), gomock.Eq(arg)).
          Times(1).
          Return(account, nil)
      },
      checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
        require.Equal(t, http.StatusOK, recorder.Code)
        requireBodyMatchAccount(t, recorder.Body, account)
      },
    },
		// case 2: no authorization.
    {
      name: "NoAuthorization",
      body: gin.H{
        "currency": account.Currency,
      },
      setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
      },
      buildStubs: func(store *mockdb.MockStore) {
        store.EXPECT().
          CreateAccount(gomock.Any(), gomock.Any()).
          Times(0)
      },
      checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
        require.Equal(t, http.StatusUnauthorized, recorder.Code)
      },
    },
		// case 3: internal sever error.
    {
      name: "InternalError",
      body: gin.H{
        "currency": account.Currency,
      },
      setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
        addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
      },
      buildStubs: func(store *mockdb.MockStore) {
        store.EXPECT().
          CreateAccount(gomock.Any(), gomock.Any()).
          Times(1).
          Return(db.Account{}, sql.ErrConnDone)
      },
      checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
        require.Equal(t, http.StatusInternalServerError, recorder.Code)
      },
    },
		// case 4: using invalid currency.
    {
      name: "InvalidCurrency",
      body: gin.H{
        "currency": "invalid",
      },
      setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
        addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
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
  
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			// build stubs
			tc.buildStubs(store)

			// start test server and send request
			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()

      // Marshal body data to JSON
      data, err := json.Marshal(tc.body)
			require.NoError(t, err)

      url := "/accounts"
      request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
      
      tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})
	}
}

// randomAccount returns a random generated account using util package.
func randomAccount(owner string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

// requireBodyMatchAccount checks if the response body matches the given account.
func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
