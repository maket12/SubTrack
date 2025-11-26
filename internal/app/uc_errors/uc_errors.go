package uc_errors

import "errors"

var (
	ErrEmptyServiceName      = errors.New("empty service name")
	ErrInvalidDate           = errors.New("bad date format, expected DD-MM-YYYY")
	ErrEmptyDate             = errors.New("empty date")
	ErrEmptyUserID           = errors.New("empty user_id")
	ErrInvalidUserID         = errors.New("user_id is not a valid UUID")
	ErrCreateSubscription    = errors.New("failed to create subscription")
	ErrInvalidSubscriptionID = errors.New("id must be positive")
	ErrGetSubscription       = errors.New("failed to get subscription")
	ErrSubscriptionNotFound  = errors.New("subscription not found")
	ErrUpdateSubscription    = errors.New("failed to update subscription")
	ErrDeleteSubscription    = errors.New("failed to delete subscription")
	ErrInvalidLimit          = errors.New("limit must be positive or 0")
	ErrInvalidOffset         = errors.New("offset must be positive or 0")
	ErrGetSubscriptionList   = errors.New("failed to get subscription list")
	ErrGetTotalSum           = errors.New("failed to get total sum")
)
