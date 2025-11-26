package http

import (
	"errors"
	"net/http"

	"github.com/maket12/SubTrack/internal/app/uc_errors"
)

func HttpError(err error) (int, string, error) {
	if w, ok := err.(*uc_errors.WrappedError); ok {
		switch w.Public {
		case uc_errors.ErrSubscriptionNotFound:
			return http.StatusNotFound, w.Public.Error(), w.Reason
		case uc_errors.ErrCreateSubscription,
			uc_errors.ErrGetSubscription,
			uc_errors.ErrUpdateSubscription,
			uc_errors.ErrDeleteSubscription,
			uc_errors.ErrGetSubscriptionList,
			uc_errors.ErrGetTotalSum:
			return http.StatusInternalServerError, w.Public.Error(), w.Reason
		default:
			return http.StatusInternalServerError, "internal error", w.Reason
		}
	}

	switch {
	case errors.Is(err, uc_errors.ErrSubscriptionNotFound):
		return http.StatusNotFound, err.Error(), nil
	case errors.Is(err, uc_errors.ErrEmptyServiceName),
		errors.Is(err, uc_errors.ErrInvalidDate),
		errors.Is(err, uc_errors.ErrEmptyDate),
		errors.Is(err, uc_errors.ErrEmptyUserID),
		errors.Is(err, uc_errors.ErrInvalidUserID),
		errors.Is(err, uc_errors.ErrInvalidSubscriptionID),
		errors.Is(err, uc_errors.ErrInvalidLimit),
		errors.Is(err, uc_errors.ErrInvalidOffset):
		return http.StatusBadRequest, err.Error(), nil
	}

	return http.StatusInternalServerError, "internal error", err
}
