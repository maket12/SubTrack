package rest

import (
	"SubTrack/app/dto"
	"SubTrack/app/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	CreateUC   *usecase.CreateSubscriptionUC
	GetUC      *usecase.GetSubscriptionUC
	UpdateUC   *usecase.UpdateSubscriptionUC
	DeleteUC   *usecase.DeleteSubscriptionUC
	ListUC     *usecase.GetSubscriptionListUC
	TotalSumUC *usecase.GetTotalSumUC
}

func NewSubscriptionHandler(
	createUC *usecase.CreateSubscriptionUC,
	getUC *usecase.GetSubscriptionUC,
	updateUC *usecase.UpdateSubscriptionUC,
	deleteUC *usecase.DeleteSubscriptionUC,
	listUC *usecase.GetSubscriptionListUC,
	totalSumUC *usecase.GetTotalSumUC,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		CreateUC:   createUC,
		GetUC:      getUC,
		UpdateUC:   updateUC,
		DeleteUC:   deleteUC,
		ListUC:     listUC,
		TotalSumUC: totalSumUC,
	}
}

// @Summary      Create subscription
// @Description  Creates a new user subscription
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        input  body      dto.CreateSubscription           true  "Subscription data"
// @Success      201    {object}  dto.CreateSubscriptionResponse
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /subscriptions [post]
func (h *SubscriptionHandler) Create(ctx *gin.Context) {
	var req dto.CreateSubscription
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	resp, err := h.CreateUC.Execute(ctx, req)
	if err != nil {
		status, msg := HttpError(err)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// @Summary      Get subscription by ID
// @Description  Returns a single subscription by ID
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      int  true  "Subscription ID"
// @Success      200  {object}  dto.GetSubscriptionResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id must be positive integer or 0"})
		return
	}

	resp, err := h.GetUC.Execute(ctx, dto.GetSubscription{ID: id})
	if err != nil {
		status, msg := HttpError(err)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary      Update subscription
// @Description  Updates subscription by ID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id     path      int                       true  "Subscription ID"
// @Param        input  body      dto.UpdateSubscription     true  "Subscription update data"
// @Success      200    {object}  dto.UpdateSubscriptionResponse
// @Failure      400    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Router       /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id must be positive integer or 0"})
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
		status, msg := HttpError(err)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary      Delete subscription
// @Description  Deletes subscription by ID
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      int  true  "Subscription ID"
// @Success      200  {object}  dto.DeleteSubscriptionResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id must be positive integer or 0"})
		return
	}

	resp, err := h.DeleteUC.Execute(ctx, dto.DeleteSubscription{ID: id})
	if err != nil {
		status, msg := HttpError(err)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary      List subscriptions
// @Description  Returns paginated list of subscriptions
// @Tags         subscriptions
// @Produce      json
// @Param        user_id      query     string  false  "User ID (UUID)"
// @Param        service_name query     string  false  "Service name"
// @Param        limit        query     int     false  "Limit (default 10)"
// @Param        offset       query     int     false  "Offset"
// @Success      200  {object} dto.GetSubscriptionListResponse
// @Failure      400  {object} map[string]string
// @Router       /subscriptions [get]
func (h *SubscriptionHandler) List(ctx *gin.Context) {
	var req dto.GetSubscriptionList
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid query"})
		return
	}

	resp, err := h.ListUC.Execute(ctx, req)
	if err != nil {
		status, msg := HttpError(err)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary      Get total subscription cost
// @Description  Calculates total sum of subscriptions for given filters
// @Tags         subscriptions
// @Produce      json
// @Param        user_id      query     string   false  "User ID (UUID)"
// @Param        service_name query     string   false  "Service name"
// @Param        start_date   query     string   false  "Start date (DD-MM-YYYY)"
// @Param        end_date     query     string   false  "End date (DD-MM-YYYY)"
// @Success      200  {object} dto.GetTotalSumResponse
// @Failure      400  {object} map[string]string
// @Router       /subscriptions/total [get]
func (h *SubscriptionHandler) GetTotalSum(ctx *gin.Context) {
	var req dto.GetTotalSum
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid query"})
		return
	}

	resp, err := h.TotalSumUC.Execute(ctx, req)
	if err != nil {
		status, msg := HttpError(err)
		ctx.JSON(status, gin.H{"error": msg})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
