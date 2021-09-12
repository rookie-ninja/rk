package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	httptest "github.com/stretchr/testify/http"
	"net/http"
	"net/url"
	"testing"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func TestGreeter_HappyCase(t *testing.T) {
	// Create a mock writer
	w := &httptest.TestResponseWriter{}

	// Create a gin context and inject query(name=ut) into it
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		URL: &url.URL{
			RawQuery: "name=ut",
		},
	}

	// Call gin handler
	Greeter(ctx)
	// Assert expected result
	assert.Contains(t, w.Output, "Hello ut!")
}
