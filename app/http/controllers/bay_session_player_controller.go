package controllers

import (
	"goravel/app/models"
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type BaySessionPlayerController struct {
}

func NewBaySessionPlayerController() *BaySessionPlayerController {
	return &BaySessionPlayerController{}
}

// Index returns a list of players in a bay session
func (r *BaySessionPlayerController) Index(ctx http.Context) http.Response {
	baySessionIDStr := ctx.Request().Route("id")
	baySessionID, err := strconv.ParseUint(baySessionIDStr, 10, 64)
	if err != nil {
		return ctx.Response().Status(400).Json(map[string]any{
			"error": "Invalid bay session ID",
		})
	}
	var players []models.BaySessionPlayer

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

	q := facades.Orm().Query().Where("bay_session_id = ?", baySessionID)

	// Get total count
	total, err := q.Count()
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err.Error(),
		})
	}

	// Fetch paginated results
	offset := (page - 1) * perPage
	q2 := facades.Orm().Query().Where("bay_session_id = ?", baySessionID)
	if err := q2.OrderBy("created_at", "desc").Offset(offset).Limit(perPage).Find(&players); err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err.Error(),
		})
	}

	if players == nil {
		players = []models.BaySessionPlayer{}
	}

	lastPage := int64(1)
	if total > 0 {
		lastPage = (total + int64(perPage) - 1) / int64(perPage)
	}

	return ctx.Response().Success().Json(map[string]any{
		"data": players,
		"pagination": map[string]any{
			"total":        total,
			"per_page":     perPage,
			"current_page": page,
			"last_page":    lastPage,
		},
	})
}

// Store adds a player to a bay session
func (r *BaySessionPlayerController) Store(ctx http.Context) http.Response {
	baySessionIDStr := ctx.Request().Route("id")
	baySessionID, err := strconv.ParseUint(baySessionIDStr, 10, 64)
	if err != nil {
		return ctx.Response().Status(400).Json(map[string]any{
			"error": "Invalid bay session ID",
		})
	}

	var request struct {
		PlayerID string `json:"player_id"`
	}

	if err := ctx.Request().Bind(&request); err != nil {
		return ctx.Response().Status(400).Json(map[string]any{
			"error": "Invalid request",
		})
	}

	if request.PlayerID == "" {
		return ctx.Response().Status(400).Json(map[string]any{
			"error": "player_id is required",
		})
	}

	// Verify bay session exists
	var baySession models.BaySession
	if err := facades.Orm().Query().Where("id = ?", baySessionID).First(&baySession); err != nil {
		return ctx.Response().Status(404).Json(map[string]any{
			"error": "Bay session not found",
		})
	}

	// Check if player already exists in this session
	var existingPlayer models.BaySessionPlayer
	checkErr := facades.Orm().Query().Where("bay_session_id = ? AND player_id = ?", baySessionID, request.PlayerID).First(&existingPlayer)
	if checkErr == nil {
		// Player already exists
		return ctx.Response().Status(409).Json(map[string]any{
			"error": "Player already exists in this session",
		})
	}

	player := models.BaySessionPlayer{
		BaySessionID: baySessionID,
		PlayerID:     request.PlayerID,
	}

	if err := facades.Orm().Query().Create(&player); err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err.Error(),
		})
	}

	return ctx.Response().Status(201).Json(player)
}

// Show retrieves a specific player in a bay session
func (r *BaySessionPlayerController) Show(ctx http.Context) http.Response {
	baySessionIDStr := ctx.Request().Route("id")
	baySessionID, err := strconv.ParseUint(baySessionIDStr, 10, 64)
	if err != nil {
		return ctx.Response().Status(400).Json(map[string]any{
			"error": "Invalid bay session ID",
		})
	}
	playerID := ctx.Request().Route("player_id")

	var player models.BaySessionPlayer
	if err := facades.Orm().Query().Where("bay_session_id = ? AND id = ?", baySessionID, playerID).First(&player); err != nil {
		return ctx.Response().Status(404).Json(map[string]any{
			"error": "Player not found in this session",
		})
	}

	return ctx.Response().Success().Json(player)
}

// Destroy removes a player from a bay session (soft delete)
func (r *BaySessionPlayerController) Destroy(ctx http.Context) http.Response {
	baySessionIDStr := ctx.Request().Route("id")
	baySessionID, err := strconv.ParseUint(baySessionIDStr, 10, 64)
	if err != nil {
		return ctx.Response().Status(400).Json(map[string]any{
			"error": "Invalid bay session ID",
		})
	}
	playerID := ctx.Request().Route("player_id")

	var player models.BaySessionPlayer
	if err := facades.Orm().Query().Where("bay_session_id = ? AND id = ?", baySessionID, playerID).First(&player); err != nil {
		return ctx.Response().Status(404).Json(map[string]any{
			"error": "Player not found in this session",
		})
	}

	_, delErr := facades.Orm().Query().Delete(&player)
	if delErr != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": delErr.Error(),
		})
	}

	return ctx.Response().Success().Json(map[string]any{
		"message": "Player removed from session successfully",
	})
}

// DeletePermanently permanently deletes a player record (force delete)
func (r *BaySessionPlayerController) DeletePermanently(ctx http.Context) http.Response {
	baySessionIDStr := ctx.Request().Route("id")
	baySessionID, err := strconv.ParseUint(baySessionIDStr, 10, 64)
	if err != nil {
		return ctx.Response().Status(400).Json(map[string]any{
			"error": "Invalid bay session ID",
		})
	}
	playerID := ctx.Request().Route("player_id")

	var player models.BaySessionPlayer
	if err := facades.Orm().Query().Where("bay_session_id = ? AND id = ?", baySessionID, playerID).First(&player); err != nil {
		return ctx.Response().Status(404).Json(map[string]any{
			"error": "Player not found in this session",
		})
	}

	_, forceErr := facades.Orm().Query().ForceDelete(&player)
	if forceErr != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": forceErr.Error(),
		})
	}

	return ctx.Response().Success().Json(map[string]any{
		"message": "Player permanently deleted",
	})
}

// Restore restores a soft-deleted player
func (r *BaySessionPlayerController) Restore(ctx http.Context) http.Response {
	baySessionIDStr := ctx.Request().Route("id")
	baySessionID, err := strconv.ParseUint(baySessionIDStr, 10, 64)
	if err != nil {
		return ctx.Response().Status(400).Json(map[string]any{
			"error": "Invalid bay session ID",
		})
	}
	playerID := ctx.Request().Route("player_id")

	var player models.BaySessionPlayer
	if err := facades.Orm().Query().WithTrashed().Where("bay_session_id = ? AND id = ?", baySessionID, playerID).First(&player); err != nil {
		return ctx.Response().Status(404).Json(map[string]any{
			"error": "Player not found in this session",
		})
	}

	_, restoreErr := facades.Orm().Query().Restore(&player)
	if restoreErr != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": restoreErr.Error(),
		})
	}

	return ctx.Response().Success().Json(player)
}
