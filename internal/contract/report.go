// internal/contract/report_contract.go
package contract

import (
	"context"
	"shopify-app/internal/exception"
	"shopify-app/internal/utils"
	"time"
)

// SalesReportData represents aggregated sales data
type SalesReportData struct {
	Date        time.Time          `json:"date"`
	TotalSales  *utils.GormDecimal `json:"total_sales"`
	OrderCount  int64              `json:"order_count"`
	ItemsSold   int64              `json:"items_sold"`
}

// BestSellingItem represents best selling menu item data
type BestSellingItem struct {
	MenuID      utils.BinaryUUID   `json:"menu_id"`
	MenuName    string             `json:"menu_name"`
	Category    string             `json:"category"`
	TotalSold   int64              `json:"total_sold"`
	TotalRevenue *utils.GormDecimal `json:"total_revenue"`
}

// SalesAnalytics represents comprehensive sales analytics
type SalesAnalytics struct {
	TotalRevenue    *utils.GormDecimal `json:"total_revenue"`
	TotalOrders     int64              `json:"total_orders"`
	TotalItems      int64              `json:"total_items"`
	AverageOrderValue *utils.GormDecimal `json:"average_order_value"`
	TopCategory     string             `json:"top_category"`
	PeriodStart     time.Time          `json:"period_start"`
	PeriodEnd       time.Time          `json:"period_end"`
}

// ReportRepository defines the contract for report data access operations
type ReportRepository interface {
	// GetDailySales retrieves daily sales data for a specific date
	GetDailySales(ctx context.Context, date time.Time) (*SalesReportData, *exception.AppError)
	
	// GetSalesByDateRange retrieves sales data for a date range
	GetSalesByDateRange(ctx context.Context, startDate, endDate time.Time) ([]SalesReportData, *exception.AppError)
	
	// GetMonthlySales retrieves monthly sales data for a specific month
	GetMonthlySales(ctx context.Context, year int, month int) (*SalesReportData, *exception.AppError)
	
	// GetYearlySales retrieves yearly sales data for a specific year
	GetYearlySales(ctx context.Context, year int) (*SalesReportData, *exception.AppError)
	
	// GetBestSellingItems retrieves best selling items within a date range
	GetBestSellingItems(ctx context.Context, startDate, endDate time.Time, limit int) ([]BestSellingItem, *exception.AppError)
	
	// GetSalesAnalytics retrieves comprehensive sales analytics for a period
	GetSalesAnalytics(ctx context.Context, startDate, endDate time.Time) (*SalesAnalytics, *exception.AppError)
	
	// GetSalesByCategory retrieves sales data grouped by category
	GetSalesByCategory(ctx context.Context, startDate, endDate time.Time) (map[string]*utils.GormDecimal, *exception.AppError)
	
	// GetCustomerOrderStats retrieves customer ordering statistics
	GetCustomerOrderStats(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, *exception.AppError)
	
	// GetHourlyOrderDistribution retrieves order distribution by hour
	GetHourlyOrderDistribution(ctx context.Context, startDate, endDate time.Time) (map[int]int64, *exception.AppError)
	
	// GetRevenueGrowth calculates revenue growth between periods
	GetRevenueGrowth(ctx context.Context, currentStart, currentEnd, previousStart, previousEnd time.Time) (float64, *exception.AppError)
}

// ReportService defines the contract for report business logic operations
type ReportService interface {
	// GenerateSalesReport generates comprehensive sales report for a date range
	GenerateSalesReport(ctx context.Context, startDate, endDate time.Time, groupBy string) ([]SalesReportData, *exception.AppError)
	
	// GetDailySalesReport retrieves daily sales report for a specific date
	GetDailySalesReport(ctx context.Context, date time.Time) (*SalesReportData, *exception.AppError)
	
	// GetMonthlySalesReport retrieves monthly sales report
	GetMonthlySalesReport(ctx context.Context, year int, month int) (*SalesReportData, *exception.AppError)
	
	// GetYearlySalesReport retrieves yearly sales report
	GetYearlySalesReport(ctx context.Context, year int) (*SalesReportData, *exception.AppError)
	
	// GetBestSellingItemsReport retrieves best selling items report
	GetBestSellingItemsReport(ctx context.Context, startDate, endDate time.Time, limit int) ([]BestSellingItem, *exception.AppError)
	
	// GetSalesAnalyticsReport retrieves comprehensive sales analytics
	GetSalesAnalyticsReport(ctx context.Context, startDate, endDate time.Time) (*SalesAnalytics, *exception.AppError)
	
	// GetCategorySalesReport retrieves sales report grouped by category
	GetCategorySalesReport(ctx context.Context, startDate, endDate time.Time) (map[string]*utils.GormDecimal, *exception.AppError)
	
	// GetCustomerInsightsReport retrieves customer behavior insights
	GetCustomerInsightsReport(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, *exception.AppError)
	
	// GetOrderTrendsReport retrieves order trends and patterns
	GetOrderTrendsReport(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, *exception.AppError)
	
	// GetRevenueGrowthReport calculates and returns revenue growth metrics
	GetRevenueGrowthReport(ctx context.Context, currentStart, currentEnd, previousStart, previousEnd time.Time) (map[string]interface{}, *exception.AppError)
	
	// ExportSalesReportCSV exports sales report data to CSV format
	ExportSalesReportCSV(ctx context.Context, startDate, endDate time.Time) ([]byte, *exception.AppError)
	
	// ValidateReportDateRange validates date range parameters
	ValidateReportDateRange(startDate, endDate time.Time) *exception.AppError
}