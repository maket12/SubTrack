package port

import "errors"

var (
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrCreateSubscription   = errors.New("failed to create subscription")
	ErrUpdateSubscription   = errors.New("failed to update subscription")
	ErrDeleteSubscription   = errors.New("failed to delete subscription")
	ErrGetSubscription      = errors.New("failed to get subscription")
	ErrGetSubscriptionList  = errors.New("failed to get subscription list")
	ErrGetTotalSum          = errors.New("failed to get total sum")
)
