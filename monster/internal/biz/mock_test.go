package biz

import (
	"context"
	"errors"
	"monster/pkg/payment"
	"time"
)

// ---------- mockOrderRepo ----------

type mockOrderRepo struct {
	createFn              func(ctx context.Context, o *Order) (*Order, error)
	getByOrderNoFn        func(ctx context.Context, orderNo string) (*Order, error)
	getCurrentByUserIDFn  func(ctx context.Context, userID int64) (*Order, error)
	getSlotByStationCabinetFn func(ctx context.Context, stationID, cabinetID, slotID int64) (*Slot, error)
	getPowerBankFn        func(ctx context.Context, id int64) (*PowerBank, error)
	getStationFn          func(ctx context.Context, id int64) (*Station, error)
	updateStatusFn        func(ctx context.Context, id int64, status string) error
	updateOrderPartialFn  func(ctx context.Context, id int64, updates map[string]any) error
	createItemFn          func(ctx context.Context, item *OrderItem) error
	updatePowerBankStatusFn func(ctx context.Context, id int64, status string) error
	releaseSlotWithOrderFn func(ctx context.Context, slotID int64, orderNo string) error
	updatePowerBankLocationFn func(ctx context.Context, powerBankID, stationID, cabinetID, slotID int64) error
	getCouponFn           func(ctx context.Context, id int64) (*Coupon, error)
	useCouponFn           func(ctx context.Context, id int64, userID int64) error
	occupySlotFn          func(ctx context.Context, slotID, powerBankID int64) error
	findEmptySlotFn       func(ctx context.Context, cabinetID int64) (*Slot, error)
	getOrderDetailFn      func(ctx context.Context, orderNo string) (*Order, error)
	listByUserIDFn        func(ctx context.Context, userID int64, page, pageSize int, status string) ([]*Order, int, error)
}

func (m *mockOrderRepo) Create(ctx context.Context, o *Order) (*Order, error) {
	if m.createFn != nil {
		return m.createFn(ctx, o)
	}
	clone := *o
	clone.ID = 1
	return &clone, nil
}
func (m *mockOrderRepo) GetByID(ctx context.Context, id int64) (*Order, error) { return nil, errors.New("not implemented") }
func (m *mockOrderRepo) GetByOrderNo(ctx context.Context, orderNo string) (*Order, error) {
	if m.getByOrderNoFn != nil {
		return m.getByOrderNoFn(ctx, orderNo)
	}
	return nil, errors.New("not implemented")
}
func (m *mockOrderRepo) ListByUserID(ctx context.Context, userID int64, page, pageSize int, status string) ([]*Order, int, error) {
	if m.listByUserIDFn != nil {
		return m.listByUserIDFn(ctx, userID, page, pageSize, status)
	}
	return nil, 0, errors.New("not implemented")
}
func (m *mockOrderRepo) Update(ctx context.Context, o *Order) error { return errors.New("not implemented") }
func (m *mockOrderRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	if m.updateStatusFn != nil {
		return m.updateStatusFn(ctx, id, status)
	}
	return nil
}
func (m *mockOrderRepo) CreateItem(ctx context.Context, item *OrderItem) error {
	if m.createItemFn != nil {
		return m.createItemFn(ctx, item)
	}
	return nil
}
func (m *mockOrderRepo) GetCurrentByUserID(ctx context.Context, userID int64) (*Order, error) {
	if m.getCurrentByUserIDFn != nil {
		return m.getCurrentByUserIDFn(ctx, userID)
	}
	return nil, nil
}
func (m *mockOrderRepo) UpdateOrderPartial(ctx context.Context, id int64, updates map[string]any) error {
	if m.updateOrderPartialFn != nil {
		return m.updateOrderPartialFn(ctx, id, updates)
	}
	return nil
}
func (m *mockOrderRepo) GetOrderByNoForUpdate(ctx context.Context, orderNo string) (*Order, error) { return nil, errors.New("not implemented") }
func (m *mockOrderRepo) GetPowerBank(ctx context.Context, id int64) (*PowerBank, error) {
	if m.getPowerBankFn != nil {
		return m.getPowerBankFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}
func (m *mockOrderRepo) UpdatePowerBankStatus(ctx context.Context, id int64, status string) error {
	if m.updatePowerBankStatusFn != nil {
		return m.updatePowerBankStatusFn(ctx, id, status)
	}
	return nil
}
func (m *mockOrderRepo) GetStation(ctx context.Context, id int64) (*Station, error) {
	if m.getStationFn != nil {
		return m.getStationFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}
func (m *mockOrderRepo) ListNearbyStations(ctx context.Context, lat, lng float64, radiusMeters int) ([]*Station, error) { return nil, errors.New("not implemented") }
func (m *mockOrderRepo) GetSlot(ctx context.Context, id int64) (*Slot, error) { return nil, errors.New("not implemented") }
func (m *mockOrderRepo) GetSlotByStationCabinet(ctx context.Context, stationID, cabinetID, slotID int64) (*Slot, error) {
	if m.getSlotByStationCabinetFn != nil {
		return m.getSlotByStationCabinetFn(ctx, stationID, cabinetID, slotID)
	}
	return nil, errors.New("not implemented")
}
func (m *mockOrderRepo) OccupySlot(ctx context.Context, slotID, powerBankID int64) error {
	if m.occupySlotFn != nil {
		return m.occupySlotFn(ctx, slotID, powerBankID)
	}
	return nil
}
func (m *mockOrderRepo) ReleaseSlot(ctx context.Context, slotID int64) error { return errors.New("not implemented") }
func (m *mockOrderRepo) ReleaseSlotWithOrder(ctx context.Context, slotID int64, orderNo string) error {
	if m.releaseSlotWithOrderFn != nil {
		return m.releaseSlotWithOrderFn(ctx, slotID, orderNo)
	}
	return nil
}
func (m *mockOrderRepo) GetUserCoupons(ctx context.Context, userID int64) ([]*Coupon, error) { return nil, errors.New("not implemented") }
func (m *mockOrderRepo) GetCoupon(ctx context.Context, id int64) (*Coupon, error) {
	if m.getCouponFn != nil {
		return m.getCouponFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}
func (m *mockOrderRepo) UseCoupon(ctx context.Context, id int64, userID int64) error {
	if m.useCouponFn != nil {
		return m.useCouponFn(ctx, id, userID)
	}
	return errors.New("not implemented")
}
func (m *mockOrderRepo) GetOrderDetail(ctx context.Context, orderNo string) (*Order, error) {
	if m.getOrderDetailFn != nil {
		return m.getOrderDetailFn(ctx, orderNo)
	}
	return nil, errors.New("not implemented")
}
func (m *mockOrderRepo) ListCabinetsByStationID(ctx context.Context, stationID int64) ([]*CabinetDetail, error) { return nil, errors.New("not implemented") }
func (m *mockOrderRepo) GetAvailableBanksByStation(ctx context.Context, stationID int64) (int32, error) { return 0, errors.New("not implemented") }
func (m *mockOrderRepo) ScanCabinet(ctx context.Context, cabinetID int64) (*CabinetDetail, *Station, error) { return nil, nil, errors.New("not implemented") }
func (m *mockOrderRepo) CountAvailableSlotsByStation(ctx context.Context, stationID int64) (int32, error) { return 0, errors.New("not implemented") }
func (m *mockOrderRepo) ListReturnCabinets(ctx context.Context, lat, lng float64, radiusMeters int) ([]*ReturnCabinet, error) { return nil, errors.New("not implemented") }
func (m *mockOrderRepo) FindCabinetByNo(ctx context.Context, cabinetNo string) (*ReturnCabinet, error) { return nil, errors.New("not implemented") }
func (m *mockOrderRepo) FindEmptySlot(ctx context.Context, cabinetID int64) (*Slot, error) {
	if m.findEmptySlotFn != nil {
		return m.findEmptySlotFn(ctx, cabinetID)
	}
	return nil, errors.New("not implemented")
}
func (m *mockOrderRepo) UpdatePowerBankLocation(ctx context.Context, powerBankID, stationID, cabinetID, slotID int64) error {
	if m.updatePowerBankLocationFn != nil {
		return m.updatePowerBankLocationFn(ctx, powerBankID, stationID, cabinetID, slotID)
	}
	return nil
}

// ---------- mockWalletRepo ----------

type mockWalletRepo struct {
	getByUserIDFn     func(ctx context.Context, userID int64) (*Wallet, error)
	getOrCreateFn     func(ctx context.Context, userID int64) (*Wallet, error)
	freezeBalanceFn   func(ctx context.Context, id int64, amount int64) error
	addBalanceFn      func(ctx context.Context, id int64, amount int64) (int64, error)
	unfreezeBalanceFn func(ctx context.Context, id int64, amount int64) error
	updateBalanceFn   func(ctx context.Context, id int64, balance, frozen int64) error
	transactionFn     func(ctx context.Context, fn func(ctx context.Context) error) error
	createTransactionFn func(ctx context.Context, t *WalletTransaction) error
}

func (m *mockWalletRepo) GetByUserID(ctx context.Context, userID int64) (*Wallet, error) {
	if m.getByUserIDFn != nil {
		return m.getByUserIDFn(ctx, userID)
	}
	return nil, errors.New("not implemented")
}
func (m *mockWalletRepo) Create(ctx context.Context, w *Wallet) (*Wallet, error) { return nil, errors.New("not implemented") }
func (m *mockWalletRepo) UpdateBalance(ctx context.Context, id int64, balance, frozen int64) error {
	if m.updateBalanceFn != nil {
		return m.updateBalanceFn(ctx, id, balance, frozen)
	}
	return nil
}
func (m *mockWalletRepo) CreateTransaction(ctx context.Context, t *WalletTransaction) error {
	if m.createTransactionFn != nil {
		return m.createTransactionFn(ctx, t)
	}
	return nil
}
func (m *mockWalletRepo) ListTransactions(ctx context.Context, userID int64, page, pageSize int) ([]*WalletTransaction, int, error) {
	return nil, 0, errors.New("not implemented")
}
func (m *mockWalletRepo) AddBalance(ctx context.Context, id int64, amount int64) (int64, error) {
	if m.addBalanceFn != nil {
		return m.addBalanceFn(ctx, id, amount)
	}
	return 0, errors.New("not implemented")
}
func (m *mockWalletRepo) AddFrozen(ctx context.Context, id int64, amount int64) (int64, error) { return 0, errors.New("not implemented") }
func (m *mockWalletRepo) FreezeBalance(ctx context.Context, id int64, amount int64) error {
	if m.freezeBalanceFn != nil {
		return m.freezeBalanceFn(ctx, id, amount)
	}
	return nil
}
func (m *mockWalletRepo) UnfreezeBalance(ctx context.Context, id int64, amount int64) error {
	if m.unfreezeBalanceFn != nil {
		return m.unfreezeBalanceFn(ctx, id, amount)
	}
	return nil
}
func (m *mockWalletRepo) GetOrCreate(ctx context.Context, userID int64) (*Wallet, error) {
	if m.getOrCreateFn != nil {
		return m.getOrCreateFn(ctx, userID)
	}
	// Default: return wallet with sufficient balance
	return &Wallet{ID: 1, UserID: userID, Balance: 100000, Frozen: 0}, nil
}
func (m *mockWalletRepo) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if m.transactionFn != nil {
		return m.transactionFn(ctx, fn)
	}
	return fn(ctx)
}

// ---------- mockLocker ----------

type mockLocker struct {
	lockFn   func(ctx context.Context, key string, ttl time.Duration) (string, bool, error)
	unlockFn func(ctx context.Context, key, token string) error
}

func (m *mockLocker) Lock(ctx context.Context, key string, ttl time.Duration) (string, bool, error) {
	if m.lockFn != nil {
		return m.lockFn(ctx, key, ttl)
	}
	return "token-xxx", true, nil
}
func (m *mockLocker) Unlock(ctx context.Context, key, token string) error {
	if m.unlockFn != nil {
		return m.unlockFn(ctx, key, token)
	}
	return nil
}

// ---------- mockPaymentRepo ----------

type mockPaymentRepo struct {
	createFn               func(ctx context.Context, p *Payment) (*Payment, error)
	getByPaymentNoFn       func(ctx context.Context, paymentNo string) (*Payment, error)
	updateStatusFn         func(ctx context.Context, id int64, status, transactionNo string, paidAt *time.Time) error
	createNotificationFn   func(ctx context.Context, n *Notification) error
}

func (m *mockPaymentRepo) Create(ctx context.Context, p *Payment) (*Payment, error) {
	if m.createFn != nil {
		return m.createFn(ctx, p)
	}
	clone := *p
	clone.ID = 1
	return &clone, nil
}
func (m *mockPaymentRepo) GetByPaymentNo(ctx context.Context, paymentNo string) (*Payment, error) {
	if m.getByPaymentNoFn != nil {
		return m.getByPaymentNoFn(ctx, paymentNo)
	}
	return nil, errors.New("not implemented")
}
func (m *mockPaymentRepo) GetByOrderNo(ctx context.Context, orderNo string) (*Payment, error) { return nil, errors.New("not implemented") }
func (m *mockPaymentRepo) UpdateStatus(ctx context.Context, id int64, status, transactionNo string, paidAt *time.Time) error {
	if m.updateStatusFn != nil {
		return m.updateStatusFn(ctx, id, status, transactionNo, paidAt)
	}
	return nil
}
func (m *mockPaymentRepo) CreateRefund(ctx context.Context, r *Refund) (*Refund, error) { return nil, errors.New("not implemented") }
func (m *mockPaymentRepo) GetRefundByRefundNo(ctx context.Context, refundNo string) (*Refund, error) { return nil, errors.New("not implemented") }
func (m *mockPaymentRepo) GetRefundByOrderNo(ctx context.Context, orderNo string) (*Refund, error) { return nil, errors.New("not implemented") }
func (m *mockPaymentRepo) UpdateRefundStatus(ctx context.Context, id int64, status, transactionNo string, refundedAt *time.Time) error { return errors.New("not implemented") }
func (m *mockPaymentRepo) CreateNotification(ctx context.Context, n *Notification) error {
	if m.createNotificationFn != nil {
		return m.createNotificationFn(ctx, n)
	}
	return nil
}
func (m *mockPaymentRepo) MarkNotificationsRead(ctx context.Context, ids []int64, userID int64) error { return nil }

// ---------- mockPaymentChannel ----------

type mockPaymentChannel struct {
	name             string
	payFn            func(ctx context.Context, req *payment.PaymentRequest) (*payment.PaymentResponse, error)
	verifyNotifyFn   func(ctx context.Context, params map[string]string) (*payment.PaymentResponse, error)
}

func (m *mockPaymentChannel) Name() string { return m.name }
func (m *mockPaymentChannel) Pay(ctx context.Context, req *payment.PaymentRequest) (*payment.PaymentResponse, error) {
	if m.payFn != nil {
		return m.payFn(ctx, req)
	}
	return &payment.PaymentResponse{PayInfo: map[string]string{"pay_url": "http://pay.example.com"}}, nil
}
func (m *mockPaymentChannel) Refund(ctx context.Context, req *payment.RefundRequest) (*payment.RefundResponse, error) {
	return nil, errors.New("not implemented")
}
func (m *mockPaymentChannel) VerifyNotify(ctx context.Context, params map[string]string) (*payment.PaymentResponse, error) {
	if m.verifyNotifyFn != nil {
		return m.verifyNotifyFn(ctx, params)
	}
	return nil, errors.New("not implemented")
}
func (m *mockPaymentChannel) VerifyRefundNotify(ctx context.Context, params map[string]string) (*payment.RefundResponse, error) {
	return nil, errors.New("not implemented")
}

// ---------- mockUserRepo ----------

type mockUserRepo struct {
	getByOpenIDFn func(ctx context.Context, openID string) (*User, error)
	createFn      func(ctx context.Context, u *User) (*User, error)
	getByIDFn     func(ctx context.Context, id int64) (*User, error)
}

func (m *mockUserRepo) Create(ctx context.Context, u *User) (*User, error) {
	if m.createFn != nil {
		return m.createFn(ctx, u)
	}
	u.ID = 1
	return u, nil
}
func (m *mockUserRepo) GetByID(ctx context.Context, id int64) (*User, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}
func (m *mockUserRepo) GetByOpenID(ctx context.Context, openID string) (*User, error) {
	if m.getByOpenIDFn != nil {
		return m.getByOpenIDFn(ctx, openID)
	}
	return nil, errors.New("not implemented")
}
func (m *mockUserRepo) GetByPhone(ctx context.Context, phone string) (*User, error) { return nil, errors.New("not implemented") }
func (m *mockUserRepo) Update(ctx context.Context, u *User) (*User, error) { return nil, errors.New("not implemented") }
func (m *mockUserRepo) BindPhone(ctx context.Context, id int64, phone string) error { return errors.New("not implemented") }
func (m *mockUserRepo) VerifySmsCode(ctx context.Context, phone, code string) bool { return false }
func (m *mockUserRepo) SendSmsCode(ctx context.Context, phone string) error { return errors.New("not implemented") }
