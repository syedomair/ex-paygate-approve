package approve

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/syedomair/ex-paygate-lib/lib/models"
	"github.com/syedomair/ex-paygate-lib/lib/tools/logger"
)

func TestApproveAction(t *testing.T) {
	c := Controller{
		Logger: logger.New("DEBUG", "TEST#", os.Stdout),
		Repo:   &mockDB{},
		Pay:    &mockPay{}}

	method := "POST"
	url := "/approve"

	res, req := MockTestServer(method, url, []byte(`{"merchant_key":"KEY1", "cc_number":"123", "cc_expiry":"12345", "currency":"USD", "amount":"88"}`))
	c.ApproveAction(res, req)
	type TestResponse struct {
		Data   string
		Result string
	}
	response := new(TestResponse)
	json.NewDecoder(res.Result().Body).Decode(response)

	expected := "success"
	if expected != response.Result {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, response.Result)
	}

}

type mockPay struct {
}

func (mdb *mockPay) ApprovePayment(*models.Approve) (string, error) {
	return "approval_key", nil
}

type mockDB struct {
}

func (mdb *mockDB) SetRequestID(requestID string) {
}

func (mdb *mockDB) GetMerchantID(merchantKey string) (int, error) {
	return 1, nil

}
func (mdb *mockDB) CreateApprove(approveModel *models.Approve, merchantID int, approveKey string) (*models.Approve, error) {
	return &models.Approve{}, nil
}
