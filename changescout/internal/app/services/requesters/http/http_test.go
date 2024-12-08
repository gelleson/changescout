package http_test

//import (
//	"errors"
//	http2 "github.com/gelleson/changescout/changescout/internal/app/services/requesters/http"
//	"github.com/stretchr/testify/mock"
//	"io/ioutil"
//	"net/http"
//	"strings"
//	"testing"
//
//	"github.com/gelleson/changescout/changescout/internal/app/services/mocks"
//	"github.com/gelleson/changescout/changescout/internal/domain"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/suite"
//)
//
//type HttpServiceSuite struct {
//	suite.Suite
//	httpService *http2.HttpService
//	mockDoer    *mocks.Doer
//}
//
//func (suite *HttpServiceSuite) SetupTest() {
//	suite.mockDoer = new(mocks.Doer)
//	suite.httpService = http2.New(suite.mockDoer)
//}
//
//func (suite *HttpServiceSuite) TestRequest_ValidResponse() {
//	expectedBody := "response body"
//	resp := &http.Response{
//		StatusCode: http.StatusOK,
//		Body:       ioutil.NopCloser(strings.NewReader(expectedBody)),
//	}
//
//	site := domain.Website{
//		URL: "http://example.com",
//		Setting: domain.Setting{
//			Method:  http.MethodGet,
//			Headers: http.Header{},
//		},
//	}
//
//	suite.mockDoer.On("Do", mock.Anything).Return(resp, nil)
//
//	body, err := suite.httpService.Request(site)
//	assert.NoError(suite.T(), err)
//	assert.Equal(suite.T(), []byte(expectedBody), body)
//
//	suite.mockDoer.AssertExpectations(suite.T())
//}
//
//func (suite *HttpServiceSuite) TestRequest_BadStatusCode() {
//	resp := &http.Response{
//		StatusCode: http.StatusBadRequest,
//		Body:       ioutil.NopCloser(strings.NewReader("")),
//	}
//
//	site := domain.Website{
//		URL: "http://example.com",
//		Setting: domain.Setting{
//			Method:  http.MethodGet,
//			Headers: http.Header{},
//		},
//	}
//
//	suite.mockDoer.On("Do", mock.Anything).Return(resp, nil)
//
//	body, err := suite.httpService.Request(site)
//	assert.Error(suite.T(), err)
//	assert.Nil(suite.T(), body)
//	assert.Equal(suite.T(), "bad status code", err.Error())
//
//	suite.mockDoer.AssertExpectations(suite.T())
//}
//
//func (suite *HttpServiceSuite) TestRequest_DoerReturnsError() {
//	site := domain.Website{
//		URL: "http://example.com",
//		Setting: domain.Setting{
//			Method:  http.MethodGet,
//			Headers: http.Header{},
//		},
//	}
//
//	suite.mockDoer.On("Do", mock.Anything).Return(nil, errors.New("doer error"))
//
//	body, err := suite.httpService.Request(site)
//	assert.Error(suite.T(), err)
//	assert.Nil(suite.T(), body)
//	assert.Equal(suite.T(), "doer error", err.Error())
//
//	suite.mockDoer.AssertExpectations(suite.T())
//}
//
//func (suite *HttpServiceSuite) TestRequest_ReadBodyError() {
//	resp := &http.Response{
//		StatusCode: http.StatusOK,
//		Body:       &errorReadCloser{},
//	}
//
//	site := domain.Website{
//		URL: "http://example.com",
//		Setting: domain.Setting{
//			Method:  http.MethodGet,
//			Headers: http.Header{},
//		},
//	}
//
//	suite.mockDoer.On("Do", mock.Anything).Return(resp, nil)
//
//	body, err := suite.httpService.Request(site)
//	assert.Error(suite.T(), err)
//	assert.Nil(suite.T(), body)
//
//	suite.mockDoer.AssertExpectations(suite.T())
//}
//
//// Custom ReadCloser to simulate a read error
//type errorReadCloser struct{}
//
//func (e *errorReadCloser) Read(p []byte) (int, error) {
//	return 0, errors.New("read error")
//}
//
//func (e *errorReadCloser) Close() error {
//	return nil
//}
//
//func TestHttpServiceSuite(t *testing.T) {
//	suite.Run(t, new(HttpServiceSuite))
//}
