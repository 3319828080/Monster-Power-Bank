package service

import (
	"context"
	orderv1 "monster/api/order/v1"
	"monster/internal/biz"
	"monster/pkg/middleware"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

var OrderProviderSet = wire.NewSet(NewOrderService)

type OrderService struct {
	uc  *biz.OrderUsecase
	log *log.Helper
}

func NewOrderService(uc *biz.OrderUsecase, logger log.Logger) *OrderService {
	return &OrderService{uc: uc, log: log.NewHelper(logger)}
}

// PreCheckBorrow 租借前预检
func (s *OrderService) PreCheckBorrow(ctx context.Context, req *orderv1.PreCheckBorrowRequest) (*orderv1.PreCheckBorrowResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	allowed, reason, pricing, err := s.uc.PreCheckBorrow(ctx, userID, req.PowerBankId, req.StationId, req.CabinetId, req.SlotId)
	if err != nil {
		return nil, err
	}
	resp := &orderv1.PreCheckBorrowResponse{
		Allowed: allowed,
		Reason:  reason,
	}
	if pricing != nil {
		resp.DepositRequired = pricing.Deposit
		resp.StartFee = pricing.StartFee
		resp.HourlyFee = pricing.HourlyFee
		resp.DailyCap = pricing.DailyCap
	}
	return resp, nil
}

// CreateOrder 创建租借订单
func (s *OrderService) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (*orderv1.CreateOrderResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	order, err := s.uc.CreateOrder(ctx, userID, req.PowerBankId, req.StationId, req.CabinetId, req.SlotId)
	if err != nil {
		return nil, err
	}
	return &orderv1.CreateOrderResponse{
		OrderNo:       order.OrderNo,
		StartFee:      order.StartFee,
		HourlyFee:     order.HourlyFee,
		DailyCap:      order.DailyCap,
		DepositFrozen: order.Deposit,
		BorrowTime:    order.BorrowTime.Format(time.RFC3339),
	}, nil
}

// CancelOrder 取消订单
func (s *OrderService) CancelOrder(ctx context.Context, req *orderv1.CancelOrderRequest) (*orderv1.CancelOrderResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	if err := s.uc.CancelOrder(ctx, userID, req.OrderNo); err != nil {
		return nil, err
	}
	return &orderv1.CancelOrderResponse{Message: "ok"}, nil
}

// ReturnPowerBank 归还充电宝
func (s *OrderService) ReturnPowerBank(ctx context.Context, req *orderv1.ReturnPowerBankRequest) (*orderv1.ReturnPowerBankResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	couponID, _ := strconv.ParseInt(req.CouponId, 10, 64)
	order, err := s.uc.ReturnPowerBank(ctx, userID, req.OrderNo, req.StationId, req.CabinetId, req.SlotId, couponID)
	if err != nil {
		return nil, err
	}
	paidAmount := order.PaidAmount
	balanceDeducted := paidAmount
	depositDeducted := int64(0)
	if order.TotalAmount > order.PaidAmount {
		depositDeducted = order.TotalAmount - order.PaidAmount
	}
	if balanceDeducted > order.TotalAmount {
		balanceDeducted = 0
	}

	return &orderv1.ReturnPowerBankResponse{
		OrderNo:         order.OrderNo,
		DurationMinutes: int64(order.DurationMinutes),
		TotalAmount:     order.TotalAmount,
		PaidAmount:      order.PaidAmount,
		DiscountAmount:  order.DiscountAmount,
		BalanceDeducted: balanceDeducted,
		DepositDeducted: depositDeducted,
		ReturnTime:      order.ReturnTime.Format(time.RFC3339),
	}, nil
}

// ExtendRental 延长租借
func (s *OrderService) ExtendRental(ctx context.Context, req *orderv1.ExtendRentalRequest) (*orderv1.ExtendRentalResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	if err := s.uc.ExtendRental(ctx, userID, req.OrderNo, req.ExtraMinutes); err != nil {
		return nil, err
	}
	return &orderv1.ExtendRentalResponse{Message: "ok"}, nil
}

// GetCurrentOrder 查询当前订单
func (s *OrderService) GetCurrentOrder(ctx context.Context, req *orderv1.GetCurrentOrderRequest) (*orderv1.GetCurrentOrderResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	order, err := s.uc.GetCurrentOrder(ctx, userID)
	if err != nil {
		return nil, err
	}
	resp := &orderv1.GetCurrentOrderResponse{}
	if order != nil {
		resp.Order = orderToInfo(order)
	}
	return resp, nil
}

// GetRealTimeFee 实时费用查询
func (s *OrderService) GetRealTimeFee(ctx context.Context, req *orderv1.GetRealTimeFeeRequest) (*orderv1.GetRealTimeFeeResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	elapsedSeconds, currentFee, hourlyFee, dailyCap, err := s.uc.GetRealTimeFee(ctx, userID, req.OrderNo)
	if err != nil {
		return nil, err
	}
	return &orderv1.GetRealTimeFeeResponse{
		OrderNo:        req.OrderNo,
		ElapsedSeconds: elapsedSeconds,
		CurrentFee:     currentFee,
		HourlyFee:      hourlyFee,
		DailyCap:       dailyCap,
	}, nil
}

// ListOrders 历史订单查询
func (s *OrderService) ListOrders(ctx context.Context, req *orderv1.ListOrdersRequest) (*orderv1.ListOrdersResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	page := int(req.Page)
	pageSize := int(req.PageSize)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	list, total, err := s.uc.ListOrders(ctx, userID, page, pageSize, req.Status)
	if err != nil {
		return nil, err
	}
	items := make([]*orderv1.OrderInfo, 0, len(list))
	for _, o := range list {
		items = append(items, orderToInfo(o))
	}
	return &orderv1.ListOrdersResponse{
		List:     items,
		Total:    int32(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}

// CreateCompensationPayment 待赔付支付入口
func (s *OrderService) CreateCompensationPayment(ctx context.Context, req *orderv1.CreateCompensationPaymentRequest) (*orderv1.CreateCompensationPaymentResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	orderNo, amount, err := s.uc.CreateCompensationPayment(ctx, userID, req.OrderNo)
	if err != nil {
		return nil, err
	}
	return &orderv1.CreateCompensationPaymentResponse{
		PaymentUrl: "",
		Amount:     amount,
		Channel:    orderNo,
	}, nil
}

// GetStation 获取站点详情
func (s *OrderService) GetStation(ctx context.Context, req *orderv1.GetStationRequest) (*orderv1.GetStationResponse, error) {
	station, err := s.uc.GetStationDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &orderv1.GetStationResponse{
		Station: &orderv1.StationInfo{
			Id:             station.ID,
			Name:           station.Name,
			Address:        station.Address,
			Latitude:       station.Latitude,
			Longitude:      station.Longitude,
			AvailableBanks: station.AvailableBanks,
			OpenTime:       station.OpenTime,
			Description:    station.Description,
			Images:         station.Images,
		},
	}, nil
}

// ListNearbyStations 查询附近站点
func (s *OrderService) ListNearbyStations(ctx context.Context, req *orderv1.ListNearbyStationsRequest) (*orderv1.ListNearbyStationsResponse, error) {
	list, err := s.uc.ListNearbyStations(ctx, req.Latitude, req.Longitude, int(req.RadiusMeters))
	if err != nil {
		return nil, err
	}
	items := make([]*orderv1.StationInfo, 0, len(list))
	for _, st := range list {
		items = append(items, &orderv1.StationInfo{
			Id:             st.ID,
			Name:           st.Name,
			Address:        st.Address,
			Latitude:       st.Latitude,
			Longitude:      st.Longitude,
			Distance:       st.Distance,
			AvailableBanks: st.AvailableBanks,
			OpenTime:       st.OpenTime,
			Description:    st.Description,
			Images:         st.Images,
		})
	}
	return &orderv1.ListNearbyStationsResponse{List: items}, nil
}

// ScanCabinet 扫码解析机柜
func (s *OrderService) ScanCabinet(ctx context.Context, req *orderv1.ScanCabinetRequest) (*orderv1.ScanCabinetResponse, error) {
	cabinet, station, err := s.uc.ScanCabinet(ctx, req.CabinetId)
	if err != nil {
		return nil, err
	}
	resp := &orderv1.ScanCabinetResponse{
		StationId:      station.ID,
		StationName:    station.Name,
		StationAddress: station.Address,
		CabinetId:      cabinet.ID,
		CabinetNo:      cabinet.CabinetNo,
		AvailableSlots: cabinet.AvailableSlots,
		TotalSlots:     cabinet.TotalSlots,
		StartFee:       s.uc.Pricing.StartFee,
		StartMins:      int32(s.uc.Pricing.StartMins),
		HourlyFee:      s.uc.Pricing.HourlyFee,
		DailyCap:       s.uc.Pricing.DailyCap,
		Deposit:        s.uc.Pricing.Deposit,
	}
	return resp, nil
}

// GetOrderDetail 查询订单详情（含机柜仓位）
func (s *OrderService) GetOrderDetail(ctx context.Context, req *orderv1.GetOrderDetailRequest) (*orderv1.GetOrderDetailResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	order, err := s.uc.GetOrderDetail(ctx, userID, req.OrderNo)
	if err != nil {
		return nil, err
	}
	return &orderv1.GetOrderDetailResponse{
		Order: orderToDetail(order),
	}, nil
}

// GetCabinetList 查询站点下所有机柜及仓位
func (s *OrderService) GetCabinetList(ctx context.Context, req *orderv1.GetCabinetListRequest) (*orderv1.GetCabinetListResponse, error) {
	cabinets, err := s.uc.GetCabinetList(ctx, req.StationId)
	if err != nil {
		return nil, err
	}
	items := make([]*orderv1.CabinetInfo, 0, len(cabinets))
	for _, c := range cabinets {
		slots := make([]*orderv1.SlotInfo, 0, len(c.Slots))
		for _, sl := range c.Slots {
			slots = append(slots, &orderv1.SlotInfo{
				Id:           sl.ID,
				SlotNo:       sl.SlotNo,
				Status:       sl.Status,
				PowerBankId:  sl.PowerBankID,
				PowerBankNo:  sl.PowerBankNo,
				BatteryLevel: sl.BatteryLevel,
				Power:        sl.Power,
			})
		}
		items = append(items, &orderv1.CabinetInfo{
			Id:             c.ID,
			CabinetNo:      c.CabinetNo,
			Status:         c.Status,
			TotalSlots:     c.TotalSlots,
			AvailableSlots: c.AvailableSlots,
			OccupiedSlots:  c.OccupiedSlots,
			Slots:          slots,
		})
	}
	return &orderv1.GetCabinetListResponse{Cabinets: items}, nil
}

// GetDelayFee 延长租借费用预览
func (s *OrderService) GetDelayFee(ctx context.Context, req *orderv1.GetDelayFeeRequest) (*orderv1.GetDelayFeeResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	additionalFee, estimatedTotal, additionalDeposit, err := s.uc.GetDelayFee(ctx, userID, req.OrderNo, req.ExtraMinutes)
	if err != nil {
		return nil, err
	}
	return &orderv1.GetDelayFeeResponse{
		OrderNo:          req.OrderNo,
		ExtraMinutes:     req.ExtraMinutes,
		AdditionalFee:    additionalFee,
		EstimatedTotal:   estimatedTotal,
		AdditionalDeposit: additionalDeposit,
		HourlyFee:        s.uc.Pricing.HourlyFee,
		DailyCap:         s.uc.Pricing.DailyCap,
	}, nil
}

// ListReturnCabinets 查询附近可归还机柜
func (s *OrderService) ListReturnCabinets(ctx context.Context, req *orderv1.ListReturnCabinetsRequest) (*orderv1.ListReturnCabinetsResponse, error) {
	cabinets, err := s.uc.ListReturnCabinets(ctx, req.Latitude, req.Longitude, int(req.RadiusMeters))
	if err != nil {
		return nil, err
	}
	items := make([]*orderv1.ReturnCabinetInfo, 0, len(cabinets))
	for _, c := range cabinets {
		items = append(items, returnCabinetToInfo(c))
	}
	return &orderv1.ListReturnCabinetsResponse{List: items}, nil
}

// SearchCabinet 根据机柜编号搜索可归还机柜
func (s *OrderService) SearchCabinet(ctx context.Context, req *orderv1.SearchCabinetRequest) (*orderv1.SearchCabinetResponse, error) {
	c, err := s.uc.SearchReturnCabinet(ctx, req.CabinetNo)
	if err != nil {
		return nil, err
	}
	return &orderv1.SearchCabinetResponse{Cabinet: returnCabinetToInfo(c)}, nil
}

func returnCabinetToInfo(c *biz.ReturnCabinet) *orderv1.ReturnCabinetInfo {
	return &orderv1.ReturnCabinetInfo{
		CabinetID:      c.CabinetID,
		CabinetNo:      c.CabinetNo,
		StationID:      c.StationID,
		StationName:    c.StationName,
		StationAddress: c.StationAddress,
		Latitude:       c.Latitude,
		Longitude:      c.Longitude,
		Distance:       c.Distance,
		EmptySlotCount: c.EmptySlotCount,
		TotalSlots:     c.TotalSlots,
	}
}

// ---------- converters ----------

func orderToDetail(o *biz.Order) *orderv1.OrderDetail {
	detail := &orderv1.OrderDetail{
		Id:                o.ID,
		OrderNo:           o.OrderNo,
		UserId:            o.UserID,
		PowerBankId:       o.PowerBankID,
		PowerBankNo:       o.PowerBankNo,
		BorrowStationId:   o.BorrowStationID,
		BorrowStationName: o.BorrowStationName,
		BorrowCabinetId:   o.BorrowCabinetID,
		BorrowSlotId:      o.BorrowSlotID,
		ReturnStationId:   o.ReturnStationID,
		ReturnStationName: o.ReturnStationName,
		ReturnCabinetId:   o.ReturnCabinetID,
		ReturnSlotId:      o.ReturnSlotID,
		DurationMinutes:   o.DurationMinutes,
		StartFee:          o.StartFee,
		HourlyFee:         o.HourlyFee,
		DailyCap:          o.DailyCap,
		TotalAmount:       o.TotalAmount,
		PaidAmount:        o.PaidAmount,
		DiscountAmount:    o.DiscountAmount,
		Deposit:           o.Deposit,
		Status:            o.Status,
		Remark:            o.Remark,
		CreatedAt:         o.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         o.UpdatedAt.Format(time.RFC3339),
		BorrowTime:        o.BorrowTime.Format(time.RFC3339),
	}
	if o.ReturnTime != nil {
		detail.ReturnTime = o.ReturnTime.Format(time.RFC3339)
	}
	return detail
}

func orderToInfo(o *biz.Order) *orderv1.OrderInfo {
	info := &orderv1.OrderInfo{
		Id:                o.ID,
		OrderNo:           o.OrderNo,
		UserId:            o.UserID,
		PowerBankId:       o.PowerBankID,
		PowerBankNo:       o.PowerBankNo,
		BorrowStationId:   o.BorrowStationID,
		BorrowStationName: o.BorrowStationName,
		ReturnStationId:   o.ReturnStationID,
		ReturnStationName: o.ReturnStationName,
		DurationMinutes:   o.DurationMinutes,
		StartFee:          o.StartFee,
		HourlyFee:         o.HourlyFee,
		TotalAmount:       o.TotalAmount,
		PaidAmount:        o.PaidAmount,
		DiscountAmount:    o.DiscountAmount,
		DailyCap:          o.DailyCap,
		Status:            o.Status,
		Remark:            o.Remark,
		CreatedAt:         o.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         o.UpdatedAt.Format(time.RFC3339),
		BorrowTime:        o.BorrowTime.Format(time.RFC3339),
	}
	if o.ReturnTime != nil {
		info.ReturnTime = o.ReturnTime.Format(time.RFC3339)
	}
	return info
}
