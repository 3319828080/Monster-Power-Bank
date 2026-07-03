package data

import (
	"context"
	"errors"
	"math"
	"monster/internal/biz"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

// ---------- Order Model ----------

type order struct {
	ID              int64          `gorm:"primaryKey;autoIncrement"`
	OrderNo         string         `gorm:"column:order_no;uniqueIndex;size:64"`
	UserID          int64          `gorm:"index"`
	PowerBankID     int64          `gorm:"column:power_bank_id"`
	PowerBankNo     string         `gorm:"column:power_bank_no;size:64"`
	BorrowStationID int64          `gorm:"column:borrow_station_id"`
	BorrowCabinetID int64          `gorm:"column:borrow_cabinet_id"`
	BorrowSlotID    int64          `gorm:"column:borrow_slot_id"`
	ReturnStationID int64          `gorm:"column:return_station_id;default:0"`
	ReturnCabinetID int64          `gorm:"column:return_cabinet_id;default:0"`
	ReturnSlotID    int64          `gorm:"column:return_slot_id;default:0"`
	BorrowTime      time.Time      `gorm:"column:borrow_time"`
	ReturnTime      *time.Time     `gorm:"column:return_time"`
	DurationMinutes int32          `gorm:"column:duration_minutes;default:0"`
	StartFee        int64          `gorm:"column:start_fee;default:0"`
	HourlyFee       int64          `gorm:"column:hourly_fee;default:0"`
	DailyCap        int64          `gorm:"column:daily_cap;default:0"`
	Deposit         int64          `gorm:"default:0"`
	TotalAmount     int64          `gorm:"column:total_amount;default:0"`
	PaidAmount      int64          `gorm:"column:paid_amount;default:0"`
	DiscountAmount  int64          `gorm:"column:discount_amount;default:0"`
	Status          string         `gorm:"size:32;index;default:pending"`
	Remark          string         `gorm:"size:255"`
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (order) TableName() string { return "orders" }

type orderItem struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	OrderID     int64     `gorm:"index"`
	OrderNo     string    `gorm:"index;size:64"`
	FeeType     string    `gorm:"column:fee_type;size:32"`
	Amount      int64     `gorm:"default:0"`
	Description string    `gorm:"size:255"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (orderItem) TableName() string { return "order_items" }

// ---------- PowerBank Model ----------

type powerBank struct {
	ID           int64  `gorm:"primaryKey;autoIncrement"`
	DeviceNo     string `gorm:"column:device_no;uniqueIndex;size:64"`
	StationID    int64  `gorm:"column:station_id;index"`
	CabinetID    int64  `gorm:"column:cabinet_id"`
	SlotID       int64  `gorm:"column:slot_id"`
	Status       string `gorm:"size:32;default:空闲"` // 空闲, 租借中, offline, fault
	BatteryLevel int32  `gorm:"column:battery_level;default:100"`
}

func (powerBank) TableName() string { return "power_banks" }

// ---------- Station Model ----------

type station struct {
	ID          int64   `gorm:"primaryKey;autoIncrement"`
	Name        string  `gorm:"size:128"`
	Address     string  `gorm:"size:255"`
	Latitude    float64 `gorm:"column:latitude;default:0"`
	Longitude   float64 `gorm:"column:longitude;default:0"`
	Status      int32   `gorm:"default:1"` // 1-normal, 2-maintenance
	OpenTime    string  `gorm:"column:open_time;size:64"`
	Description string  `gorm:"column:description;size:512"`
	Images      string  `gorm:"column:images;size:1024"`
}

func (station) TableName() string { return "stations" }

// ---------- Cabinet Model ----------

type cabinet struct {
	ID        int64  `gorm:"primaryKey;autoIncrement"`
	StationID int64  `gorm:"column:station_id;index"`
	CabinetNo string `gorm:"column:cabinet_no;size:32"`
	Status    int32  `gorm:"default:1"` // 1-normal, 2-offline, 3-fault
}

func (cabinet) TableName() string { return "cabinets" }

// ---------- Slot Model ----------

type slot struct {
	ID          int64  `gorm:"primaryKey;autoIncrement"`
	CabinetID   int64  `gorm:"column:cabinet_id;index"`
	StationID   int64  `gorm:"column:station_id;index"`
	SlotNo      string `gorm:"column:slot_no;size:32"`
	Status      string `gorm:"size:32;default:已借出"` // 空闲, 已借出, disabled
	PowerBankID int64  `gorm:"column:power_bank_id;default:0"`
	OrderNo     string `gorm:"column:order_no;size:64"`
}

func (slot) TableName() string { return "slots" }

// ---------- Coupon Model ----------

type coupon struct {
	ID        int64      `gorm:"primaryKey;autoIncrement"`
	UserID    int64      `gorm:"column:user_id;index"`
	Type      string     `gorm:"size:32"`   // full_reduce, percent, fixed
	Amount    int64      `gorm:"default:0"` // discount amount (cents) or percentage
	MinAmount int64      `gorm:"column:min_amount;default:0"`
	Status    string     `gorm:"size:32;default:unused"` // unused, used, expired
	ExpireAt  time.Time  `gorm:"column:expire_at"`
	UsedAt    *time.Time `gorm:"column:used_at"`
}

func (coupon) TableName() string { return "coupons" }

// ---------- OrderRepo ----------

type orderRepo struct {
	data *Data
	log  *log.Helper
}

func NewOrderRepo(data *Data, logger log.Logger) biz.OrderRepo {
	return &orderRepo{data: data, log: log.NewHelper(logger)}
}

// ----- Basic CRUD -----

func (r *orderRepo) Create(ctx context.Context, o *biz.Order) (*biz.Order, error) {
	m := r.toModel(o)
	if err := r.data.DB(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return r.toBiz(m), nil
}

func (r *orderRepo) GetByID(ctx context.Context, id int64) (*biz.Order, error) {
	var m order
	if err := r.data.DB(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return r.toBiz(&m), nil
}

func (r *orderRepo) GetByOrderNo(ctx context.Context, orderNo string) (*biz.Order, error) {
	var m order
	if err := r.data.DB(ctx).Where("order_no = ?", orderNo).First(&m).Error; err != nil {
		return nil, err
	}
	return r.toBiz(&m), nil
}

func (r *orderRepo) ListByUserID(ctx context.Context, userID int64, page, pageSize int, status string) ([]*biz.Order, int, error) {
	query := r.data.DB(ctx).Model(&order{}).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	var total int64
	query.Count(&total)
	var list []order
	offset := (page - 1) * pageSize
	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	result := make([]*biz.Order, 0, len(list))
	for _, m := range list {
		result = append(result, r.toBiz(&m))
	}
	return result, int(total), nil
}

func (r *orderRepo) Update(ctx context.Context, o *biz.Order) error {
	return r.data.DB(ctx).Save(r.toModel(o)).Error
}

func (r *orderRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	return r.data.DB(ctx).Model(&order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *orderRepo) CreateItem(ctx context.Context, item *biz.OrderItem) error {
	return r.data.DB(ctx).Create(&orderItem{
		OrderID:     item.OrderID,
		OrderNo:     item.OrderNo,
		FeeType:     item.FeeType,
		Amount:      item.Amount,
		Description: item.Description,
	}).Error
}

// ----- Extended Operations -----

func (r *orderRepo) GetCurrentByUserID(ctx context.Context, userID int64) (*biz.Order, error) {
	var m order
	if err := r.data.DB(ctx).
		Where("user_id = ? AND status IN ('pending','租借中')", userID).
		Order("id DESC").First(&m).Error; err != nil {
		return nil, err
	}
	return r.toBiz(&m), nil
}

func (r *orderRepo) UpdateOrderPartial(ctx context.Context, id int64, updates map[string]interface{}) error {
	return r.data.DB(ctx).Model(&order{}).Where("id = ?", id).Updates(updates).Error
}

func (r *orderRepo) GetOrderByNoForUpdate(ctx context.Context, orderNo string) (*biz.Order, error) {
	var m order
	tx := r.data.DB(ctx).Set("gorm:query_option", "FOR UPDATE")
	if err := tx.Where("order_no = ?", orderNo).First(&m).Error; err != nil {
		return nil, err
	}
	return r.toBiz(&m), nil
}

// ----- PowerBank -----

func (r *orderRepo) GetPowerBank(ctx context.Context, id int64) (*biz.PowerBank, error) {
	var m powerBank
	if err := r.data.DB(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return r.toPowerBankBiz(&m), nil
}

func (r *orderRepo) UpdatePowerBankStatus(ctx context.Context, id int64, status string) error {
	return r.data.DB(ctx).Model(&powerBank{}).Where("id = ?", id).Update("status", status).Error
}

// ----- Station -----

func (r *orderRepo) GetStation(ctx context.Context, id int64) (*biz.Station, error) {
	var m station
	if err := r.data.DB(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &biz.Station{ID: m.ID, Name: m.Name, Address: m.Address, Latitude: m.Latitude, Longitude: m.Longitude, OpenTime: m.OpenTime, Description: m.Description, Images: m.Images}, nil
}

func (r *orderRepo) ListNearbyStations(ctx context.Context, lat, lng float64, radiusMeters int) ([]*biz.Station, error) {
	latDelta := float64(radiusMeters) / 111320.0
	lngDelta := float64(radiusMeters) / (111320.0 * cosDeg(lat))

	minLat := lat - latDelta
	maxLat := lat + latDelta
	minLng := lng - lngDelta
	maxLng := lng + lngDelta

	var list []station
	if err := r.data.DB(ctx).
		Where("latitude BETWEEN ? AND ? AND longitude BETWEEN ? AND ? AND status = 1", minLat, maxLat, minLng, maxLng).
		Find(&list).Error; err != nil {
		return nil, err
	}

	type stationDist struct {
		station  *biz.Station
		distance float64
	}
	var sorted []stationDist
	for _, m := range list {
		dist := haversine(lat, lng, m.Latitude, m.Longitude)
		sorted = append(sorted, stationDist{
			station:  &biz.Station{ID: m.ID, Name: m.Name, Address: m.Address, Latitude: m.Latitude, Longitude: m.Longitude, Distance: dist, OpenTime: m.OpenTime, Description: m.Description, Images: m.Images},
			distance: dist,
		})
	}
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i].distance > sorted[j].distance {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	result := make([]*biz.Station, 0, len(sorted))
	for _, s := range sorted {
		result = append(result, s.station)
	}
	return result, nil
}

func cosDeg(deg float64) float64 {
	return math.Cos(deg * math.Pi / 180)
}

func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371000.0
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

// ----- Slot -----

func (r *orderRepo) GetSlot(ctx context.Context, id int64) (*biz.Slot, error) {
	var m slot
	if err := r.data.DB(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return r.toSlotBiz(&m), nil
}

func (r *orderRepo) GetSlotByStationCabinet(ctx context.Context, stationID, cabinetID, slotID int64) (*biz.Slot, error) {
	var m slot
	if err := r.data.DB(ctx).
		Where("station_id = ? AND cabinet_id = ? AND id = ?", stationID, cabinetID, slotID).
		First(&m).Error; err != nil {
		return nil, err
	}
	return r.toSlotBiz(&m), nil
}

func (r *orderRepo) OccupySlot(ctx context.Context, slotID, powerBankID int64) error {
	result := r.data.DB(ctx).Model(&slot{}).
		Where("id = ? AND status = '已借出'", slotID).
		Updates(map[string]interface{}{"status": "空闲", "power_bank_id": powerBankID, "order_no": ""})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("slot is not empty")
	}
	return nil
}

func (r *orderRepo) ReleaseSlot(ctx context.Context, slotID int64) error {
	return r.data.DB(ctx).Model(&slot{}).Where("id = ?", slotID).
		Updates(map[string]interface{}{"status": "已借出", "power_bank_id": 0}).Error
}

func (r *orderRepo) ReleaseSlotWithOrder(ctx context.Context, slotID int64, orderNo string) error {
	return r.data.DB(ctx).Model(&slot{}).Where("id = ?", slotID).
		Updates(map[string]interface{}{"status": "已借出", "power_bank_id": 0, "order_no": orderNo}).Error
}

// ----- Slots Count by Station -----

func (r *orderRepo) CountAvailableSlotsByStation(ctx context.Context, stationID int64) (int32, error) {
	var count int64
	if err := r.data.DB(ctx).Model(&powerBank{}).
		Where("station_id = ? AND status = '空闲'", stationID).Count(&count).Error; err != nil {
		return 0, err
	}
	return int32(count), nil
}

func (r *orderRepo) GetAvailableBanksByStation(ctx context.Context, stationID int64) (int32, error) {
	var count int64
	if err := r.data.DB(ctx).Model(&powerBank{}).
		Where("station_id = ? AND status = '空闲'", stationID).Count(&count).Error; err != nil {
		return 0, err
	}
	return int32(count), nil
}

// ----- Coupon -----

func (r *orderRepo) GetUserCoupons(ctx context.Context, userID int64) ([]*biz.Coupon, error) {
	var list []coupon
	if err := r.data.DB(ctx).
		Where("user_id = ? AND status = 'unused' AND expire_at > NOW()", userID).
		Find(&list).Error; err != nil {
		return nil, err
	}
	result := make([]*biz.Coupon, 0, len(list))
	for _, m := range list {
		result = append(result, r.toCouponBiz(&m))
	}
	return result, nil
}

func (r *orderRepo) GetCoupon(ctx context.Context, id int64) (*biz.Coupon, error) {
	var m coupon
	if err := r.data.DB(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return r.toCouponBiz(&m), nil
}

func (r *orderRepo) UseCoupon(ctx context.Context, id int64, userID int64) error {
	now := time.Now()
	result := r.data.DB(ctx).Model(&coupon{}).
		Where("id = ? AND user_id = ? AND status = 'unused'", id, userID).
		Updates(map[string]interface{}{"status": "used", "used_at": &now})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("coupon not available")
	}
	return nil
}

// ----- Converters -----

func (r *orderRepo) toBiz(m *order) *biz.Order {
	return &biz.Order{
		ID:              m.ID,
		OrderNo:         m.OrderNo,
		UserID:          m.UserID,
		PowerBankID:     m.PowerBankID,
		PowerBankNo:     m.PowerBankNo,
		BorrowStationID: m.BorrowStationID,
		BorrowCabinetID: m.BorrowCabinetID,
		BorrowSlotID:    m.BorrowSlotID,
		ReturnStationID: m.ReturnStationID,
		ReturnCabinetID: m.ReturnCabinetID,
		ReturnSlotID:    m.ReturnSlotID,
		BorrowTime:      m.BorrowTime,
		ReturnTime:      m.ReturnTime,
		DurationMinutes: m.DurationMinutes,
		StartFee:        m.StartFee,
		HourlyFee:       m.HourlyFee,
		DailyCap:        m.DailyCap,
		Deposit:         m.Deposit,
		TotalAmount:     m.TotalAmount,
		PaidAmount:      m.PaidAmount,
		DiscountAmount:  m.DiscountAmount,
		Status:          m.Status,
		Remark:          m.Remark,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func (r *orderRepo) toModel(o *biz.Order) *order {
	return &order{
		ID:              o.ID,
		OrderNo:         o.OrderNo,
		UserID:          o.UserID,
		PowerBankID:     o.PowerBankID,
		PowerBankNo:     o.PowerBankNo,
		BorrowStationID: o.BorrowStationID,
		BorrowCabinetID: o.BorrowCabinetID,
		BorrowSlotID:    o.BorrowSlotID,
		ReturnStationID: o.ReturnStationID,
		ReturnCabinetID: o.ReturnCabinetID,
		ReturnSlotID:    o.ReturnSlotID,
		BorrowTime:      o.BorrowTime,
		ReturnTime:      o.ReturnTime,
		DurationMinutes: o.DurationMinutes,
		StartFee:        o.StartFee,
		HourlyFee:       o.HourlyFee,
		DailyCap:        o.DailyCap,
		Deposit:         o.Deposit,
		TotalAmount:     o.TotalAmount,
		PaidAmount:      o.PaidAmount,
		DiscountAmount:  o.DiscountAmount,
		Status:          o.Status,
		Remark:          o.Remark,
	}
}

func (r *orderRepo) toPowerBankBiz(m *powerBank) *biz.PowerBank {
	return &biz.PowerBank{
		ID:           m.ID,
		DeviceNo:     m.DeviceNo,
		StationID:    m.StationID,
		CabinetID:    m.CabinetID,
		SlotID:       m.SlotID,
		Status:       m.Status,
		BatteryLevel: m.BatteryLevel,
	}
}

func (r *orderRepo) toSlotBiz(m *slot) *biz.Slot {
	return &biz.Slot{
		ID:          m.ID,
		CabinetID:   m.CabinetID,
		StationID:   m.StationID,
		SlotNo:      m.SlotNo,
		Status:      m.Status,
		PowerBankID: m.PowerBankID,
	}
}

// ----- Cabinet & Slot Details -----

func (r *orderRepo) ListCabinetsByStationID(ctx context.Context, stationID int64) ([]*biz.CabinetDetail, error) {
	var cabinets []cabinet
	if err := r.data.DB(ctx).Where("station_id = ?", stationID).Find(&cabinets).Error; err != nil {
		return nil, err
	}

	result := make([]*biz.CabinetDetail, 0, len(cabinets))
	for _, c := range cabinets {
		var slots []slot
		r.data.DB(ctx).Where("cabinet_id = ?", c.ID).Find(&slots)

		occupiedCount := 0
		slotDetails := make([]*biz.SlotDetail, 0, len(slots))
		for _, s := range slots {
			sd := &biz.SlotDetail{
				ID:     s.ID,
				SlotNo: s.SlotNo,
				Status: s.Status,
			}
			if s.Status == "空闲" && s.PowerBankID > 0 {
				occupiedCount++
				var pb powerBank
				if err := r.data.DB(ctx).First(&pb, s.PowerBankID).Error; err == nil {
					sd.PowerBankID = pb.ID
					sd.PowerBankNo = pb.DeviceNo
					sd.BatteryLevel = pb.BatteryLevel
					sd.Power = "22.5W"
				}
			}
			slotDetails = append(slotDetails, sd)
		}

		result = append(result, &biz.CabinetDetail{
			ID:             c.ID,
			CabinetNo:      c.CabinetNo,
			Status:         c.Status,
			TotalSlots:     int32(len(slots)),
			OccupiedSlots:  int32(occupiedCount),
			AvailableSlots: int32(len(slots) - occupiedCount),
			Slots:          slotDetails,
		})
	}

	return result, nil
}

func (r *orderRepo) GetOrderDetail(ctx context.Context, orderNo string) (*biz.Order, error) {
	var m order
	if err := r.data.DB(ctx).Where("order_no = ?", orderNo).First(&m).Error; err != nil {
		return nil, err
	}
	return r.toBiz(&m), nil
}

// ----- Scan Cabinet -----

func (r *orderRepo) ScanCabinet(ctx context.Context, cabinetID int64) (*biz.CabinetDetail, *biz.Station, error) {
	var c cabinet
	if err := r.data.DB(ctx).First(&c, cabinetID).Error; err != nil {
		return nil, nil, errors.New("机柜不存在")
	}
	if c.Status != 1 {
		return nil, nil, errors.New("机柜暂不可用")
	}

	// Get cabinet detail with slots
	cabinets, err := r.ListCabinetsByStationID(ctx, c.StationID)
	if err != nil {
		return nil, nil, err
	}
	var targetCabinet *biz.CabinetDetail
	for _, cb := range cabinets {
		if cb.ID == cabinetID {
			targetCabinet = cb
			break
		}
	}
	if targetCabinet == nil {
		return nil, nil, errors.New("机柜数据异常")
	}

	// Get station info
	var s station
	if err := r.data.DB(ctx).First(&s, c.StationID).Error; err != nil {
		return nil, nil, errors.New("站点不存在")
	}
	st := &biz.Station{
		ID:        s.ID,
		Name:      s.Name,
		Address:   s.Address,
		Latitude:  s.Latitude,
		Longitude: s.Longitude,
	}

	return targetCabinet, st, nil
}

// ----- Return Cabinets -----

func (r *orderRepo) ListReturnCabinets(ctx context.Context, lat, lng float64, radiusMeters int) ([]*biz.ReturnCabinet, error) {
	latDelta := float64(radiusMeters) / 111320.0
	lngDelta := float64(radiusMeters) / (111320.0 * cosDeg(lat))

	minLat := lat - latDelta
	maxLat := lat + latDelta
	minLng := lng - lngDelta
	maxLng := lng + lngDelta

	type cabinetRow struct {
		CabinetID      int64
		CabinetNo      string
		StationID      int64
		StationName    string
		StationAddress string
		Latitude       float64
		Longitude      float64
		TotalSlots     int32
		EmptyCount     int64
	}

	var rows []cabinetRow
	r.data.DB(ctx).Raw(`
		SELECT c.id AS cabinet_id, c.cabinet_no, s.id AS station_id, s.name AS station_name,
		       s.address AS station_address, s.latitude, s.longitude,
		       COUNT(sl.id) AS total_slots,
		       SUM(CASE WHEN sl.status = '已借出' THEN 1 ELSE 0 END) AS empty_count
		FROM cabinets c
		JOIN stations s ON s.id = c.station_id
		LEFT JOIN slots sl ON sl.cabinet_id = c.id
		WHERE c.status = 1 AND s.status = 1
		  AND s.latitude BETWEEN ? AND ?
		  AND s.longitude BETWEEN ? AND ?
		GROUP BY c.id, c.cabinet_no, s.id, s.name, s.address, s.latitude, s.longitude
		HAVING empty_count > 0
	`, minLat, maxLat, minLng, maxLng).Scan(&rows)

	type cabinetDist struct {
		cabinet  *biz.ReturnCabinet
		distance float64
	}
	var sorted []cabinetDist
	for _, row := range rows {
		dist := haversine(lat, lng, row.Latitude, row.Longitude)
		sorted = append(sorted, cabinetDist{
			cabinet: &biz.ReturnCabinet{
				CabinetID:      row.CabinetID,
				CabinetNo:      row.CabinetNo,
				StationID:      row.StationID,
				StationName:    row.StationName,
				StationAddress: row.StationAddress,
				Latitude:       row.Latitude,
				Longitude:      row.Longitude,
				Distance:       dist,
				EmptySlotCount: int32(row.EmptyCount),
				TotalSlots:     row.TotalSlots,
			},
			distance: dist,
		})
	}
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i].distance > sorted[j].distance {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	result := make([]*biz.ReturnCabinet, 0, len(sorted))
	for _, s := range sorted {
		result = append(result, s.cabinet)
	}
	return result, nil
}

func (r *orderRepo) FindCabinetByNo(ctx context.Context, cabinetNo string) (*biz.ReturnCabinet, error) {
	var c cabinet
	if err := r.data.DB(ctx).Where("cabinet_no = ? AND status = 1", cabinetNo).First(&c).Error; err != nil {
		return nil, err
	}
	var s station
	if err := r.data.DB(ctx).First(&s, c.StationID).Error; err != nil {
		return nil, err
	}
	// Count empty slots
	var totalSlots, emptyCount int64
	r.data.DB(ctx).Model(&slot{}).Where("cabinet_id = ?", c.ID).Count(&totalSlots)
	r.data.DB(ctx).Model(&slot{}).Where("cabinet_id = ? AND status = '已借出'", c.ID).Count(&emptyCount)
	if emptyCount == 0 {
		return nil, errors.New("该机柜暂无空闲仓位")
	}
	return &biz.ReturnCabinet{
		CabinetID:      c.ID,
		CabinetNo:      c.CabinetNo,
		StationID:      s.ID,
		StationName:    s.Name,
		StationAddress: s.Address,
		Latitude:       s.Latitude,
		Longitude:      s.Longitude,
		EmptySlotCount: int32(emptyCount),
		TotalSlots:     int32(totalSlots),
	}, nil
}

func (r *orderRepo) FindEmptySlot(ctx context.Context, cabinetID int64) (*biz.Slot, error) {
	var m slot
	if err := r.data.DB(ctx).
		Where("cabinet_id = ? AND status = '已借出'", cabinetID).
		Order("id ASC").First(&m).Error; err != nil {
		return nil, err
	}
	return r.toSlotBiz(&m), nil
}

func (r *orderRepo) UpdatePowerBankLocation(ctx context.Context, powerBankID, stationID, cabinetID, slotID int64) error {
	return r.data.DB(ctx).Model(&powerBank{}).Where("id = ?", powerBankID).
		Updates(map[string]interface{}{"station_id": stationID, "cabinet_id": cabinetID, "slot_id": slotID}).Error
}

func (r *orderRepo) toCouponBiz(m *coupon) *biz.Coupon {
	return &biz.Coupon{
		ID:        m.ID,
		UserID:    m.UserID,
		Type:      m.Type,
		Amount:    m.Amount,
		MinAmount: m.MinAmount,
		Status:    m.Status,
		ExpireAt:  m.ExpireAt,
	}
}
