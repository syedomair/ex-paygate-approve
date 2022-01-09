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
)

func TestApproveAction(t *testing.T) {
	c := Controller{
		Logger: logger.New("DEBUG", "TEST#", os.Stdout),
		Repo:   &mockDB{},
		Pay:    &mockPay{}}

	method := "POST"
	url := "/approve"

	type TestResponse struct {
		Data   string
		Result string
	}

	//Invalid CC number
	res, req := mockserver.MockTestServer(method, url, []byte(`{"merchant_key":"KEY1", "cc_number":"123", "cc_expiry":"12345", "currency":"USD", "amount":"88"}`))
	c.ApproveAction(res, req)
	response := new(TestResponse)
	json.NewDecoder(res.Result().Body).Decode(response)

	expected := "failure"
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}

	//Valid CC number
	res, req = mockserver.MockTestServer(method, url, []byte(`{"merchant_key":"KEY1", "cc_number":"4000000000000000", "cc_expiry":"12345", "currency":"USD", "amount":"88"}`))
	c.ApproveAction(res, req)
	response = new(TestResponse)
	json.NewDecoder(res.Result().Body).Decode(response)

	expected = "success"
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}

}

type mockPay struct {
}

func (mdb *mockPay) ApprovePayment(approveObj *models.Approve) (string, error) {
	if approveObj.CCNumber != "4000000000000000" {
		return "", errors.New("invalid credit card")
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
	if merchantKey == "KEY1" {
		return 1, nil
	}
	return 0, errors.New("merchant not found")
}
func (mdb *mockDB) CreateApprove(approveModel *models.Approve, merchantID int, approveKey string) (*models.Approve, error) {
	return &models.Approve{}, nil
}
