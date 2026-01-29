package iaps

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestIAPOneStore(t *testing.T) {
	suite.Run(t, new(SuiteIAPOneStore))
}

type SuiteIAPOneStore struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteIAPOneStore) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-iaps-iapOneStore"))
}

func (this *SuiteIAPOneStore) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteIAPOneStore) TestIAPOneStore() {
	target := NewIAPOneStore(&IAPOneStoreConfig{})
	this.NotNil(target)
	this.Nil(target.Initialize())
	this.NotNil(target.Client())
	target.Finalize()
	this.Nil(target.Initialize(&testIAPOneStoreClient{}))
	this.NotNil(target.Client())
	target.Finalize()
}

func (this *SuiteIAPOneStore) TestVerify() {
	target := NewIAPOneStore(&IAPOneStoreConfig{})
	_ = target.Initialize(&testIAPOneStoreClient{
		token:  true,
		verify: true,
	})
	result := target.Verify(testdata.Unknown, testdata.Unknown)
	this.Nil(result.Err)
	result = target.Verify(testdata.Unknown, testdata.Unknown)
	this.Nil(result.Err)
	target.Finalize()
	this.NotNil(target.Verify(testdata.Unknown, testdata.Unknown))

	_ = target.Initialize(&testIAPOneStoreClient{
		token:  false,
		verify: true,
	})
	result = target.Verify(testdata.Unknown, testdata.Unknown)
	this.NotNil(result.Err)
	target.Finalize()

	_ = target.Initialize(&testIAPOneStoreClient{
		token:  true,
		verify: false,
	})
	result = target.Verify(testdata.Unknown, testdata.Unknown)
	this.NotNil(result.Err)
	target.Finalize()
}

type testIAPOneStoreClient struct {
	token  bool
	verify bool
}

func (this *testIAPOneStoreClient) Do(req *http.Request) (*http.Response, error) {
	body := &bytes.Buffer{}

	if strings.Contains(req.URL.Path, "/v7/oauth/token") {
		if this.token {
			_, _ = fmt.Fprintf(body, `{"access_token":"test-token","expires_in":3600}`)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(body),
				Header:     make(http.Header),
			}, nil
		} // if

		_, _ = fmt.Fprintf(body, `{"error":{"code":"401","message":"invalid client"}}`)
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(body),
			Header:     make(http.Header),
		}, nil
	} // if

	if strings.Contains(req.URL.Path, "/purchases/inapp") {
		if this.verify {
			_, _ = fmt.Fprintf(body, `{"purchaseTime":1,"purchaseState":0}`)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(body),
				Header:     make(http.Header),
			}, nil
		} // if

		_, _ = fmt.Fprintf(body, `{"error":{"code":"404","message":"not found"}}`)
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(body),
			Header:     make(http.Header),
		}, nil
	} // if

	return nil, fmt.Errorf("unexpected request: %v", req.URL.Path)
}
