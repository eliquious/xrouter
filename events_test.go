package xrouter

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouterEventHandler(t *testing.T) {
	routes := make(map[string]bool)
	routes["GET /api/v1/settings"] = false
	routes["GET /api/v1/group1/hello"] = false
	routes["GET /api/v1/group1/group2/hello"] = false
	routes["GET /api/v1/apps/:app/clients/users/:userid/info"] = false
	routes["OPTIONS /api/v1/settings"] = false
	routes["HEAD /api/v1/settings"] = false
	routes["DELETE /api/v1/settings"] = false
	routes["POST /api/v1/settings"] = false
	routes["PUT /api/v1/settings"] = false
	routes["PATCH /api/v1/settings"] = false
	routes["NotFound"] = false
	routes["MethodNotAllowed"] = false

	r := New()
	r.EventHandler(func(evt Event) {
		switch e := evt.(type) {
		case AddHandlerEvent:
			routes[e.Method+" "+e.Path] = true
			// t.Logf("%s %s", e.Method, e.Path)
		case NotFoundHandlerEvent:
			routes["NotFound"] = true
		case MethodNotAllowedHandlerEvent:
			routes["MethodNotAllowed"] = true
		}
	})
	r.NotFound(http.NotFoundHandler())
	r.MethodNotAllowed(http.NotFoundHandler())

	// api group
	api := r.Group("/api/v1")
	api.GET("/settings", GetTest)
	api.OPTIONS("/settings", OptionsTest)
	api.HEAD("/settings", HeadTest)
	api.DELETE("/settings", DeleteTest)
	api.POST("/settings", PostTest)
	api.PUT("/settings", PutTest)
	api.PATCH("/settings", PatchTest)

	// Create router group and handler
	group1 := api.Group("/group1")
	group1.GET("/hello", GetTest)

	// Create nested group and route
	group2 := group1.Group("/group2")
	group2.GET("/hello", GetTest)

	appGroup := api.Group("/apps/:app")
	clientGroup := appGroup.Group("/clients")
	clientGroup.GET("/users/:userid/info", GetTest)

	for k, v := range routes {
		assert.True(t, v, "%s should be true", k)
	}
}
