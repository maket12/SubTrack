package http

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/maket12/SubTrack/internal/app/dto"
	"github.com/maket12/SubTrack/internal/app/usecase"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	log        *slog.Logger
	CreateUC   *usecase.CreateSubscriptionUC
	GetUC      *usecase.GetSubscriptionUC
	UpdateUC   *usecase.UpdateSubscriptionUC
	DeleteUC   *usecase.DeleteSubscriptionUC
	ListUC     *usecase.GetSubscriptionListUC
	TotalSumUC *usecase.GetTotalSumUC
}

func NewSubscriptionHandler(
	log *slog.Logger,
	createUC *usecase.CreateSubscriptionUC,
	getUC *usecase.GetSubscriptionUC,
	updateUC *usecase.UpdateSubscriptionUC,
	deleteUC *usecase.DeleteSubscriptionUC,
	listUC *usecase.GetSubscriptionListUC,
	totalSumUC *usecase.GetTotalSumUC,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		log:        log,
		CreateUC:   createUC,
		GetUC:      getUC,
		UpdateUC:   updateUC,
		DeleteUC:   deleteUC,
		ListUC:     listUC,
		TotalSumUC: totalSumUC,
	}
}

func (h *SubscriptionHandler) Create(ctx *gin.Context) {
	var req dto.CreateSubscription
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	resp, err := h.CreateUC.Execute(ctx, req)
	if err != nil {
		status, msg, internalErr := HttpError(err)
		h.log.ErrorContext(ctx, "failed to create subscription",
			slog.Int("status", status),
			slog.String("public_msg", msg),
			slog.Any("cause", internalErr),
		)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	h.log.InfoContext(ctx, "created subscription",
		slog.Int("id", resp.ID),
	)

	ctx.JSON(http.StatusCreated, resp)
}

func (h *SubscriptionHandler) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id must be positive integer"})
		return
	}

	resp, err := h.GetUC.Execute(ctx, dto.GetSubscription{ID: id})
	if err != nil {
		status, msg, internalErr := HttpError(err)
		h.log.ErrorContext(ctx, "failed to get subscription",
			slog.Int("status", status),
			slog.String("public_msg", msg),
			slog.Any("cause", internalErr),
		)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *SubscriptionHandler) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id must be positive integer"})
		return
	}

	var req dto.UpdateSubscription
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}
	req.ID = id

	resp, err := h.UpdateUC.Execute(ctx, req)
	if err != nil {
		status, msg, internalErr := HttpError(err)
		h.log.ErrorContext(ctx, "failed to update subscription",
			slog.Int("status", status),
			slog.String("public_msg", msg),
			slog.Any("cause", internalErr),
		)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	h.log.InfoContext(ctx, "updated subscription",
		slog.Int("id", req.ID),
	)

	ctx.JSON(http.StatusOK, resp)
}

func (h *SubscriptionHandler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id must be positive integer"})
		return
	}

	resp, err := h.DeleteUC.Execute(ctx, dto.DeleteSubscription{ID: id})
	if err != nil {
		status, msg, internalErr := HttpError(err)
		h.log.ErrorContext(ctx, "failed to delete subscription",
			slog.Int("status", status),
			slog.String("public_msg", msg),
			slog.Any("cause", internalErr),
		)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	h.log.InfoContext(ctx, "deleted subscription",
		slog.Int("id", id),
	)

	ctx.JSON(http.StatusOK, resp)
}

func (h *SubscriptionHandler) List(ctx *gin.Context) {
	var req dto.GetSubscriptionList
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid query"})
		return
	}

	resp, err := h.ListUC.Execute(ctx, req)
	if err != nil {
		status, msg, internalErr := HttpError(err)
		h.log.Error("failed to get subscriptions list",
			slog.Int("status", status),
			slog.String("public_msg", msg),
			slog.Any("cause", internalErr),
		)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *SubscriptionHandler) GetTotalSum(ctx *gin.Context) {
	var req dto.GetTotalSum
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid query"})
		return
	}

	resp, err := h.TotalSumUC.Execute(ctx, req)
	if err != nil {
		status, msg, internalErr := HttpError(err)
		h.log.Error("failed to get total sum",
			slog.Int("status", status),
			slog.String("public_msg", msg),
			slog.Any("cause", internalErr),
		)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
