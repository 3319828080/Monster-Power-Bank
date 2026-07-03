package service

import (
	"context"
	userv1 "monster/api/user/v1"
	"monster/internal/biz"
	"monster/pkg/jwt"
	"monster/pkg/middleware"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

var UserProviderSet = wire.NewSet(NewUserService)

type UserService struct {
	uc  *biz.UserUsecase
	wuc *biz.WalletUsecase
	jwt *jwt.JWT
	log *log.Helper
}

func NewUserService(uc *biz.UserUsecase, wuc *biz.WalletUsecase, j *jwt.JWT, logger log.Logger) *UserService {
	return &UserService{uc: uc, wuc: wuc, jwt: j, log: log.NewHelper(logger)}
}

// Login 微信小程序登录
func (s *UserService) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	user, err := s.uc.LoginOrRegister(ctx, "mock_"+req.Code, "", req.Nickname, req.Avatar)
	if err != nil {
		return nil, err
	}
	s.wuc.EnsureWallet(ctx, user.ID)
	token, err := s.jwt.GenerateToken(user.ID, user.OpenID, user.Phone)
	if err != nil {
		return nil, err
	}
	return &userv1.LoginResponse{Token: token, User: s.userToInfo(user)}, nil
}

// LoginPhone 手机号验证码登录
func (s *UserService) LoginPhone(ctx context.Context, req *userv1.LoginPhoneRequest) (*userv1.LoginResponse, error) {
	user, err := s.uc.LoginByPhone(ctx, req.Phone, req.Code, req.Nickname, req.Avatar)
	if err != nil {
		return nil, err
	}
	s.wuc.EnsureWallet(ctx, user.ID)
	token, err := s.jwt.GenerateToken(user.ID, user.OpenID, user.Phone)
	if err != nil {
		return nil, err
	}
	return &userv1.LoginResponse{Token: token, User: s.userToInfo(user)}, nil
}

// BindPhone 绑定手机号
func (s *UserService) BindPhone(ctx context.Context, req *userv1.BindPhoneRequest) (*userv1.BindPhoneResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	if err := s.uc.BindPhone(ctx, userID, req.Phone); err != nil {
		return nil, err
	}
	return &userv1.BindPhoneResponse{Message: "ok"}, nil
}

// GetProfile 获取个人信息
func (s *UserService) GetProfile(ctx context.Context, req *userv1.GetProfileRequest) (*userv1.GetProfileResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	user, err := s.uc.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &userv1.GetProfileResponse{User: s.userToInfo(user)}, nil
}

// UpdateProfile 更新个人信息
func (s *UserService) UpdateProfile(ctx context.Context, req *userv1.UpdateProfileRequest) (*userv1.UpdateProfileResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	user, err := s.uc.UpdateProfile(ctx, &biz.User{
		ID:       userID,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Gender:   req.Gender,
	})
	if err != nil {
		return nil, err
	}
	return &userv1.UpdateProfileResponse{User: s.userToInfo(user)}, nil
}

// SendSmsCode 发送短信验证码
func (s *UserService) SendSmsCode(ctx context.Context, req *userv1.SendSmsCodeRequest) (*userv1.SendSmsCodeResponse, error) {
	if err := s.uc.SendSmsCode(ctx, req.Phone); err != nil {
		return nil, err
	}
	return &userv1.SendSmsCodeResponse{Message: "ok"}, nil
}

// Logout 退出登录
func (s *UserService) Logout(ctx context.Context, req *userv1.LogoutRequest) (*userv1.LogoutResponse, error) {
	return &userv1.LogoutResponse{Message: "ok"}, nil
}

// GetWallet 查询钱包资产（余额、押金、累计消费）
func (s *UserService) GetWallet(ctx context.Context, req *userv1.GetWalletRequest) (*userv1.GetWalletResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	wallet, err := s.wuc.EnsureWallet(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &userv1.GetWalletResponse{Wallet: s.walletToInfo(wallet)}, nil
}

// ListTransactions 查询资金流水账单
func (s *UserService) ListTransactions(ctx context.Context, req *userv1.ListTransactionsRequest) (*userv1.ListTransactionsResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	page := int(req.Page)
	pageSize := int(req.PageSize)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	list, total, err := s.wuc.ListTransactions(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	items := make([]*userv1.TransactionInfo, 0, len(list))
	for _, t := range list {
		items = append(items, &userv1.TransactionInfo{
			Id:            t.ID,
			UserId:        t.UserID,
			WalletId:      t.WalletID,
			OrderId:       t.OrderID,
			TradeType:     t.TradeType,
			Amount:        t.Amount,
			BalanceBefore: t.BalanceBefore,
			BalanceAfter:  t.BalanceAfter,
			Remark:        t.Remark,
			CreatedAt:     t.CreatedAt.Format(time.RFC3339),
		})
	}
	return &userv1.ListTransactionsResponse{
		List:     items,
		Total:    int32(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}

// PayDeposit 缴纳押金
func (s *UserService) PayDeposit(ctx context.Context, req *userv1.PayDepositRequest) (*userv1.PayDepositResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	if err := s.wuc.PayDeposit(ctx, userID, req.Amount); err != nil {
		return nil, err
	}
	return &userv1.PayDepositResponse{Message: "ok"}, nil
}

// RechargeDeposit 押金快速充值
func (s *UserService) RechargeDeposit(ctx context.Context, req *userv1.RechargeDepositRequest) (*userv1.RechargeDepositResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	if err := s.wuc.RechargeDeposit(ctx, userID, req.Amount); err != nil {
		return nil, err
	}
	return &userv1.RechargeDepositResponse{Message: "ok"}, nil
}

// GetDepositInfo 查询押金规则和缴纳状态
func (s *UserService) GetDepositInfo(ctx context.Context, req *userv1.GetDepositInfoRequest) (*userv1.GetDepositInfoResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	hasDeposit, depositAmount, depositFrozen, depositEnough, rules, err := s.wuc.GetDepositInfo(ctx, userID, 9900) // default deposit 99 yuan
	if err != nil {
		return nil, err
	}
	wallet, _ := s.wuc.GetWallet(ctx, userID)
	walletBalance := int64(0)
	if wallet != nil {
		walletBalance = wallet.Balance
	}
	return &userv1.GetDepositInfoResponse{
		HasDeposit:    hasDeposit,
		DepositAmount: depositAmount,
		DepositFrozen: depositFrozen,
		WalletBalance: walletBalance,
		DepositEnough: depositEnough,
		Rules:         rules,
	}, nil
}

// RefundDeposit 退还押金
func (s *UserService) RefundDeposit(ctx context.Context, req *userv1.RefundDepositRequest) (*userv1.RefundDepositResponse, error) {
	userID := middleware.UserIDFromContext(ctx)
	if err := s.wuc.RefundDeposit(ctx, userID); err != nil {
		return nil, err
	}
	return &userv1.RefundDepositResponse{Message: "ok"}, nil
}

func (s *UserService) walletToInfo(w *biz.Wallet) *userv1.WalletInfo {
	return &userv1.WalletInfo{
		Id:            w.ID,
		UserId:        w.UserID,
		Balance:       w.Balance,
		Frozen:        w.Frozen,
		TotalRecharge: w.TotalRecharge,
		TotalConsume:  w.TotalConsume,
		Status:        w.Status,
		CreatedAt:     w.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     w.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *UserService) userToInfo(u *biz.User) *userv1.UserInfo {
	return &userv1.UserInfo{
		Id:        u.ID,
		OpenId:    u.OpenID,
		UnionId:   u.UnionID,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Phone:     u.Phone,
		Gender:    u.Gender,
		Status:    u.Status,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}
