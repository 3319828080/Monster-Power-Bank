package payment

// NewIntegratedManager creates a Manager and registers all available channels.
func NewIntegratedManager(wechatAppID, wechatMchID, wechatAPIKey, wechatNotifyURL string,
	alipayAppID, alipayPrivateKey, alipayReturnURL, alipayNotifyURL string) (*Manager, error) {

	m := NewManager()

	// WeChat Pay channel
	if wechatAppID != "" && wechatMchID != "" {
		m.Register(NewWechatPayChannel(wechatAppID, wechatMchID, wechatAPIKey, wechatNotifyURL))
	}

	// Alipay channel
	if alipayAppID != "" && alipayPrivateKey != "" {
		alipayCh, err := NewAlipayChannel(alipayAppID, alipayPrivateKey, alipayReturnURL, alipayNotifyURL)
		if err != nil {
			return nil, err
		}
		m.Register(alipayCh)
	}

	return m, nil
}
