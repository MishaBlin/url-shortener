package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	slogdiscard "url-service/internal/lib/logger"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"url-service/internal/http-server/handlers/url/save"
	"url-service/internal/http-server/handlers/url/save/mocks"
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
			url:      "https://www.google.com/",
			respCode: http.StatusCreated,
		},
		{
			name:      "Empty URL",
			url:       "",
			respCode:  http.StatusBadRequest,
			respError: "field URL is required",
		},
		{
			name:      "Invalid URL",
			url:       "some invalid URL",
			respCode:  http.StatusBadRequest,
			respError: "field URL is not a valid url",
		},
		{
			name:      "SaveURL Error",
			url:       "https://www.google.com/",
			respCode:  http.StatusInternalServerError,
			respError: "failed to save url",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mocks.NewURLSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).
					Return(tc.mockError).
					Once()
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaverMock)

			input := fmt.Sprintf(`{"url": "%s"}`, tc.url)

			req, err := http.NewRequest(http.MethodPost, "/url", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, tc.respCode)

			body := rr.Body.String()

			var resp save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
