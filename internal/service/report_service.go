package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"shopify-app/internal/contract"
	"shopify-app/internal/exception"
	"shopify-app/internal/utils"
	"strconv"
	"time"
)

type reportService struct {
	reportRepo contract.ReportRepository
}

func NewReportService(reportRepo contract.ReportRepository) contract.ReportService {
	return &reportService{reportRepo: reportRepo}
}

func (s *reportService) GenerateSalesReport(ctx context.Context, startDate, endDate time.Time, groupBy string) ([]contract.SalesReportData, *exception.AppError) {
	if err := s.ValidateReportDateRange(startDate, endDate); err != nil {
		return nil, err
	}
	// Simple implementation, ignores groupBy for now
	return s.reportRepo.GetSalesByDateRange(ctx, startDate, endDate)
}

func (s *reportService) GetDailySalesReport(ctx context.Context, date time.Time) (*contract.SalesReportData, *exception.AppError) {
	return s.reportRepo.GetDailySales(ctx, date)
}

func (s *reportService) GetMonthlySalesReport(ctx context.Context, year int, month int) (*contract.SalesReportData, *exception.AppError) {
	return s.reportRepo.GetMonthlySales(ctx, year, month)
}

func (s *reportService) GetYearlySalesReport(ctx context.Context, year int) (*contract.SalesReportData, *exception.AppError) {
	return s.reportRepo.GetYearlySales(ctx, year)
}

func (s *reportService) GetBestSellingItemsReport(ctx context.Context, startDate, endDate time.Time, limit int) ([]contract.BestSellingItem, *exception.AppError) {
	if err := s.ValidateReportDateRange(startDate, endDate); err != nil {
		return nil, err
	}
	return s.reportRepo.GetBestSellingItems(ctx, startDate, endDate, limit)
}

func (s *reportService) GetSalesAnalyticsReport(ctx context.Context, startDate, endDate time.Time) (*contract.SalesAnalytics, *exception.AppError) {
	if err := s.ValidateReportDateRange(startDate, endDate); err != nil {
		return nil, err
	}
	return s.reportRepo.GetSalesAnalytics(ctx, startDate, endDate)
}

func (s *reportService) GetCategorySalesReport(ctx context.Context, startDate, endDate time.Time) (map[string]*utils.GormDecimal, *exception.AppError) {
	if err := s.ValidateReportDateRange(startDate, endDate); err != nil {
		return nil, err
	}
	return s.reportRepo.GetSalesByCategory(ctx, startDate, endDate)
}

func (s *reportService) GetCustomerInsightsReport(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, *exception.AppError) {
	if err := s.ValidateReportDateRange(startDate, endDate); err != nil {
		return nil, err
	}
	return s.reportRepo.GetCustomerOrderStats(ctx, startDate, endDate)
}

func (s *reportService) GetOrderTrendsReport(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, *exception.AppError) {
	if err := s.ValidateReportDateRange(startDate, endDate); err != nil {
		return nil, err
	}
	dist, err := s.reportRepo.GetHourlyOrderDistribution(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"hourly_distribution": dist}, nil
}

func (s *reportService) GetRevenueGrowthReport(ctx context.Context, currentStart, currentEnd, previousStart, previousEnd time.Time) (map[string]interface{}, *exception.AppError) {
	growth, err := s.reportRepo.GetRevenueGrowth(ctx, currentStart, currentEnd, previousStart, previousEnd)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"revenue_growth_percentage": growth}, nil
}

func (s *reportService) ExportSalesReportCSV(ctx context.Context, startDate, endDate time.Time) ([]byte, *exception.AppError) {
	salesData, err := s.reportRepo.GetSalesByDateRange(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	w := csv.NewWriter(&b)

	// Write header
	if err := w.Write([]string{"Date", "TotalSales", "OrderCount", "ItemsSold"}); err != nil {
		return nil, exception.NewAppError(err, "failed to write csv header")
	}

	// Write rows
	for _, record := range salesData {
		row := []string{
			record.Date.Format("2006-01-02"),
			utils.GormDecimalToString(*record.TotalSales),
			strconv.FormatInt(record.OrderCount, 10),
			strconv.FormatInt(record.ItemsSold, 10),
		}
		if err := w.Write(row); err != nil {
			return nil, exception.NewAppError(err, "failed to write csv row")
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, exception.NewAppError(err, "failed to flush csv writer")
	}

	return b.Bytes(), nil
}

func (s *reportService) ValidateReportDateRange(startDate, endDate time.Time) *exception.AppError {
	if startDate.IsZero() || endDate.IsZero() {
		return exception.NewAppError(nil, "start and end dates are required")
	}
	if startDate.After(endDate) {
		return exception.NewAppError(nil, "start date cannot be after end date")
	}
	return nil
}
