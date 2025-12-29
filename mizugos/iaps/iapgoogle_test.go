package iaps

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/api/androidpublisher/v3"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestIAPGoogle(t *testing.T) {
	suite.Run(t, new(SuiteIAPGoogle))
}

type SuiteIAPGoogle struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteIAPGoogle) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-iaps-iapGoogle"))
}

func (this *SuiteIAPGoogle) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteIAPGoogle) TestIAPGoogle() {
	target := NewIAPGoogle(&IAPGoogleConfig{})
	this.NotNil(target)
	this.Nil(target.Initialize(&testIAPGoogleClient{}))
	target.Finalize()
}

func (this *SuiteIAPGoogle) TestVerify() {
	target := NewIAPGoogle(&IAPGoogleConfig{})
	_ = target.Initialize(&testIAPGoogleClient{
		verify: true,
	})
	result := target.Verify(testdata.Unknown, testdata.Unknown)
	this.Nil(result.Err)
	target.Finalize()
	this.NotNil(target.Verify(testdata.Unknown, testdata.Unknown))

	_ = target.Initialize(&testIAPGoogleClient{
		verify: false,
	})
	result = target.Verify(testdata.Unknown, testdata.Unknown)
	this.NotNil(result.Err)
	target.Finalize()
}

type testIAPGoogleClient struct {
	verify bool
}

func (this *testIAPGoogleClient) VerifyProduct(context.Context, string, string, string) (*androidpublisher.ProductPurchase, error) {
	if this.verify {
		return &androidpublisher.ProductPurchase{}, nil
	} // if

	return nil, fmt.Errorf("fail")
}
