package utils_test

import (
	"annotate-x/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(path string, handler gin.HandlerFunc) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET(path, handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", path, nil)
	router.ServeHTTP(w, req)
	return w
}

// Generic unmarshal helper
func parseResponse[T any](t *testing.T, body []byte) utils.Response[T] {
	var resp utils.Response[T]
	err := json.Unmarshal(body, &resp)
	assert.NoError(t, err)
	return resp
}

func TestOK(t *testing.T) {
	w := setupRouter("/ok", func(c *gin.Context) {
		utils.OK(c, gin.H{"foo": "bar"})
	})
	assert.Equal(t, http.StatusOK, w.Code)

	resp := parseResponse[map[string]string](t, w.Body.Bytes())
	assert.Equal(t, "Success", resp.Message)
	assert.Equal(t, "bar", resp.Data["foo"])
}

func TestCreated(t *testing.T) {
	w := setupRouter("/created", func(c *gin.Context) {
		utils.Created(c, gin.H{"id": 123})
	})
	assert.Equal(t, http.StatusCreated, w.Code)

	resp := parseResponse[map[string]int](t, w.Body.Bytes())
	assert.Equal(t, "Created", resp.Message)
	assert.Equal(t, 123, resp.Data["id"])
}

func TestError(t *testing.T) {
	w := setupRouter("/error", func(c *gin.Context) {
		utils.Error(c, http.StatusTeapot, "I'm a teapot")
	})
	assert.Equal(t, http.StatusTeapot, w.Code)

	resp := parseResponse[map[string]any](t, w.Body.Bytes())
	assert.Equal(t, "I'm a teapot", resp.Message)
	assert.Empty(t, resp.Data)
}

func TestAbortJSON(t *testing.T) {
	w := setupRouter("/abort", func(c *gin.Context) {
		utils.AbortJSON[string](c, http.StatusUnauthorized, "not logged in")
	})
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	resp := parseResponse[string](t, w.Body.Bytes())
	assert.Equal(t, "not logged in", resp.Message)
	assert.Equal(t, "", resp.Data)
}

func TestAbortError(t *testing.T) {
	w := setupRouter("/aborterr", func(c *gin.Context) {
		utils.AbortError(c, http.StatusForbidden, "forbidden")
	})
	assert.Equal(t, http.StatusForbidden, w.Code)

	resp := parseResponse[map[string]any](t, w.Body.Bytes())
	assert.Equal(t, "forbidden", resp.Message)
	assert.Empty(t, resp.Data)
}

// ---------- Status code shortcuts ----------

func TestBadRequest(t *testing.T) {
	w := setupRouter("/bad", func(c *gin.Context) {
		utils.BadRequest(c, "bad")
	})
	assert.Equal(t, http.StatusBadRequest, w.Code)

	resp := parseResponse[map[string]any](t, w.Body.Bytes())
	assert.Equal(t, "bad", resp.Message)
}

func TestUnauthorized(t *testing.T) {
	w := setupRouter("/unauth", func(c *gin.Context) {
		utils.Unauthorized(c, "unauthorized")
	})
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	resp := parseResponse[map[string]any](t, w.Body.Bytes())
	assert.Equal(t, "unauthorized", resp.Message)
}

func TestForbidden(t *testing.T) {
	w := setupRouter("/forbidden", func(c *gin.Context) {
		utils.Forbidden(c, "forbidden")
	})
	assert.Equal(t, http.StatusForbidden, w.Code)

	resp := parseResponse[map[string]any](t, w.Body.Bytes())
	assert.Equal(t, "forbidden", resp.Message)
}

func TestNotFound(t *testing.T) {
	w := setupRouter("/notfound", func(c *gin.Context) {
		utils.NotFound(c, "not found")
	})
	assert.Equal(t, http.StatusNotFound, w.Code)

	resp := parseResponse[map[string]any](t, w.Body.Bytes())
	assert.Equal(t, "not found", resp.Message)
}

func TestInternalServerError(t *testing.T) {
	w := setupRouter("/ise", func(c *gin.Context) {
		utils.InternalServerError(c, "oops")
	})
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	resp := parseResponse[map[string]any](t, w.Body.Bytes())
	assert.Equal(t, "oops", resp.Message)
}
