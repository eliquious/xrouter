package xrouter

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouterGroupPath(t *testing.T) {
	r := New()

	api := r.Group("/api/v1")
	assert.Equal(t, "/api/v1", api.Path(), "API has wrong path")

	g1 := api.Group("/group1")
	assert.Equal(t, "/api/v1/group1", g1.Path(), "API group1 has wrong path")

	g2 := api.Group("/group2")
	assert.Equal(t, "/api/v1/group2", g2.Path(), "API group2 has wrong path")

	g3 := g2.Group("/group3")
	assert.Equal(t, "/api/v1/group2/group3", g3.Path(), "API group3 has wrong path")
}

func (suite *RouterTestSuite) TestGetGroup() {

	res, err := http.Get(suite.server.URL + "/api/v1/settings")
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

func (suite *RouterTestSuite) TestGetNestedGroup() {

	res, err := http.Get(suite.server.URL + "/api/v1/group1/hello")
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

func (suite *RouterTestSuite) TestGetNestedNestedGroup() {

	res, err := http.Get(suite.server.URL + "/api/v1/group1/group2/hello")
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

func (suite *RouterTestSuite) TestHeadGroup() {
	res, err := http.Head(suite.server.URL + "/api/v1/settings")
	if err != nil {
		suite.Fail("HEAD request failed", err)
	}
	res.Body.Close()
	suite.Equal(http.StatusOK, res.StatusCode)
}

func (suite *RouterTestSuite) TestDeleteGroup() {
	client := http.Client{}
	res, err := client.Do(&http.Request{
		Method: "DELETE",
		URL:    suite.GroupURL,
	})
	if err != nil {
		suite.Fail("DELETE request failed", err.Error())
		return
	}
	res.Body.Close()
	suite.Equal(http.StatusOK, res.StatusCode)
}

func (suite *RouterTestSuite) TestPostGroup() {
	client := http.Client{}
	res, err := client.Post(suite.server.URL+"/api/v1/settings", "application/text", strings.NewReader(ResponseBody))
	if err != nil {
		suite.Fail("POST request failed", err)
	}
	res.Body.Close()
	suite.Equal(http.StatusOK, res.StatusCode)
}

func (suite *RouterTestSuite) TestPostFailGroup() {
	client := http.Client{}
	res, err := client.Post(suite.server.URL+"/api/v1/settings", "application/text", nil)
	if err != nil {
		suite.Fail("POST request failed", err)
	}
	res.Body.Close()
	suite.Equal(http.StatusBadRequest, res.StatusCode)
}

func (suite *RouterTestSuite) TestPutGroup() {
	client := http.Client{}
	res, err := client.Do(&http.Request{
		Method:        "PUT",
		URL:           suite.GroupURL,
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

func (suite *RouterTestSuite) TestPutFailGroup() {
	client := http.Client{}
	res, err := client.Do(&http.Request{
		Method: "PUT",
		URL:    suite.GroupURL,
	})

	if err != nil {
		suite.Fail("PUT request failed", err.Error())
		return
	}
	res.Body.Close()
	suite.Equal(http.StatusBadRequest, res.StatusCode)
}

func (suite *RouterTestSuite) TestPatchGroup() {
	client := http.Client{}
	res, err := client.Do(&http.Request{
		Method:        "PATCH",
		URL:           suite.GroupURL,
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

func (suite *RouterTestSuite) TestPatchFailGroup() {
	client := http.Client{}
	res, err := client.Do(&http.Request{
		Method: "PATCH",
		URL:    suite.GroupURL,
	})

	if err != nil {
		suite.Fail("PATCH request failed", err.Error())
		return
	}
	res.Body.Close()
	suite.Equal(http.StatusBadRequest, res.StatusCode)
}

func (suite *RouterTestSuite) TestOptionsGroup() {
	client := http.Client{}
	res, err := client.Do(&http.Request{
		Method: "OPTIONS",
		URL:    suite.GroupURL,
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
