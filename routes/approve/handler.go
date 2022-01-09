package approve

import (
	"net/http"
	"time"

	"github.com/syedomair/ex-paygate-lib/lib/models"
	"github.com/syedomair/ex-paygate-lib/lib/tools/logger"
	"github.com/syedomair/ex-paygate-lib/lib/tools/request"
	"github.com/syedomair/ex-paygate-lib/lib/tools/response"
)

const (
	errorCodePrefix = "01"
)

// Controller Public
type Controller struct {
	Logger logger.Logger
	Repo   Repository
	Pay    Payment
}

//var httpClient = &http.Client{}

// Ping Public
func (c *Controller) Ping(w http.ResponseWriter, r *http.Request) {
	methodName := "Ping"
	c.Logger.Debug(request.GetRequestID(r), "M:%v start", methodName)
	start := time.Now()
	responseToken := map[string]string{"response": "authController pong"}
	c.Logger.Debug(request.GetRequestID(r), "M:%v ts %+v", methodName, time.Since(start))
	response.SuccessResponseHelper(w, responseToken, http.StatusOK)
}

// ApproveAction Public
func (c *Controller) ApproveAction(w http.ResponseWriter, r *http.Request) {
	methodName := "ApproveAction"
	c.Logger.Debug(request.GetRequestID(r), "M:%v start", methodName)
	start := time.Now()

	paramConf := make(map[string]models.ParamConf)
	paramConf["merchant_key"] = models.ParamConf{Required: true, Type: request.STRING, EmptyAllowed: false}
	paramConf["cc_number"] = models.ParamConf{Required: true, Type: request.STRING, EmptyAllowed: false}
	paramConf["cc_expiry"] = models.ParamConf{Required: true, Type: request.STRING, EmptyAllowed: false}
	paramConf["currency"] = models.ParamConf{Required: true, Type: request.STRING, EmptyAllowed: false}
	paramConf["amount"] = models.ParamConf{Required: true, Type: request.STRING, EmptyAllowed: false}

	paramMap, errCode, err := request.ValidateInputParameters(r, request.GetRequestID(r), c.Logger, paramConf, nil)
	if err != nil {
		response.ErrorResponseHelper(request.GetRequestID(r), methodName, c.Logger, w, errorCodePrefix+errCode, err.Error(), http.StatusBadRequest)
		return
	}

	merchantKey := ""
	if merchantKeyValue, ok := paramMap["merchant_key"]; ok {
		merchantKey = merchantKeyValue.(string)
	}

	merchantID, err := c.Repo.GetMerchantID(merchantKey)
	if err != nil {
		response.ErrorResponseHelper(request.GetRequestID(r), methodName, c.Logger, w, errorCodePrefix+"1", err.Error(), http.StatusBadRequest)
		return
	}

	approveObj := createApproveObject(paramMap)
	approveKey, err := c.Pay.ApprovePayment(approveObj)
	if err != nil {
		response.ErrorResponseHelper(request.GetRequestID(r), methodName, c.Logger, w, errorCodePrefix+"2", err.Error(), http.StatusBadRequest)
		return
	}

	approveObj, err = c.Repo.CreateApprove(approveObj, merchantID, approveKey)
	if err != nil {
		response.ErrorResponseHelper(request.GetRequestID(r), methodName, c.Logger, w, errorCodePrefix+"2", err.Error(), http.StatusBadRequest)
		return
	}

	responseActionID := map[string]string{"approve_key": approveObj.ApproveKey, "approved_amount_balance": approveObj.AmountBalance, "currency": approveObj.Currency}
	c.Logger.Debug(request.GetRequestID(r), "M:%v ts %+v", methodName, time.Since(start))
	response.SuccessResponseHelper(w, responseActionID, http.StatusOK)
}

// createApproveObject Public
func createApproveObject(inputApprove map[string]interface{}) *models.Approve {

	ccNumber := ""
	if ccNumberValue, ok := inputApprove["cc_number"]; ok {
		ccNumber = ccNumberValue.(string)
	}
	ccExpiry := ""
	if ccExpiryValue, ok := inputApprove["cc_expiry"]; ok {
		ccExpiry = ccExpiryValue.(string)
	}
	currency := ""
	if currencyValue, ok := inputApprove["currency"]; ok {
		currency = currencyValue.(string)
	}
	amount := ""
	if amountValue, ok := inputApprove["amount"]; ok {
		amount = amountValue.(string)
	}
	newApprove := &models.Approve{}
	newApprove.CCNumber = ccNumber
	newApprove.CCExpiry = ccExpiry
	newApprove.Currency = currency
	newApprove.Amount = amount
	newApprove.Status = 1
	newApprove.AmountBalance = amount
	newApprove.CreatedAt = time.Now().Format(time.RFC3339)

	return newApprove
}
