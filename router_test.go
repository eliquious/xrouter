package xrouter

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"context"

	"github.com/stretchr/testify/suite"
)

const ResponseBody = "Hello, World!"

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type RouterTestSuite struct {
	suite.Suite
	router   Router
	server   *httptest.Server
	URL      *url.URL
	GroupURL *url.URL
}

func OptionsTest(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "HEAD,GET,PUT,POST,DELETE,OPTIONS")
}

func HeadTest(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func GetTest(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(ResponseBody))
}

func DeleteTest(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func PostTest(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil || string(body) != ResponseBody {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func PutTest(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil || string(body) != ResponseBody {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func PatchTest(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil || string(body) != ResponseBody {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// SetupTest creates the HTTP server for test.
func (suite *RouterTestSuite) SetupSuite() {
	suite.router = New()
	suite.router.OPTIONS("/", OptionsTest)
	suite.router.HEAD("/", HeadTest)
	suite.router.GET("/", GetTest)
	suite.router.DELETE("/", DeleteTest)
	suite.router.POST("/", PostTest)
	suite.router.PUT("/", PutTest)
	suite.router.PATCH("/", PatchTest)

	// api group
	api := suite.router.Group("/api/v1")
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

	suite.server = httptest.NewServer(suite.router.Handler())
	url, _ := url.Parse(suite.server.URL)
	suite.URL = url

	url, _ = url.Parse(suite.server.URL + "/api/v1/settings")
	suite.GroupURL = url
}

// TearDownSuite stops the server
func (suite *RouterTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *RouterTestSuite) TestGet() {

	res, err := http.Get(suite.server.URL)
	if err != nil {
		suite.Fail("GET request failed", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		suite.Fail("GET: failed to read response body", err)
	}
	suite.Equal(ResponseBody, string(body))
	suite.Equal(http.StatusOK, res.StatusCode)
}

func (suite *RouterTestSuite) TestNestedGet() {

	res, err := http.Get(suite.server.URL + "/api/v1/apps/1/clients/users/2/info")
	if err != nil {
		suite.Fail("GET request failed", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		suite.Fail("GET: failed to read response body", err)
	}
	suite.Equal(ResponseBody, string(body))
	suite.Equal(http.StatusOK, res.StatusCode)
}

func (suite *RouterTestSuite) TestHead() {
	res, err := http.Head(suite.server.URL)
	if err != nil {
		suite.Fail("HEAD request failed", err)
	}
	res.Body.Close()
	suite.Equal(http.StatusOK, res.StatusCode)
}

func (suite *RouterTestSuite) TestDelete() {
	client := http.Client{}
	res, err := client.Do(&http.Request{
		Method: "DELETE",
		URL:    suite.URL,
	})
	if err != nil {
		suite.Fail("DELETE request failed", err.Error())
		return
	}
	res.Body.Close()
	suite.Equal(http.StatusOK, res.StatusCode)
}

func (suite *RouterTestSuite) TestPost() {
	client := http.Client{}
	res, err := client.Post(suite.server.URL, "application/text", strings.NewReader(ResponseBody))
	if err != nil {
		suite.Fail("POST request failed", err)
	}
	res.Body.Close()
	suite.Equal(http.StatusOK, res.StatusCode)
}

func (suite *RouterTestSuite) TestPostFail() {
	client := http.Client{}
	res, err := client.Post(suite.server.URL, "application/text", nil)
	if err != nil {
		suite.Fail("POST request failed", err)
	}
	res.Body.Close()
	suite.Equal(http.StatusBadRequest, res.StatusCode)
}

func (suite *RouterTestSuite) TestPut() {
	client := http.Client{}
	res, err := client.Do(&http.Request{
		Method:        "PUT",
		URL:           suite.URL,
		Body:          ioutil.NopCloser(strings.NewReader(ResponseBody)),
		ContentLength: int64(len(ResponseBody)),
	})
	if err != nil {
		suite.Fail("PUT request failed", err.Error())
		return
	}
	res.Body.Close()
	suite.Equal(http.StatusOK, res.StatusCode)
}

func (suite *RouterTestSuite) TestPutFail() {
	client := http.Client{}
	res, err := client.Do(&http.Request{
		Method: "PUT",
		URL:    suite.URL,
	})

	if err != nil {
		suite.Fail("PUT request failed", err.Error())
		return
	}
	res.Body.Close()
	suite.Equal(http.StatusBadRequest, res.StatusCode)
}

func (suite *RouterTestSuite) TestPatch() {
	client := http.Client{}
	res, err := client.Do(&http.Request{
		Method:        "PATCH",
		URL:           suite.URL,
		Body:          ioutil.NopCloser(strings.NewReader(ResponseBody)),
		ContentLength: int64(len(ResponseBody)),
	})
	if err != nil {
		suite.Fail("PATCH request failed", err.Error())
		return
	}
	res.Body.Close()
	suite.Equal(http.StatusOK, res.StatusCode)
}

func (suite *RouterTestSuite) TestPatchFail() {
	client := http.Client{}
	res, err := client.Do(&http.Request{
		Method: "PATCH",
		URL:    suite.URL,
	})

	if err != nil {
		suite.Fail("PATCH request failed", err.Error())
		return
	}
	res.Body.Close()
	suite.Equal(http.StatusBadRequest, res.StatusCode)
}

func (suite *RouterTestSuite) TestOptions() {
	client := http.Client{}
	res, err := client.Do(&http.Request{
		Method: "OPTIONS",
		URL:    suite.URL,
	})
	if err != nil {
		suite.Fail("OPTIONS request failed", err.Error())
		return
	}
	res.Body.Close()
	suite.Equal(http.StatusOK, res.StatusCode)
	suite.Equal(1, len(res.Header["Allow"]))
	suite.Equal("HEAD,GET,PUT,POST,DELETE,OPTIONS", res.Header["Allow"][0])
}

func (suite *RouterTestSuite) TestHandler() {
	suite.NotNil(suite.router.Handler())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestRouterTestSuite(t *testing.T) {
	suite.Run(t, new(RouterTestSuite))
}
