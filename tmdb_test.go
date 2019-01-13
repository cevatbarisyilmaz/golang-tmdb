package tmdb

import (
	"bytes"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TMBDTestSuite struct {
	suite.Suite
	Client
}

func (suite *TMBDTestSuite) SetupTest() {
	suite.APIKey = os.Getenv("APIKey")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TMBDTestSuite))
}

func (suite *TMBDTestSuite) TestGetFail() {
	err := suite.get("http://www.testfakewebsite.org", nil)
	suite.Contains(err.Error(), "no such host")
	err = suite.get("https://api.themoviedb.org/3/movie/7578000?language=en-US", nil)
	suite.Equal("Invalid API key: You must be granted a valid key.", err.Error())
	err = suite.get("", nil)
	suite.Equal("url field is empty", err.Error())
}

func (suite *TMBDTestSuite) TestDecodeDataFail() {
	b := []byte(`{}`)
	err := suite.get("https://www.google.com.br", b)
	suite.Contains(err.Error(), "could not decode the data")
}

func (suite *TMBDTestSuite) TestDecodeErrorFail() {
	r, err := http.Get("https://golang.org/")
	err = suite.decodeError(r)
	defer r.Body.Close()
	suite.Contains(err.Error(), "couldn't decode error")
}

func (suite *TMBDTestSuite) TestDecodeErrorEmptyBodyFail() {
	r, err := http.Get("https://golang.org/")
	r.Write(bytes.NewBuffer([]byte("")))
	err = suite.decodeError(r)
	defer r.Body.Close()
	suite.Contains(err.Error(), "empty body")
}

func (suite *TMBDTestSuite) TestDecodeErrorReadBodyFail() {
	r, err := http.Get("https://golang.org/")
	r.Body.Close()
	err = suite.decodeError(r)
	suite.Contains(err.Error(), "could not read body response")
}
