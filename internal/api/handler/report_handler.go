package handler

import (
	"shopify-app/internal/contract"
	"shopify-app/pkg/web_response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportService contract.ReportService
}

func NewReportHandler(reportService contract.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

func (h *ReportHandler) GetSalesReport(c *gin.Context) {
	startDate, err := time.Parse("2006-01-02", c.Query("start_date"))
	if err != nil {
		web_response.HandleError(c, err)
		return
	}
	endDate, err := time.Parse("2006-01-02", c.Query("end_date"))
	if err != nil {
		web_response.HandleError(c, err)
		return
	}
	groupBy := c.Query("group_by")

	report, appErr := h.reportService.GenerateSalesReport(c.Request.Context(), startDate, endDate, groupBy)
	if appErr != nil {
		web_response.HandleError(c, appErr)
		return
	}
	web_response.Success(c, report)
}

func (h *ReportHandler) GetBestSellingItems(c *gin.Context) {
	startDate, err := time.Parse("2006-01-02", c.Query("start_date"))
	if err != nil {
		web_response.HandleError(c, err)
		return
	}
	endDate, err := time.Parse("2006-01-02", c.Query("end_date"))
	if err != nil {
		web_response.HandleError(c, err)
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))

	report, appErr := h.reportService.GetBestSellingItemsReport(c.Request.Context(), startDate, endDate, limit)
	if appErr != nil {
		web_response.HandleError(c, appErr)
		return
	}
	web_response.Success(c, report)
}
