package approve

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/syedomair/ex-paygate-lib/lib/models"
	"github.com/syedomair/ex-paygate-lib/lib/tools/logger"
	"github.com/syedomair/ex-paygate-lib/lib/tools/mockserver"
	"github.com/syedomair/ex-paygate-lib/lib/tools/request"
)

const (
	ValidMerchanKey    = "KEY1"
	InValidMerchanKey  = "KEYxyz"
	ValidCCNumber      = "4242424242424242"
	InValidCCNumber    = "0002424242424242"
	AuthFailedCCNumber = "4000000000000119"
)

func TestApproveAction(t *testing.T) {
	c := Controller{
		Logger: logger.New("DEBUG", "TEST#", os.Stdout),
		Repo:   &mockDB{},
		Pay:    &mockPay{}}

	method := "POST"
	url := "/authorize"

	type TestResponse struct {
		Data   string
		Result string
	}

	//Invalid Merchant KEY
	res, req := mockserver.MockTestServer(method, url, []byte(`{"merchant_key":"`+InValidMerchanKey+
		`", "cc_number":"123", "cc_cvv":"1234", "cc_month":"12", "cc_year":"2025", "currency":"USD", "amount":"88"}`))
	c.ApproveAction(res, req)
	response := new(TestResponse)
	json.NewDecoder(res.Result().Body).Decode(response)
	expected := request.FAILURE
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}

	//InValid CC number
	res, req = mockserver.MockTestServer(method, url, []byte(`{"merchant_key":"`+ValidMerchanKey+
		`", "cc_number":"`+InValidCCNumber+`", "cc_cvv":"1234", "cc_month":"12", "cc_year":"2025" ,"currency":"USD", "amount":"88"}`))
	c.ApproveAction(res, req)
	response = new(TestResponse)
	json.NewDecoder(res.Result().Body).Decode(response)
	expected = request.FAILURE
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}

	//Valid CC number
	res, req = mockserver.MockTestServer(method, url, []byte(`{"merchant_key":"`+ValidMerchanKey+
		`", "cc_number":"`+ValidCCNumber+`", "cc_cvv":"1234", "cc_month":"12", "cc_year":"2025" ,"currency":"USD", "amount":"88"}`))
	c.ApproveAction(res, req)
	response = new(TestResponse)
	json.NewDecoder(res.Result().Body).Decode(response)
	expected = request.SUCCESS
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}

	//Authorization Failed
	res, req = mockserver.MockTestServer(method, url, []byte(`{"merchant_key":"`+ValidMerchanKey+
		`", "cc_number":"`+AuthFailedCCNumber+`", "cc_cvv":"1234", "cc_month":"12", "cc_year":"2025" ,"currency":"USD", "amount":"88"}`))
	c.ApproveAction(res, req)
	response = new(TestResponse)
	json.NewDecoder(res.Result().Body).Decode(response)
	expected = request.FAILURE
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}

}

type mockPay struct {
}

func (mdb *mockPay) ApprovePayment(approveObj *models.Approve) (string, error) {
	if approveObj.CCNumber == AuthFailedCCNumber {
		return "", errors.New("authorisation failure")
	}
	key := make([]byte, 10)
	_, _ = rand.Read(key)
	return fmt.Sprintf("%X", key), nil
}

type mockDB struct {
}

func (mdb *mockDB) SetRequestID(requestID string) {
}

func (mdb *mockDB) GetMerchantID(merchantKey string) (int, error) {
	if merchantKey == ValidMerchanKey {
		return 1, nil
	}
	return 0, errors.New("merchant not found")
}
func (mdb *mockDB) CreateApprove(approveModel *models.Approve, merchantID int, approveKey string) (*models.Approve, error) {
	return &models.Approve{}, nil
}
