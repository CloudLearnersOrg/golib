package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTest() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up a mock request
	req := httptest.NewRequest("GET", "/", nil)
	c.Request = req

	return c, w
}

func TestStatusResponses(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testErr := errors.New("test error")

	tests := []struct {
		name       string
		execute    func(*gin.Context)
		wantStatus int
		wantBody   Response
	}{
		{
			name: "StatusOK",
			execute: func(c *gin.Context) {
				StatusOK(c, "success", map[string]string{"key": "value"})
			},
			wantStatus: http.StatusOK,
			wantBody: Response{
				Code:    http.StatusOK,
				Message: "success",
				Data:    map[string]interface{}{"key": "value"},
			},
		},
		{
			name: "StatusCreated",
			execute: func(c *gin.Context) {
				StatusCreated(c, "created", map[string]string{"id": "123"})
			},
			wantStatus: http.StatusCreated,
			wantBody: Response{
				Code:    http.StatusCreated,
				Message: "created",
				Data:    map[string]interface{}{"id": "123"},
			},
		},
		// 4xx status codes
		{
			name: "StatusBadRequest",
			execute: func(c *gin.Context) {
				StatusBadRequest(c, "bad request", testErr)
			},
			wantStatus: http.StatusBadRequest,
			wantBody: Response{
				Code:    http.StatusBadRequest,
				Message: "bad request",
				Error:   testErr.Error(),
			},
		},
		{
			name: "StatusUnauthorized",
			execute: func(c *gin.Context) {
				StatusUnauthorized(c, "unauthorized", testErr)
			},
			wantStatus: http.StatusUnauthorized,
			wantBody: Response{
				Code:    http.StatusUnauthorized,
				Message: "unauthorized",
				Error:   testErr.Error(),
			},
		},
		{
			name: "StatusForbidden",
			execute: func(c *gin.Context) {
				StatusForbidden(c, "forbidden", testErr)
			},
			wantStatus: http.StatusForbidden,
			wantBody: Response{
				Code:    http.StatusForbidden,
				Message: "forbidden",
				Error:   testErr.Error(),
			},
		},
		{
			name: "StatusNotFound",
			execute: func(c *gin.Context) {
				StatusNotFound(c, "not found", testErr)
			},
			wantStatus: http.StatusNotFound,
			wantBody: Response{
				Code:    http.StatusNotFound,
				Message: "not found",
				Error:   testErr.Error(),
			},
		},
		{
			name: "StatusRequestTimeout",
			execute: func(c *gin.Context) {
				StatusRequestTimeout(c, "timeout", testErr)
			},
			wantStatus: http.StatusRequestTimeout,
			wantBody: Response{
				Code:    http.StatusRequestTimeout,
				Message: "timeout",
				Error:   testErr.Error(),
			},
		},
		{
			name: "StatusConflict",
			execute: func(c *gin.Context) {
				StatusConflict(c, "conflict", testErr)
			},
			wantStatus: http.StatusConflict,
			wantBody: Response{
				Code:    http.StatusConflict,
				Message: "conflict",
				Error:   testErr.Error(),
			},
		},
		{
			name: "StatusUnprocessableEntity",
			execute: func(c *gin.Context) {
				StatusUnprocessableEntity(c, "unprocessable", testErr)
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantBody: Response{
				Code:    http.StatusUnprocessableEntity,
				Message: "unprocessable",
				Error:   testErr.Error(),
			},
		},
		{
			name: "StatusTooManyRequests",
			execute: func(c *gin.Context) {
				StatusTooManyRequests(c, "too many requests", testErr)
			},
			wantStatus: http.StatusTooManyRequests,
			wantBody: Response{
				Code:    http.StatusTooManyRequests,
				Message: "too many requests",
				Error:   testErr.Error(),
			},
		},
		// 5xx status codes
		{
			name: "StatusInternalServerError",
			execute: func(c *gin.Context) {
				StatusInternalServerError(c, "internal error", testErr)
			},
			wantStatus: http.StatusInternalServerError,
			wantBody: Response{
				Code:    http.StatusInternalServerError,
				Message: "internal error",
				Error:   testErr.Error(),
			},
		},
		{
			name: "StatusServiceUnavailable",
			execute: func(c *gin.Context) {
				StatusServiceUnavailable(c, "service down", testErr)
			},
			wantStatus: http.StatusServiceUnavailable,
			wantBody: Response{
				Code:    http.StatusServiceUnavailable,
				Message: "service down",
				Error:   testErr.Error(),
			},
		},
		{
			name: "StatusBadGateway",
			execute: func(c *gin.Context) {
				StatusBadGateway(c, "bad gateway", testErr)
			},
			wantStatus: http.StatusBadGateway,
			wantBody: Response{
				Code:    http.StatusBadGateway,
				Message: "bad gateway",
				Error:   testErr.Error(),
			},
		},
		{
			name: "StatusGatewayTimeout",
			execute: func(c *gin.Context) {
				StatusGatewayTimeout(c, "gateway timeout", testErr)
			},
			wantStatus: http.StatusGatewayTimeout,
			wantBody: Response{
				Code:    http.StatusGatewayTimeout,
				Message: "gateway timeout",
				Error:   testErr.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			c, w := setupTest()

			// When
			tt.execute(c)

			// Then
			assert.Equal(t, tt.wantStatus, w.Code)

			var response Response
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Equal(t, tt.wantBody.Code, response.Code)
			assert.Equal(t, tt.wantBody.Message, response.Message)

			if tt.wantBody.Data != nil {
				assert.Equal(t, tt.wantBody.Data, response.Data)
			}
			if tt.wantBody.Error != "" {
				assert.Equal(t, tt.wantBody.Error, response.Error)
			}
		})
	}
}

func TestRedirectResponses(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		execute      func(*gin.Context)
		wantStatus   int
		wantLocation string
	}{
		{
			name: "StatusMovedPermanently",
			execute: func(c *gin.Context) {
				StatusMovedPermanently(c, "/new-location")
			},
			wantStatus:   http.StatusMovedPermanently,
			wantLocation: "/new-location",
		},
		{
			name: "StatusFound",
			execute: func(c *gin.Context) {
				StatusFound(c, "/temp-location")
			},
			wantStatus:   http.StatusFound,
			wantLocation: "/temp-location",
		},
		{
			name: "StatusTemporaryRedirect",
			execute: func(c *gin.Context) {
				StatusTemporaryRedirect(c, "/temp")
			},
			wantStatus:   http.StatusTemporaryRedirect,
			wantLocation: "/temp",
		},
		{
			name: "StatusPermanentRedirect",
			execute: func(c *gin.Context) {
				StatusPermanentRedirect(c, "/perm")
			},
			wantStatus:   http.StatusPermanentRedirect,
			wantLocation: "/perm",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			c, w := setupTest()

			// When
			tt.execute(c)

			// Then
			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantLocation, w.Header().Get("Location"))
		})
	}
}
