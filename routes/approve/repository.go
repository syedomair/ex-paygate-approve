package approve

// Repository interface
type Repository interface {
	SetRequestID(requestID string)
	GetMerchantID(merchantKey string) (int, error)
	CreateApprove(inputApprove map[string]interface{}, merchantID int) (string, error)
}
