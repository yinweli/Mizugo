package iaps

import (
	"context"
	"fmt"
	"testing"

	"github.com/awa/go-iap/appstore/api"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/helps"
	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestIAPApple(t *testing.T) {
	suite.Run(t, new(SuiteIAPApple))
}

type SuiteIAPApple struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteIAPApple) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-iaps-iapApple"))
}

func (this *SuiteIAPApple) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteIAPApple) TestIAPApple() {
	target := NewIAPApple(&IAPAppleConfig{})
	this.NotNil(target)
	this.Nil(target.Initialize())
	target.Finalize()
	this.Nil(target.Initialize(&testIAPAppleClient{}))
	target.Finalize()
}

func (this *SuiteIAPApple) TestVerify() {
	target := NewIAPApple(&IAPAppleConfig{})
	_ = target.Initialize(&testIAPAppleClient{
		info:          true,
		parse:         true,
		productID:     testdata.Unknown,
		transactionID: testdata.Unknown,
	})
	result := target.Verify(testdata.Unknown, testdata.Unknown)
	this.Nil(result.Err)
	target.Finalize()
	this.NotNil(target.Verify(testdata.Unknown, testdata.Unknown))

	_ = target.Initialize(&testIAPAppleClient{
		info:          false,
		parse:         true,
		productID:     testdata.Unknown,
		transactionID: testdata.Unknown,
	})
	result = target.Verify(testdata.Unknown, testdata.Unknown)
	this.NotNil(result.Err)
	target.Finalize()

	_ = target.Initialize(&testIAPAppleClient{
		info:          true,
		parse:         false,
		productID:     testdata.Unknown,
		transactionID: testdata.Unknown,
	})
	result = target.Verify(testdata.Unknown, testdata.Unknown)
	this.NotNil(result.Err)
	target.Finalize()

	_ = target.Initialize(&testIAPAppleClient{
		info:          true,
		parse:         true,
		productID:     helps.RandStringDefault(),
		transactionID: testdata.Unknown,
	})
	result = target.Verify(testdata.Unknown, testdata.Unknown)
	this.NotNil(result.Err)
	target.Finalize()

	_ = target.Initialize(&testIAPAppleClient{
		info:          true,
		parse:         true,
		productID:     testdata.Unknown,
		transactionID: helps.RandStringDefault(),
	})
	result = target.Verify(testdata.Unknown, testdata.Unknown)
	this.NotNil(result.Err)
	target.Finalize()
}

type testIAPAppleClient struct {
	info          bool
	parse         bool
	productID     string
	transactionID string
}

func (this *testIAPAppleClient) GetTransactionInfo(context.Context, string) (*api.TransactionInfoResponse, error) {
	if this.info {
		return &api.TransactionInfoResponse{}, nil
	} else {
		return nil, fmt.Errorf("fail")
	} // if
}

func (this *testIAPAppleClient) ParseSignedTransaction(string) (*api.JWSTransaction, error) {
	if this.parse {
		return &api.JWSTransaction{
			ProductID:     this.productID,
			TransactionID: this.transactionID,
		}, nil
	} else {
		return nil, fmt.Errorf("fail")
	} // if
}
