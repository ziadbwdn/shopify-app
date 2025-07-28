package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"shopify-app/internal/contract"
	"shopify-app/internal/entities"
	"shopify-app/internal/exception"
	"shopify-app/internal/utils"
	"time"
)

// reportRepository implements the contract.ReportRepository interface
type reportRepository struct {
	db *gorm.DB
}

// NewReportRepository creates a new instance of the report repository
func NewReportRepository(db *gorm.DB) contract.ReportRepository {
	return &reportRepository{db: db}
}

// GetDailySales retrieves daily sales data for a specific date
func (r *reportRepository) GetDailySales(ctx context.Context, date time.Time) (*contract.SalesReportData, *exception.AppError) {
	var result contract.SalesReportData
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := r.db.WithContext(ctx).Model(&entities.Order{}).
		Select("? as date, COALESCE(SUM(total_amount), 0) as total_sales, COUNT(id) as order_count", startOfDay).
		Where("created_at >= ? AND created_at < ? AND status = ?", startOfDay, endOfDay, entities.StatusDelivered).
		First(&result).Error

	if err != nil {
		return nil, exception.NewAppError(err, "failed to get daily sales")
	}

	// Get total items sold
	err = r.db.WithContext(ctx).Model(&entities.OrderItem{}).
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Select("COALESCE(SUM(order_items.quantity), 0) as items_sold").
		Where("orders.created_at >= ? AND orders.created_at < ? AND orders.status = ?", startOfDay, endOfDay, entities.StatusDelivered).
		Scan(&result.ItemsSold).Error

	if err != nil {
		return nil, exception.NewAppError(err, "failed to get total items sold for daily sales")
	}

	return &result, nil
}

// GetSalesByDateRange retrieves sales data for a date range
func (r *reportRepository) GetSalesByDateRange(ctx context.Context, startDate, endDate time.Time) ([]contract.SalesReportData, *exception.AppError) {
	var results []contract.SalesReportData
	err := r.db.WithContext(ctx).Model(&entities.Order{}).
		Select("DATE(created_at) as date, SUM(total_amount) as total_sales, COUNT(id) as order_count").
		Where("created_at BETWEEN ? AND ? AND status = ?", startDate, endDate, entities.StatusDelivered).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&results).Error

	if err != nil {
		return nil, exception.NewAppError(err, "failed to get sales by date range")
	}

	// This would be inefficient (N+1), better to do a separate query and map
	for i := range results {
		var itemsSold int64
		startOfDay := results[i].Date
		endOfDay := startOfDay.Add(24 * time.Hour)
		err = r.db.WithContext(ctx).Model(&entities.OrderItem{}).
			Joins("JOIN orders ON orders.id = order_items.order_id").
			Select("COALESCE(SUM(order_items.quantity), 0)").
			Where("orders.created_at >= ? AND orders.created_at < ? AND orders.status = ?", startOfDay, endOfDay, entities.StatusDelivered).
			Scan(&itemsSold).Error
		if err != nil {
			return nil, exception.NewAppError(err, "failed to get items sold for date range")
		}
		results[i].ItemsSold = itemsSold
	}

	return results, nil
}

// GetMonthlySales retrieves monthly sales data
func (r *reportRepository) GetMonthlySales(ctx context.Context, year int, month int) (*contract.SalesReportData, *exception.AppError) {
	var result contract.SalesReportData
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	err := r.db.WithContext(ctx).Model(&entities.Order{}).
		Select("? as date, COALESCE(SUM(total_amount), 0) as total_sales, COUNT(id) as order_count", startDate).
		Where("created_at >= ? AND created_at < ? AND status = ?", startDate, endDate, entities.StatusDelivered).
		First(&result).Error

	if err != nil {
		return nil, exception.NewAppError(err, "failed to get monthly sales")
	}

	err = r.db.WithContext(ctx).Model(&entities.OrderItem{}).
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Select("COALESCE(SUM(order_items.quantity), 0) as items_sold").
		Where("orders.created_at >= ? AND orders.created_at < ? AND orders.status = ?", startDate, endDate, entities.StatusDelivered).
		Scan(&result.ItemsSold).Error

	if err != nil {
		return nil, exception.NewAppError(err, "failed to get items sold for monthly sales")
	}

	return &result, nil
}

// GetYearlySales retrieves yearly sales data
func (r *reportRepository) GetYearlySales(ctx context.Context, year int) (*contract.SalesReportData, *exception.AppError) {
	var result contract.SalesReportData
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(1, 0, 0)

	err := r.db.WithContext(ctx).Model(&entities.Order{}).
		Select("? as date, COALESCE(SUM(total_amount), 0) as total_sales, COUNT(id) as order_count", startDate).
		Where("created_at >= ? AND created_at < ? AND status = ?", startDate, endDate, entities.StatusDelivered).
		First(&result).Error

	if err != nil {
		return nil, exception.NewAppError(err, "failed to get yearly sales")
	}

	err = r.db.WithContext(ctx).Model(&entities.OrderItem{}).
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Select("COALESCE(SUM(order_items.quantity), 0) as items_sold").
		Where("orders.created_at >= ? AND orders.created_at < ? AND orders.status = ?", startDate, endDate, entities.StatusDelivered).
		Scan(&result.ItemsSold).Error

	if err != nil {
		return nil, exception.NewAppError(err, "failed to get items sold for yearly sales")
	}

	return &result, nil
}

// GetBestSellingItems retrieves best selling items within a date range
func (r *reportRepository) GetBestSellingItems(ctx context.Context, startDate, endDate time.Time, limit int) ([]contract.BestSellingItem, *exception.AppError) {
	var results []contract.BestSellingItem
	err := r.db.WithContext(ctx).Model(&entities.OrderItem{}).
		Select("order_items.menu_id, menus.name as menu_name, menus.category, SUM(order_items.quantity) as total_sold, SUM(order_items.price * order_items.quantity) as total_revenue").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Joins("JOIN menus ON menus.id = order_items.menu_id").
		Where("orders.created_at BETWEEN ? AND ? AND orders.status = ?", startDate, endDate, entities.StatusDelivered).
		Group("order_items.menu_id, menus.name, menus.category").
		Order("total_sold DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		return nil, exception.NewAppError(err, "failed to get best selling items")
	}
	return results, nil
}

// GetSalesAnalytics retrieves comprehensive sales analytics for a period
func (r *reportRepository) GetSalesAnalytics(ctx context.Context, startDate, endDate time.Time) (*contract.SalesAnalytics, *exception.AppError) {
	var analytics contract.SalesAnalytics
	analytics.PeriodStart = startDate
	analytics.PeriodEnd = endDate

	err := r.db.WithContext(ctx).Model(&entities.Order{}).
		Select("COALESCE(SUM(total_amount), 0) as total_revenue, COUNT(id) as total_orders").
		Where("created_at BETWEEN ? AND ? AND status = ?", startDate, endDate, entities.StatusDelivered).
		First(&analytics).Error
	if err != nil {
		return nil, exception.NewAppError(err, "failed to get sales analytics revenue and orders")
	}

	err = r.db.WithContext(ctx).Model(&entities.OrderItem{}).
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Select("COALESCE(SUM(order_items.quantity), 0) as total_items").
		Where("orders.created_at BETWEEN ? AND ? AND orders.status = ?", startDate, endDate, entities.StatusDelivered).
		Scan(&analytics.TotalItems).Error
	if err != nil {
		return nil, exception.NewAppError(err, "failed to get sales analytics total items")
	}

	if analytics.TotalOrders > 0 {
		totalRevenue, _ := utils.GormDecimalToFloat64(*analytics.TotalRevenue)
		avg, _ := utils.Float64ToGormDecimal(totalRevenue / float64(analytics.TotalOrders))
		analytics.AverageOrderValue = avg
	} else {
		analytics.AverageOrderValue = utils.MustNewGormDecimal("0")
	}

	err = r.db.WithContext(ctx).Model(&entities.OrderItem{}).
		Select("menus.category").
		Joins("JOIN menus ON menus.id = order_items.menu_id").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.created_at BETWEEN ? AND ? AND orders.status = ?", startDate, endDate, entities.StatusDelivered).
		Group("menus.category").
		Order("SUM(order_items.quantity) DESC").
		Limit(1).
		Scan(&analytics.TopCategory).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.NewAppError(err, "failed to get sales analytics top category")
	}

	return &analytics, nil
}

// GetSalesByCategory retrieves sales data grouped by category
func (r *reportRepository) GetSalesByCategory(ctx context.Context, startDate, endDate time.Time) (map[string]*utils.GormDecimal, *exception.AppError) {
	type CategorySale struct {
		Category string
		Total    *utils.GormDecimal
	}
	var results []CategorySale
	err := r.db.WithContext(ctx).Model(&entities.OrderItem{}).
		Select("menus.category, SUM(order_items.price * order_items.quantity) as total").
		Joins("JOIN menus ON menus.id = order_items.menu_id").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.created_at BETWEEN ? AND ? AND orders.status = ?", startDate, endDate, entities.StatusDelivered).
		Group("menus.category").
		Scan(&results).Error

	if err != nil {
		return nil, exception.NewAppError(err, "failed to get sales by category")
	}

	salesMap := make(map[string]*utils.GormDecimal)
	for _, res := range results {
		salesMap[res.Category] = res.Total
	}

	return salesMap, nil
}

// GetCustomerOrderStats retrieves customer ordering statistics
func (r *reportRepository) GetCustomerOrderStats(ctx context.Context, startDate, endDate time.Time) (map[string]interface{}, *exception.AppError) {
	// This is a complex query, placeholder for now
	// You would typically calculate new vs returning customers, etc.
	return map[string]interface{}{"status": "not implemented"}, nil
}

// GetHourlyOrderDistribution retrieves order distribution by hour
func (r *reportRepository) GetHourlyOrderDistribution(ctx context.Context, startDate, endDate time.Time) (map[int]int64, *exception.AppError) {
	type HourlyResult struct {
		Hour  int
		Count int64
	}
	var results []HourlyResult
	err := r.db.WithContext(ctx).Model(&entities.Order{}).
		Select("EXTRACT(HOUR FROM created_at) as hour, COUNT(id) as count").
		Where("created_at BETWEEN ? AND ? AND status = ?", startDate, endDate, entities.StatusDelivered).
		Group("hour").
		Scan(&results).Error

	if err != nil {
		return nil, exception.NewAppError(err, "failed to get hourly order distribution")
	}

	distMap := make(map[int]int64)
	for _, res := range results {
		distMap[res.Hour] = res.Count
	}

	return distMap, nil
}

// GetRevenueGrowth calculates revenue growth between two periods
func (r *reportRepository) GetRevenueGrowth(ctx context.Context, currentStart, currentEnd, previousStart, previousEnd time.Time) (float64, *exception.AppError) {
	var currentRevenue, previousRevenue utils.GormDecimal

	err := r.db.WithContext(ctx).Model(&entities.Order{}).
		Select("COALESCE(SUM(total_amount), 0)").
		Where("created_at BETWEEN ? AND ? AND status = ?", currentStart, currentEnd, entities.StatusDelivered).
		Scan(&currentRevenue).Error
	if err != nil {
		return 0, exception.NewAppError(err, "failed to get current period revenue for growth calculation")
	}

	err = r.db.WithContext(ctx).Model(&entities.Order{}).
		Select("COALESCE(SUM(total_amount), 0)").
		Where("created_at BETWEEN ? AND ? AND status = ?", previousStart, previousEnd, entities.StatusDelivered).
		Scan(&previousRevenue).Error
	if err != nil {
		return 0, exception.NewAppError(err, "failed to get previous period revenue for growth calculation")
	}

	current, _ := utils.GormDecimalToFloat64(currentRevenue)
	previous, _ := utils.GormDecimalToFloat64(previousRevenue)

	if previous == 0 {
		if current > 0 {
			return 100.0, nil // Infinite growth, represent as 100%
		}
		return 0.0, nil
	}

	growth := ((current - previous) / previous) * 100
	return growth, nil
}
