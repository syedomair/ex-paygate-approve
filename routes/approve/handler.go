package approve

import (
	"errors"
	"net/http"
	"time"

	creditcard "github.com/durango/go-credit-card"
	"github.com/syedomair/ex-paygate-lib/lib/models"
	"github.com/syedomair/ex-paygate-lib/lib/tools/logger"
	"github.com/syedomair/ex-paygate-lib/lib/tools/request"
	"github.com/syedomair/ex-paygate-lib/lib/tools/response"
)

const (
	errorCodePrefix       = "01"
	MerchantKey           = "merchant_key"
	CcNumber              = "cc_number"
	CcCVV                 = "cc_cvv"
	CcMonth               = "cc_month"
	CcYear                = "cc_year"
	Currency              = "currency"
	Amount                = "amount"
	ApproveKey            = "approve_key"
	ApprovedAmountBalance = "approved_amount_balance"
)

// Controller Public
type Controller struct {
	Logger logger.Logger
	Repo   Repository
	Pay    Payment
}

// Ping Public
func (c *Controller) Ping(w http.ResponseWriter, r *http.Request) {
	methodName := "Ping"
	c.Logger.Debug(request.GetRequestID(r), "M:%v start", methodName)
	start := time.Now()
	responseToken := map[string]string{"response": "approveController pong"}
	c.Logger.Debug(request.GetRequestID(r), "M:%v ts %+v", methodName, time.Since(start))
	response.SuccessResponseHelper(w, responseToken, http.StatusOK)
}

// ApproveAction Public
func (c *Controller) ApproveAction(w http.ResponseWriter, r *http.Request) {
	methodName := "ApproveAction"
	c.Logger.Debug(request.GetRequestID(r), "M:%v start", methodName)
	start := time.Now()

	paramConf := make(map[string]models.ParamConf)
	paramConf[MerchantKey] = models.ParamConf{Required: true, Type: request.STRING, EmptyAllowed: false}
	paramConf[CcNumber] = models.ParamConf{Required: true, Type: request.STRING, EmptyAllowed: false}
	paramConf[CcCVV] = models.ParamConf{Required: true, Type: request.STRING, EmptyAllowed: false}
	paramConf[CcMonth] = models.ParamConf{Required: true, Type: request.STRING, EmptyAllowed: false}
	paramConf[CcYear] = models.ParamConf{Required: true, Type: request.STRING, EmptyAllowed: false}
	paramConf[Currency] = models.ParamConf{Required: true, Type: request.STRING, EmptyAllowed: false}
	paramConf[Amount] = models.ParamConf{Required: true, Type: request.STRING, EmptyAllowed: false}

	paramMap, errCode, err := request.ValidateInputParameters(r, request.GetRequestID(r), c.Logger, paramConf, nil)
	if err != nil {
		response.ErrorResponseHelper(request.GetRequestID(r), methodName, c.Logger, w, errorCodePrefix+errCode, err.Error(), http.StatusBadRequest)
		return
	}

	merchantKey := ""
	if merchantKeyValue, ok := paramMap[MerchantKey]; ok {
		merchantKey = merchantKeyValue.(string)
	}

	merchantID, err := c.Repo.GetMerchantID(merchantKey)
	if err != nil {
		response.ErrorResponseHelper(request.GetRequestID(r), methodName, c.Logger, w, errorCodePrefix+"1", err.Error(), http.StatusBadRequest)
		return
	}

	approveObj := createApproveObject(paramMap)

	//Luhn Check
	card := creditcard.Card{Number: approveObj.CCNumber, Cvv: approveObj.CCCVV, Month: approveObj.CCMonth, Year: approveObj.CCYear}
	err = card.Validate(true)
	if err != nil {
		response.ErrorResponseHelper(request.GetRequestID(r), methodName, c.Logger, w, errorCodePrefix+"2", errors.New("invalid credit card").Error(), http.StatusBadRequest)
		return

	}

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

	responseAction := map[string]string{ApproveKey: approveObj.ApproveKey, ApprovedAmountBalance: approveObj.AmountBalance, Currency: approveObj.Currency}
	c.Logger.Debug(request.GetRequestID(r), "M:%v ts %+v", methodName, time.Since(start))
	response.SuccessResponseHelper(w, responseAction, http.StatusOK)
}

// createApproveObject Public
func createApproveObject(inputApprove map[string]interface{}) *models.Approve {

	ccNumber := ""
	if ccNumberValue, ok := inputApprove[CcNumber]; ok {
		ccNumber = ccNumberValue.(string)
	}
	ccCVV := ""
	if ccCVVValue, ok := inputApprove[CcCVV]; ok {
		ccCVV = ccCVVValue.(string)
	}
	ccMonth := ""
	if ccMonthValue, ok := inputApprove[CcMonth]; ok {
		ccMonth = ccMonthValue.(string)
	}
	ccYear := ""
	if ccYearValue, ok := inputApprove[CcYear]; ok {
		ccYear = ccYearValue.(string)
	}
	currency := ""
	if currencyValue, ok := inputApprove[Currency]; ok {
		currency = currencyValue.(string)
	}
	amount := ""
	if amountValue, ok := inputApprove[Amount]; ok {
		amount = amountValue.(string)
	}
	newApprove := &models.Approve{}
	newApprove.CCNumber = ccNumber
	newApprove.CCCVV = ccCVV
	newApprove.CCMonth = ccMonth
	newApprove.CCYear = ccYear
	newApprove.Currency = currency
	newApprove.Amount = amount
	newApprove.Status = 1
	newApprove.AmountBalance = amount
	newApprove.CreatedAt = time.Now().Format(time.RFC3339)
	return newApprove
}
