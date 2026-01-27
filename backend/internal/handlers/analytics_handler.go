package handlers

import (
	"encoding/csv"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"thank-you-card-backend/internal/handlers/dto"
	"thank-you-card-backend/internal/services"
)

type AnalyticsHandler struct {
	analyticsService services.AnalyticsService
}

func NewAnalyticsHandler(analyticsService services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

func (h *AnalyticsHandler) GetDashboard(c *gin.Context) {
	from, to := parseTimeRange(c)

	ctx := c.Request.Context()
	overview, err := h.analyticsService.GetDashboardOverview(ctx, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load dashboard"})
		return
	}

	c.JSON(http.StatusOK, dto.DashboardResponse{
		TotalCards:        overview.TotalCards,
		TotalEmployees:    overview.TotalEmployees,
		TotalDepartments:  overview.TotalDepartments,
		ActiveRecognizers: overview.ActiveRecognizers,
		Period:            overview.Period,
	})
}

func (h *AnalyticsHandler) GetMostActiveRecognizers(c *gin.Context) {
	from, to := parseTimeRange(c)
	limit := parseIntQuery(c, "limit", 10)

	ctx := c.Request.Context()
	recognizers, err := h.analyticsService.GetMostActiveRecognizers(ctx, from, to, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load active recognizers"})
		return
	}

	resp := make([]dto.TopEmployeeResponse, len(recognizers))
	for i, r := range recognizers {
		resp[i] = dto.TopEmployeeResponse{
			EmployeeID: r.EmployeeID,
			Name:       r.Name,
			Department: r.Department,
			Count:      r.Count,
		}
	}

	c.JSON(http.StatusOK, gin.H{"items": resp})
}

func (h *AnalyticsHandler) GetTeamPatterns(c *gin.Context) {
	from, to := parseTimeRange(c)

	ctx := c.Request.Context()
	patterns, err := h.analyticsService.GetTeamRecognitionPatterns(ctx, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load team patterns"})
		return
	}

	resp := make([]dto.TeamPatternResponse, len(patterns))
	for i, p := range patterns {
		resp[i] = dto.TeamPatternResponse{
			Department:              p.Department,
			TotalSent:               p.TotalSent,
			TotalReceived:           p.TotalReceived,
			AverageCardsPerEmployee: p.AverageCardsPerEmployee,
		}
	}

	c.JSON(http.StatusOK, gin.H{"items": resp})
}

func (h *AnalyticsHandler) GetValuesDistribution(c *gin.Context) {
	from, to := parseTimeRange(c)

	ctx := c.Request.Context()
	distribution, err := h.analyticsService.GetValuesDistribution(ctx, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load values distribution"})
		return
	}

	var total int64
	for _, item := range distribution {
		total += item.Count
	}

	resp := make([]dto.ValueDistributionResponse, len(distribution))
	for i, item := range distribution {
		percentage := float64(0)
		if total > 0 {
			percentage = float64(item.Count) / float64(total) * 100
		}
		resp[i] = dto.ValueDistributionResponse{
			ValueID:    item.CompanyValueID,
			Code:       item.Code,
			Name:       item.Name,
			Type:       item.Type,
			Count:      item.Count,
			Percentage: percentage,
		}
	}

	c.JSON(http.StatusOK, gin.H{"items": resp, "total": total})
}

func (h *AnalyticsHandler) ExportCards(c *gin.Context) {
	filters := parseCardFilters(c)
	from, to := parseTimeRange(c)
	filters.From = &from
	filters.To = &to

	ctx := c.Request.Context()
	rows, err := h.analyticsService.ExportCards(ctx, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to export cards"})
		return
	}

	// Set CSV headers
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=cards_export.csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write header
	header := []string{
		"Card ID",
		"Sender Name",
		"Sender Email",
		"Recipient Name",
		"Recipient Email",
		"Reason",
		"Values",
		"Credos",
		"Created At",
		"Reaction Counts",
	}
	if err := writer.Write(header); err != nil {
		return
	}

	// Write data rows
	for _, row := range rows {
		record := []string{
			strconv.FormatUint(uint64(row.CardID), 10),
			row.SenderName,
			row.SenderEmail,
			row.RecipientName,
			row.RecipientEmail,
			row.Reason,
			row.Values,
			row.Credos,
			row.CreatedAt.Format(time.RFC3339),
			row.ReactionCounts,
		}
		if err := writer.Write(record); err != nil {
			return
		}
	}
}

func (h *AnalyticsHandler) GetTopRecognizedEmployees(c *gin.Context) {
	from, to := parseTimeRange(c)
	limit := parseIntQuery(c, "limit", 10)

	ctx := c.Request.Context()
	employees, err := h.analyticsService.GetTopRecognizedEmployees(ctx, from, to, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load top recognized employees"})
		return
	}

	resp := make([]dto.TopEmployeeResponse, len(employees))
	for i, e := range employees {
		resp[i] = dto.TopEmployeeResponse{
			EmployeeID: e.EmployeeID,
			Name:       e.Name,
			Department: e.Department,
			Count:      e.Count,
		}
	}

	c.JSON(http.StatusOK, gin.H{"items": resp})
}

func parseTimeRange(c *gin.Context) (time.Time, time.Time) {
	to := time.Now()
	from := to.AddDate(0, 0, -30) // default: last 30 days

	if fromStr := c.Query("from"); fromStr != "" {
		if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = t
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = t
		}
	}

	return from, to
}

func parseIntQuery(c *gin.Context, key string, defaultValue int) int {
	if str := c.Query(key); str != "" {
		if val, err := strconv.Atoi(str); err == nil && val > 0 {
			return val
		}
	}
	return defaultValue
}
