package http

import (
	"errors"
	"net/http"

	"github.com/maket12/SubTrack/internal/app/uc_errors"
)

func HttpError(err error) (int, string) {
	switch {
	case errors.Is(err, uc_errors.ErrSubscriptionNotFound):
		return http.StatusNotFound, "subscription not found"

	case errors.Is(err, uc_errors.ErrEmptyServiceName):
		return http.StatusBadRequest, "empty service name"

	case errors.Is(err, uc_errors.ErrInvalidDate):
		return http.StatusBadRequest, "bad date format, expected DD-MM-YYYY"

	case errors.Is(err, uc_errors.ErrEmptyDate):
		return http.StatusBadRequest, "empty date"

	case errors.Is(err, uc_errors.ErrEmptyUserID):
		return http.StatusBadRequest, "empty user_id"

	case errors.Is(err, uc_errors.ErrInvalidUserID):
		return http.StatusBadRequest, "user_id is not a valid UUID"

	case errors.Is(err, uc_errors.ErrCreateSubscription):
		return http.StatusInternalServerError, "failed to create subscription"

	case errors.Is(err, uc_errors.ErrInvalidSubscriptionID):
		return http.StatusBadRequest, "id must be positive"

	case errors.Is(err, uc_errors.ErrGetSubscription):
		return http.StatusInternalServerError, "failed to get subscription"

	case errors.Is(err, uc_errors.ErrUpdateSubscription):
		return http.StatusInternalServerError, "failed to update subscription"

	case errors.Is(err, uc_errors.ErrDeleteSubscription):
		return http.StatusInternalServerError, "failed to delete subscription"

	case errors.Is(err, uc_errors.ErrInvalidLimit):
		return http.StatusBadRequest, "limit must be positive or 0"

	case errors.Is(err, uc_errors.ErrInvalidOffset):
		return http.StatusBadRequest, "offset must be positive or 0"

	case errors.Is(err, uc_errors.ErrGetSubscriptionList):
		return http.StatusInternalServerError, "failed to get subscription list"

	case errors.Is(err, uc_errors.ErrGetTotalSum):
		return http.StatusInternalServerError, "failed to get total sum"

	default:
		return http.StatusInternalServerError, "internal error"
	}
}
