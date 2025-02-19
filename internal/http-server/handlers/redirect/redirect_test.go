package redirect_test

import (
	"errors"
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-service/internal/lib/api"
	slogdiscard "url-service/internal/lib/logger"
	"url-service/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"url-service/internal/http-server/handlers/redirect"
	"url-service/internal/http-server/handlers/redirect/mocks"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respCode  int
		respError string
		mockError error
	}{
		{
			name:     "Success",
			alias:    "test_alias",
			respCode: http.StatusFound,
			url:      "https://www.google.com/",
		},
		{
			name:      "URL not found",
			alias:     "test_alias",
			respCode:  http.StatusBadRequest,
			respError: "url not found",
			mockError: storage.ErrURLNotFound,
		},
		{
			name:      "Unexpected error",
			alias:     "test_alias",
			respCode:  http.StatusInternalServerError,
			respError: "Internal server error",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlGetterMock := mocks.NewURLGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlGetterMock.On("GetURL", tc.alias).
					Return(tc.url, tc.mockError).Once()
			}

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			redirectedToURL, err := api.GetRedirect(ts.URL+"/"+tc.alias, tc.respCode)
			require.NoError(t, err)

			assert.Equal(t, tc.url, redirectedToURL)
		})
	}
}
