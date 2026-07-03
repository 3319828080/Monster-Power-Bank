package biz

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/google/uuid"
)

var OrderProviderSet = wire.NewSet(NewOrderUsecase)

// ---------- Business Models ----------

type Order struct {
	ID                int64
	OrderNo           string
	UserID            int64
	PowerBankID       int64
	PowerBankNo       string
	BorrowStationID   int64
	BorrowStationName string
	BorrowCabinetID   int64
	BorrowSlotID      int64
	ReturnStationID   int64
	ReturnStationName string
	ReturnCabinetID   int64
	ReturnSlotID      int64
	BorrowTime        time.Time
	ReturnTime        *time.Time
	DurationMinutes   int32
	StartFee          int64
	HourlyFee         int64
	DailyCap          int64
	Deposit           int64
	TotalAmount       int64
	PaidAmount        int64
	DiscountAmount    int64
	Status            string
	Remark            string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type OrderItem struct {
	ID          int64
	OrderID     int64
	OrderNo     string
	FeeType     string
	Amount      int64
	Description string
	CreatedAt   time.Time
}

type PowerBank struct {
	ID           int64
	DeviceNo     string
	StationID    int64
	CabinetID    int64
	SlotID       int64
	Status       string
	BatteryLevel int32
}

type Station struct {
	ID             int64
	Name           string
	Address        string
	Latitude       float64
	Longitude      float64
	Distance       float64
	AvailableBanks int32
	OpenTime       string
	Description    string
	Images         string
}

type Slot struct {
	ID          int64
	CabinetID   int64
	StationID   int64
	SlotNo      string
	Status      string
	PowerBankID int64
}

type ReturnCabinet struct {
	CabinetID      int64
	CabinetNo      string
	StationID      int64
	StationName    string
	StationAddress string
	Latitude       float64
	Longitude      float64
	Distance       float64
	EmptySlotCount int32
	TotalSlots     int32
}

type Coupon struct {
	ID        int64
	UserID    int64
	Type      string
	Amount    int64
	MinAmount int64
	Status    string
	ExpireAt  time.Time
}

type CabinetDetail struct {
	ID             int64
	CabinetNo      string
	Status         int32
	TotalSlots     int32
	OccupiedSlots  int32
	AvailableSlots int32
	Slots          []*SlotDetail
}

type SlotDetail struct {
	ID           int64
	SlotNo       string
	Status       string
	PowerBankID  int64
	PowerBankNo  string
	BatteryLevel int32
	Power        string
}

// ---------- Repo Interfaces ----------

type Locker interface {
	Lock(ctx context.Context, key string, ttl time.Duration) (string, bool, error)
	Unlock(ctx context.Context, key, token string) error
}

type OrderRepo interface {
	Create(ctx context.Context, o *Order) (*Order, error)
	GetByID(ctx context.Context, id int64) (*Order, error)
	GetByOrderNo(ctx context.Context, orderNo string) (*Order, error)
	ListByUserID(ctx context.Context, userID int64, page, pageSize int, status string) ([]*Order, int, error)
	Update(ctx context.Context, o *Order) error
	UpdateStatus(ctx context.Context, id int64, status string) error
	CreateItem(ctx context.Context, item *OrderItem) error

	// Extended operations
	GetCurrentByUserID(ctx context.Context, userID int64) (*Order, error)
	UpdateOrderPartial(ctx context.Context, id int64, updates map[string]interface{}) error
	GetOrderByNoForUpdate(ctx context.Context, orderNo string) (*Order, error)
	GetPowerBank(ctx context.Context, id int64) (*PowerBank, error)
	UpdatePowerBankStatus(ctx context.Context, id int64, status string) error
	GetStation(ctx context.Context, id int64) (*Station, error)
	ListNearbyStations(ctx context.Context, lat, lng float64, radiusMeters int) ([]*Station, error)
	GetSlot(ctx context.Context, id int64) (*Slot, error)
	GetSlotByStationCabinet(ctx context.Context, stationID, cabinetID, slotID int64) (*Slot, error)
	OccupySlot(ctx context.Context, slotID, powerBankID int64) error
	ReleaseSlot(ctx context.Context, slotID int64) error
	ReleaseSlotWithOrder(ctx context.Context, slotID int64, orderNo string) error
	GetUserCoupons(ctx context.Context, userID int64) ([]*Coupon, error)
	GetCoupon(ctx context.Context, id int64) (*Coupon, error)
	UseCoupon(ctx context.Context, id int64, userID int64) error
	GetOrderDetail(ctx context.Context, orderNo string) (*Order, error)
	ListCabinetsByStationID(ctx context.Context, stationID int64) ([]*CabinetDetail, error)
	GetAvailableBanksByStation(ctx context.Context, stationID int64) (int32, error)
	ScanCabinet(ctx context.Context, cabinetID int64) (*CabinetDetail, *Station, error)
	CountAvailableSlotsByStation(ctx context.Context, stationID int64) (int32, error)
	ListReturnCabinets(ctx context.Context, lat, lng float64, radiusMeters int) ([]*ReturnCabinet, error)
	FindCabinetByNo(ctx context.Context, cabinetNo string) (*ReturnCabinet, error)
	FindEmptySlot(ctx context.Context, cabinetID int64) (*Slot, error)
	UpdatePowerBankLocation(ctx context.Context, powerBankID, stationID, cabinetID, slotID int64) error
}

// ---------- Pricing Config ----------

type PricingConfig struct {
	StartFee  int64 // 起步价（分）
	StartMins int   // 起步时长（分钟）
	HourlyFee int64 // 每小时费用（分）
	DailyCap  int64 // 每日封顶（分）
	Deposit   int64 // 默认押金（分）
}

// ---------- OrderUsecase ----------

type OrderUsecase struct {
	repo    OrderRepo
	wallet  WalletRepo
	locker  Locker
	Pricing *PricingConfig
	log     *log.Helper
}

func NewOrderUsecase(repo OrderRepo, wallet WalletRepo, locker Locker, pricing *PricingConfig, logger log.Logger) *OrderUsecase {
	return &OrderUsecase{
		repo:    repo,
		wallet:  wallet,
		locker:  locker,
		Pricing: pricing,
		log:     log.NewHelper(logger),
	}
}

// ---------- 1. PreCheckBorrow ----------

func (uc *OrderUsecase) PreCheckBorrow(ctx context.Context, userID, powerBankID, stationID, cabinetID, slotID int64) (bool, string, *PricingConfig, error) {
	// Check for existing active order
	existing, _ := uc.repo.GetCurrentByUserID(ctx, userID)
	if existing != nil {
		return false, "您有正在进行中的订单，请先归还", nil, nil
	}

	// Check slot status
	slot, err := uc.repo.GetSlotByStationCabinet(ctx, stationID, cabinetID, slotID)
	if err != nil {
		return false, "仓位信息不存在", nil, nil
	}
	if slot.Status != "空闲" {
		return false, "该仓位暂无充电宝", nil, nil
	}

	// Check power bank
	pb, err := uc.repo.GetPowerBank(ctx, powerBankID)
	if err != nil {
		return false, "充电宝信息不存在", nil, nil
	}
	if pb.Status != "空闲" {
		return false, "充电宝不可用", nil, nil
	}

	// Check user wallet / deposit
	wallet, err := uc.wallet.GetByUserID(ctx, userID)
	if err != nil {
		wallet, _ = uc.wallet.GetOrCreate(ctx, userID)
	}
	if wallet != nil && wallet.Balance+wallet.Frozen < uc.Pricing.Deposit {
		return false, fmt.Sprintf("余额不足，需要押金¥%.2f，当前余额¥%.2f", float64(uc.Pricing.Deposit)/100, float64(wallet.Balance+wallet.Frozen)/100), uc.Pricing, nil
	}

	return true, "", uc.Pricing, nil
}

// ---------- 2. CreateOrder ----------

func (uc *OrderUsecase) CreateOrder(ctx context.Context, userID, powerBankID, stationID, cabinetID, slotID int64) (*Order, error) {
	// Distributed lock
	token, ok, err := uc.locker.Lock(ctx, fmt.Sprintf("order:create:%d", userID), 10*time.Second)
	if err != nil {
		return nil, errors.New("系统繁忙，请稍后再试")
	}
	if !ok {
		return nil, errors.New("操作过于频繁，请稍后再试")
	}
	defer uc.locker.Unlock(ctx, fmt.Sprintf("order:create:%d", userID), token)

	// Check for existing active order
	existing, _ := uc.repo.GetCurrentByUserID(ctx, userID)
	if existing != nil {
		return nil, errors.New("您有正在进行中的订单，请先归还")
	}

	// Validate slot
	slot, err := uc.repo.GetSlotByStationCabinet(ctx, stationID, cabinetID, slotID)
	if err != nil {
		return nil, errors.New("仓位信息不存在")
	}
	if slot.Status != "空闲" || slot.PowerBankID != powerBankID {
		return nil, errors.New("充电宝不可用")
	}

	// Validate power bank
	pb, err := uc.repo.GetPowerBank(ctx, powerBankID)
	if err != nil {
		return nil, errors.New("充电宝信息不存在")
	}
	if pb.Status != "空闲" {
		return nil, errors.New("充电宝不可用")
	}

	// Get station name for reference
	station, _ := uc.repo.GetStation(ctx, stationID)
	stationName := ""
	if station != nil {
		stationName = station.Name
	}

	// Check deposit
	wallet, err := uc.wallet.GetOrCreate(ctx, userID)
	if err != nil {
		return nil, errors.New("获取钱包信息失败")
	}
	if wallet.Balance+wallet.Frozen < uc.Pricing.Deposit {
		return nil, errors.New("余额不足以冻结押金，请先充值")
	}

	// Create order (pending)
	now := time.Now()
	orderNo := generateOrderNo(userID)
	order := &Order{
		OrderNo:           orderNo,
		UserID:            userID,
		PowerBankID:       powerBankID,
		PowerBankNo:       pb.DeviceNo,
		BorrowStationID:   stationID,
		BorrowStationName: stationName,
		BorrowCabinetID:   cabinetID,
		BorrowSlotID:      slotID,
		BorrowTime:        now,
		StartFee:          uc.Pricing.StartFee,
		HourlyFee:         uc.Pricing.HourlyFee,
		DailyCap:          uc.Pricing.DailyCap,
		Deposit:           uc.Pricing.Deposit,
		Status:            "租借中",
	}

	// Execute in transaction
	txErr := uc.wallet.Transaction(ctx, func(ctx context.Context) error {
		created, err := uc.repo.Create(ctx, order)
		if err != nil {
			return err
		}
		order.ID = created.ID

		// Freeze deposit
		if err := uc.wallet.FreezeBalance(ctx, wallet.ID, uc.Pricing.Deposit); err != nil {
			return errors.New("冻结押金失败")
		}

		// Release slot with order binding
		if err := uc.repo.ReleaseSlotWithOrder(ctx, slotID, orderNo); err != nil {
			return err
		}

		// Update power bank status and clear location (physically taken away)
		if err := uc.repo.UpdatePowerBankStatus(ctx, powerBankID, "租借中"); err != nil {
			return err
		}
		if err := uc.repo.UpdatePowerBankLocation(ctx, powerBankID, 0, 0, 0); err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	_ = uc.repo.CreateItem(ctx, &OrderItem{
		OrderID:     order.ID,
		OrderNo:     orderNo,
		FeeType:     "deposit_freeze",
		Amount:      -uc.Pricing.Deposit,
		Description: fmt.Sprintf("借用充电宝冻结押金（%s）", station.Name),
	})

	return order, nil
}

// ---------- 3. CancelOrder ----------

func (uc *OrderUsecase) CancelOrder(ctx context.Context, userID int64, orderNo string) error {
	order, err := uc.repo.GetByOrderNo(ctx, orderNo)
	if err != nil {
		return errors.New("订单不存在")
	}
	if order.UserID != userID {
		return errors.New("无权操作此订单")
	}
	if order.Status != "pending" && order.Status != "租借中" {
		return errors.New("当前状态不允许取消")
	}
	// For borrowed orders, we need to return the power bank first
	if order.Status == "租借中" {
		return errors.New("充电宝已借出，请通过归还流程操作")
	}

	return uc.repo.UpdateStatus(ctx, order.ID, "cancelled")
}

// ---------- 4. ReturnPowerBank ----------

func (uc *OrderUsecase) ReturnPowerBank(ctx context.Context, userID int64, orderNo string, stationID, cabinetID, slotID int64, couponID int64) (*Order, error) {
	var result *Order
	err := uc.wallet.Transaction(ctx, func(ctx context.Context) error {
		order, err := uc.repo.GetByOrderNo(ctx, orderNo)
		if err != nil {
			return errors.New("订单不存在")
		}
		if order.UserID != userID {
			return errors.New("无权操作此订单")
		}
		if order.Status != "租借中" {
			return errors.New("当前订单状态不允许归还")
		}

		// Auto-assign or validate return slot
		var slot *Slot
		if stationID == 0 && slotID == 0 {
			slot, err = uc.repo.FindEmptySlot(ctx, cabinetID)
			if err != nil {
				return errors.New("该机柜暂无空闲仓位")
			}
			stationID = slot.StationID
			slotID = slot.ID
		} else {
			slot, err = uc.repo.GetSlotByStationCabinet(ctx, stationID, cabinetID, slotID)
			if err != nil {
				return errors.New("归还仓位不存在")
			}
			if slot.Status != "已借出" {
				return errors.New("该仓位已被占用，请选择其他仓位")
			}
		}

		// Calculate fee
		now := time.Now()
		duration := now.Sub(order.BorrowTime)
		totalMinutes := int32(duration.Minutes())
		if totalMinutes < 1 {
			totalMinutes = 1
		}
		totalAmount := calculateFee(duration, order.StartFee, uc.Pricing.StartMins, order.HourlyFee, order.DailyCap)
		discountAmount := int64(0)

		// Apply coupon
		if couponID > 0 {
			coupon, err := uc.repo.GetCoupon(ctx, couponID)
			if err == nil && coupon.UserID == userID && coupon.Status == "unused" {
				if totalAmount >= coupon.MinAmount {
					switch coupon.Type {
					case "full_reduce":
						discountAmount = coupon.Amount
					case "percent":
						discountAmount = totalAmount * coupon.Amount / 100
					case "fixed":
						discountAmount = coupon.Amount
					}
					if discountAmount > totalAmount {
						discountAmount = totalAmount
					}
					uc.repo.UseCoupon(ctx, couponID, userID)
				}
			}
		}

		paidAmount := totalAmount - discountAmount
		if paidAmount < 0 {
			paidAmount = 0
		}

		// Deduct: balance first, then deposit
		wallet, err := uc.wallet.GetOrCreate(ctx, userID)
		if err != nil {
			return errors.New("获取钱包信息失败")
		}
		depositDeducted := int64(0)
		remaining := paidAmount

		// Deduct from balance
		if remaining > 0 {
			deduct := remaining
			if wallet.Balance < deduct {
				deduct = wallet.Balance
			}
			_ = deduct
			newBal := wallet.Balance - deduct
			if err := uc.wallet.UpdateBalance(ctx, wallet.ID, newBal, wallet.Frozen); err != nil {
				return errors.New("扣费失败")
			}
			remaining -= deduct

			_ = uc.repo.CreateItem(ctx, &OrderItem{
				OrderID:     order.ID,
				OrderNo:     orderNo,
				FeeType:     "rental_fee",
				Amount:      -deduct,
				Description: "充电宝租借费用（余额支付）",
			})
		}

		// Deduct remaining from deposit (unfreeze + deduct)
		if remaining > 0 && order.Deposit > 0 {
			depositToRelease := order.Deposit
			if remaining >= depositToRelease {
				depositDeducted = depositToRelease
				depositToRelease = 0
			} else {
				depositDeducted = remaining
				depositToRelease = depositToRelease - remaining
			}
			remaining -= depositDeducted

			// Unfreeze remaining deposit, deduct owed
			if depositToRelease > 0 {
				if err := uc.wallet.UnfreezeBalance(ctx, wallet.ID, order.Deposit); err != nil {
					return errors.New("退还押金失败")
				}
				if depositDeducted > 0 {
					if _, err := uc.wallet.AddBalance(ctx, wallet.ID, depositToRelease); err != nil {
						return errors.New("退还押金失败")
					}
				}
			} else {
				// All deposit consumed
				if err := uc.wallet.UnfreezeBalance(ctx, wallet.ID, order.Deposit); err != nil {
					return errors.New("释放押金失败")
				}
			}

			_ = uc.repo.CreateItem(ctx, &OrderItem{
				OrderID:     order.ID,
				OrderNo:     orderNo,
				FeeType:     "deposit_deduct",
				Amount:      -depositDeducted,
				Description: "押金扣除租借费用",
			})
		} else if order.Deposit > 0 {
			// No remaining fee, fully unfreeze deposit
			if err := uc.wallet.UnfreezeBalance(ctx, wallet.ID, order.Deposit); err != nil {
				return errors.New("退还押金失败")
			}
			_ = uc.repo.CreateItem(ctx, &OrderItem{
				OrderID:     order.ID,
				OrderNo:     orderNo,
				FeeType:     "deposit_refund",
				Amount:      order.Deposit,
				Description: "退还押金",
			})
		}

		// Update order
		updates := map[string]interface{}{
			"return_station_id": stationID,
			"return_cabinet_id": cabinetID,
			"return_slot_id":    slotID,
			"return_time":       &now,
			"duration_minutes":  totalMinutes,
			"total_amount":      totalAmount,
			"paid_amount":       paidAmount,
			"discount_amount":   discountAmount,
			"status":            "completed",
		}
		if err := uc.repo.UpdateOrderPartial(ctx, order.ID, updates); err != nil {
			return errors.New("更新订单失败")
		}

		// Occupy return slot
		if err := uc.repo.OccupySlot(ctx, slotID, order.PowerBankID); err != nil {
			return err
		}

		// Update power bank status and location
		if err := uc.repo.UpdatePowerBankStatus(ctx, order.PowerBankID, "空闲"); err != nil {
			return err
		}
		if err := uc.repo.UpdatePowerBankLocation(ctx, order.PowerBankID, stationID, cabinetID, slotID); err != nil {
			return err
		}

		order.ReturnStationID = stationID
		order.ReturnCabinetID = cabinetID
		order.ReturnSlotID = slotID
		order.ReturnTime = &now
		order.DurationMinutes = totalMinutes
		order.TotalAmount = totalAmount
		order.PaidAmount = paidAmount
		order.DiscountAmount = discountAmount
		order.Status = "completed"

		result = order
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ---------- 5. ExtendRental ----------

func (uc *OrderUsecase) ExtendRental(ctx context.Context, userID int64, orderNo string, extraMinutes int32) error {
	return uc.wallet.Transaction(ctx, func(ctx context.Context) error {
		order, err := uc.repo.GetByOrderNo(ctx, orderNo)
		if err != nil {
			return errors.New("订单不存在")
		}
		if order.UserID != userID {
			return errors.New("无权操作此订单")
		}
		if order.Status != "租借中" {
			return errors.New("当前订单状态不允许续租")
		}
		if extraMinutes <= 0 {
			return errors.New("延长时长无效")
		}

		// Calculate additional deposit needed
		extraHours := (extraMinutes + 59) / 60
		additionalFee := int64(extraHours) * order.HourlyFee
		additionalDeposit := additionalFee - order.Deposit
		if additionalDeposit < 0 {
			additionalDeposit = 0
		}

		if additionalDeposit > 0 {
			wallet, err := uc.wallet.GetOrCreate(ctx, userID)
			if err != nil {
				return errors.New("获取钱包信息失败")
			}
			if wallet.Balance < additionalDeposit {
				return errors.New("余额不足以续租，请先充值")
			}
			if err := uc.wallet.FreezeBalance(ctx, wallet.ID, additionalDeposit); err != nil {
				return errors.New("冻结押金失败")
			}
		}

		// Update order deposit
		if err := uc.repo.UpdateOrderPartial(ctx, order.ID, map[string]interface{}{
			"deposit": order.Deposit + additionalDeposit,
			"remark":  fmt.Sprintf("续租%d分钟", extraMinutes),
		}); err != nil {
			return err
		}

		return nil
	})
}

// ---------- 6. GetCurrentOrder ----------

func (uc *OrderUsecase) GetCurrentOrder(ctx context.Context, userID int64) (*Order, error) {
	order, err := uc.repo.GetCurrentByUserID(ctx, userID)
	if err != nil {
		return nil, nil // no active order
	}

	// Enrich with station name
	borrowStation, _ := uc.repo.GetStation(ctx, order.BorrowStationID)
	if borrowStation != nil {
		order.BorrowStationName = borrowStation.Name
	}

	return order, nil
}

// ---------- 7. GetRealTimeFee ----------

func (uc *OrderUsecase) GetRealTimeFee(ctx context.Context, userID int64, orderNo string) (int64, int64, int64, int64, error) {
	order, err := uc.repo.GetByOrderNo(ctx, orderNo)
	if err != nil {
		return 0, 0, 0, 0, errors.New("订单不存在")
	}
	if order.UserID != userID {
		return 0, 0, 0, 0, errors.New("无权查看此订单")
	}
	if order.Status != "租借中" {
		return 0, 0, 0, 0, errors.New("订单已结束")
	}

	elapsed := time.Since(order.BorrowTime)
	elapsedSeconds := int64(elapsed.Seconds())
	currentFee := calculateFee(elapsed, order.StartFee, uc.Pricing.StartMins, order.HourlyFee, order.DailyCap)

	return elapsedSeconds, currentFee, order.HourlyFee, order.DailyCap, nil
}

// ---------- 8. ListOrders ----------

func (uc *OrderUsecase) ListOrders(ctx context.Context, userID int64, page, pageSize int, status string) ([]*Order, int, error) {
	list, total, err := uc.repo.ListByUserID(ctx, userID, page, pageSize, status)
	if err != nil {
		return nil, 0, err
	}
	// Enrich with station names
	for _, o := range list {
		if o.BorrowStationName == "" && o.BorrowStationID > 0 {
			if s, e := uc.repo.GetStation(ctx, o.BorrowStationID); e == nil {
				o.BorrowStationName = s.Name
			}
		}
		if o.ReturnStationName == "" && o.ReturnStationID > 0 {
			if s, e := uc.repo.GetStation(ctx, o.ReturnStationID); e == nil {
				o.ReturnStationName = s.Name
			}
		}
	}
	return list, total, nil
}

// ---------- 9. CreateCompensationPayment ----------

func (uc *OrderUsecase) CreateCompensationPayment(ctx context.Context, userID int64, orderNo string) (string, int64, error) {
	order, err := uc.repo.GetByOrderNo(ctx, orderNo)
	if err != nil {
		return "", 0, errors.New("订单不存在")
	}
	if order.UserID != userID {
		return "", 0, errors.New("无权操作此订单")
	}
	if order.Status != "租借中" {
		return "", 0, errors.New("当前状态不可赔付")
	}

	// Calculate overdue + compensation: use hourly fee as base
	overdueMinutes := int32(time.Since(order.BorrowTime).Minutes())
	compAmount := int64(overdueMinutes/60+1) * order.HourlyFee
	if compAmount < order.Deposit {
		compAmount = order.Deposit
	}

	// Update status to pending_compensation
	if err := uc.repo.UpdateStatus(ctx, order.ID, "pending_compensation"); err != nil {
		return "", 0, errors.New("更新订单失败")
	}

	// Return payment info (channel integration later)
	return orderNo, compAmount, nil
}

// ---------- 10. GetOrderDetail ----------

func (uc *OrderUsecase) GetOrderDetail(ctx context.Context, userID int64, orderNo string) (*Order, error) {
	order, err := uc.repo.GetOrderDetail(ctx, orderNo)
	if err != nil {
		return nil, errors.New("订单不存在")
	}
	if order.UserID != userID {
		return nil, errors.New("无权查看此订单")
	}
	return order, nil
}

// ---------- 11. GetCabinetList ----------

func (uc *OrderUsecase) GetCabinetList(ctx context.Context, stationID int64) ([]*CabinetDetail, error) {
	return uc.repo.ListCabinetsByStationID(ctx, stationID)
}

// ---------- 12. GetDelayFee ----------

func (uc *OrderUsecase) GetDelayFee(ctx context.Context, userID int64, orderNo string, extraMinutes int32) (int64, int64, int64, error) {
	order, err := uc.repo.GetByOrderNo(ctx, orderNo)
	if err != nil {
		return 0, 0, 0, errors.New("订单不存在")
	}
	if order.UserID != userID {
		return 0, 0, 0, errors.New("无权操作此订单")
	}
	if order.Status != "租借中" {
		return 0, 0, 0, errors.New("当前订单状态不允许续租")
	}
	if extraMinutes <= 0 {
		return 0, 0, 0, errors.New("延长时长无效")
	}

	// Calculate additional fee for the extended period
	extraHours := int32(extraMinutes+59) / 60
	additionalFee := int64(extraHours) * order.HourlyFee

	// Estimate total: current elapsed fee + extension fee
	elapsed := time.Since(order.BorrowTime)
	currentFee := calculateFee(elapsed, order.StartFee, uc.Pricing.StartMins, order.HourlyFee, order.DailyCap)
	estimatedTotal := currentFee + additionalFee
	if order.DailyCap > 0 && estimatedTotal > order.DailyCap {
		estimatedTotal = order.DailyCap
		additionalFee = estimatedTotal - currentFee
		if additionalFee < 0 {
			additionalFee = 0
		}
	}

	// Calculate additional deposit needed
	additionalDeposit := additionalFee - order.Deposit
	if additionalDeposit < 0 {
		additionalDeposit = 0
	}

	return additionalFee, estimatedTotal, additionalDeposit, nil
}

// ---------- Station Map ----------

func (uc *OrderUsecase) GetStationDetail(ctx context.Context, id int64) (*Station, error) {
	station, err := uc.repo.GetStation(ctx, id)
	if err != nil {
		return nil, err
	}
	// Compute available banks count for this station
	if cabinets, err := uc.repo.ListCabinetsByStationID(ctx, id); err == nil {
		var avail int32
		for _, c := range cabinets {
			avail += c.OccupiedSlots
		}
		station.AvailableBanks = avail
	}
	return station, nil
}

func (uc *OrderUsecase) ListNearbyStations(ctx context.Context, lat, lng float64, radiusMeters int) ([]*Station, error) {
	if radiusMeters <= 0 || radiusMeters > 50000 {
		radiusMeters = 5000
	}
	stations, err := uc.repo.ListNearbyStations(ctx, lat, lng, radiusMeters)
	if err != nil {
		return nil, err
	}
	// Populate available banks for each station
	for _, s := range stations {
		count, _ := uc.repo.CountAvailableSlotsByStation(ctx, s.ID)
		s.AvailableBanks = count
	}
	return stations, nil
}

// ---------- 13. ScanCabinet ----------

func (uc *OrderUsecase) ScanCabinet(ctx context.Context, cabinetID int64) (*CabinetDetail, *Station, error) {
	cabinet, station, err := uc.repo.ScanCabinet(ctx, cabinetID)
	if err != nil {
		return nil, nil, err
	}
	return cabinet, station, nil
}

// ---------- 14. ListReturnCabinets ----------

func (uc *OrderUsecase) ListReturnCabinets(ctx context.Context, lat, lng float64, radiusMeters int) ([]*ReturnCabinet, error) {
	if radiusMeters <= 0 || radiusMeters > 50000 {
		radiusMeters = 5000
	}
	return uc.repo.ListReturnCabinets(ctx, lat, lng, radiusMeters)
}

// ---------- 15. SearchReturnCabinet ----------

func (uc *OrderUsecase) SearchReturnCabinet(ctx context.Context, cabinetNo string) (*ReturnCabinet, error) {
	return uc.repo.FindCabinetByNo(ctx, cabinetNo)
}

// ---------- Helpers ----------

func generateOrderNo(userID int64) string {
	now := time.Now()
	return fmt.Sprintf("MP%s%05d%s", now.Format("20060102150405"), userID%100000, uuid.New().String()[:8])
}

func calculateFee(duration time.Duration, startFee int64, startMins int, hourlyFee int64, dailyCap int64) int64 {
	totalMinutes := duration.Minutes()
	if totalMinutes <= float64(startMins) {
		return startFee
	}

	extraMinutes := totalMinutes - float64(startMins)
	extraHours := int64(extraMinutes+59) / 60 // ceil division

	total := startFee + extraHours*hourlyFee

	if dailyCap > 0 && total > dailyCap {
		total = dailyCap
	}
	return total
}
