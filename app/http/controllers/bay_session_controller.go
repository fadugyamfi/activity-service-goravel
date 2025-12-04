package controllers

import (
	"goravel/app/models"
	"strconv"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type BaySessionController struct {
}

func NewBaySessionController() *BaySessionController {
	return &BaySessionController{}
}

// Index returns a list of bay sessions with pagination
func (r *BaySessionController) Index(ctx http.Context) http.Response {
	var baySessions []models.BaySession

	// Get pagination parameters
	page := 1
	if p := ctx.Request().Query("page"); p != "" {
		if pageNum, err := strconv.Atoi(p); err == nil && pageNum > 0 {
			page = pageNum
		}
	}

	perPage := 15
	if pp := ctx.Request().Query("per_page"); pp != "" {
		if perPageNum, err := strconv.Atoi(pp); err == nil && perPageNum > 0 {
			perPage = perPageNum
		}
	}

	q := facades.Orm().Query()

	// Filter by visit_id if provided
	if visitID := ctx.Request().Query("visit_id"); visitID != "" {
		q = q.Where("visit_id = ?", visitID)
	}

	// Get total count
	total, err := q.Table("bay_sessions").Count()
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err.Error(),
		})
	}

	// Fetch paginated results
	offset := (page - 1) * perPage
	if err := q.OrderBy("created_at", "desc").Offset(offset).Limit(perPage).Find(&baySessions); err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err.Error(),
		})
	}

	if baySessions == nil {
		baySessions = []models.BaySession{}
	}

	lastPage := int64(1)
	if total > 0 {
		lastPage = (total + int64(perPage) - 1) / int64(perPage)
	}

	return ctx.Response().Success().Json(map[string]any{
		"data": baySessions,
		"pagination": map[string]any{
			"total":        total,
			"per_page":     perPage,
			"current_page": page,
			"last_page":    lastPage,
		},
	})
}

// Store creates a new bay session
func (r *BaySessionController) Store(ctx http.Context) http.Response {
	var request struct {
		VisitID   string     `json:"visit_id"`
		StartTime *time.Time `json:"start_time"`
		Duration  int        `json:"duration"`
	}

	if err := ctx.Request().Bind(&request); err != nil {
		return ctx.Response().Status(400).Json(map[string]any{
			"error": "Invalid request",
		})
	}

	// Validate required fields
	if request.VisitID == "" || request.Duration == 0 || request.StartTime == nil {
		return ctx.Response().Status(400).Json(map[string]any{
			"error": "visit_id, start_time, and duration are required",
		})
	}

	baySession := models.BaySession{
		VisitID:   request.VisitID,
		StartTime: *request.StartTime,
		Duration:  request.Duration,
	}

	if err := facades.Orm().Query().Create(&baySession); err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err.Error(),
		})
	}

	return ctx.Response().Status(201).Json(baySession)
}

// Show retrieves a specific bay session by ID
func (r *BaySessionController) Show(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	var baySession models.BaySession

	err := facades.Orm().Query().Where("id = ?", id).First(&baySession)
	if err != nil {
		return ctx.Response().Status(404).Json(map[string]any{
			"error": "Bay session not found",
		})
	}

	return ctx.Response().Success().Json(map[string]any{
		"data": baySession,
	})
}

// Update updates a bay session
func (r *BaySessionController) Update(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	var request struct {
		VisitID   string     `json:"visit_id"`
		StartTime *time.Time `json:"start_time"`
		Duration  int        `json:"duration"`
	}

	if err := ctx.Request().Bind(&request); err != nil {
		return ctx.Response().Status(400).Json(map[string]any{
			"error": "Invalid request",
		})
	}

	var baySession models.BaySession
	if err := facades.Orm().Query().Where("id = ?", id).First(&baySession); err != nil {
		return ctx.Response().Status(404).Json(map[string]any{
			"error": "Bay session not found",
		})
	}

	if request.VisitID != "" {
		baySession.VisitID = request.VisitID
	}
	if request.Duration != 0 {
		baySession.Duration = request.Duration
	}
	if request.StartTime != nil {
		baySession.StartTime = *request.StartTime
	}

	_, err := facades.Orm().Query().Model(&baySession).Update("visit_id", baySession.VisitID)
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err.Error(),
		})
	}

	_, err = facades.Orm().Query().Model(&baySession).Update("start_time", baySession.StartTime)
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err.Error(),
		})
	}

	_, err = facades.Orm().Query().Model(&baySession).Update("duration", baySession.Duration)
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err.Error(),
		})
	}

	// Refresh to get updated values
	facades.Orm().Query().Where("id = ?", id).First(&baySession)

	return ctx.Response().Success().Json(baySession)
}

// Destroy deletes a bay session
func (r *BaySessionController) Destroy(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	var baySession models.BaySession
	if err := facades.Orm().Query().Where("id = ?", id).First(&baySession); err != nil {
		return ctx.Response().Status(404).Json(map[string]any{
			"error": "Bay session not found",
		})
	}

	_, err := facades.Orm().Query().Delete(&baySession)
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err.Error(),
		})
	}

	return ctx.Response().Success().Json(map[string]any{
		"message": "Bay session deleted successfully",
	})
}
