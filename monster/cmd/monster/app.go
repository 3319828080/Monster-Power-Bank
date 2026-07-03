package main

import (
	"monster/internal/biz"
	"monster/internal/conf"
)

func defaultPricing() *biz.PricingConfig {
	return &biz.PricingConfig{
		StartFee:  200,
		StartMins: 60,
		HourlyFee: 100,
		DailyCap:  2000,
		Deposit:   9900,
	}
}

// pricingConfigProvider converts conf.Business to biz.PricingConfig.
// Falls back to defaults when config is nil or fields are zero.
func pricingConfigProvider(b *conf.Business) *biz.PricingConfig {
	d := defaultPricing()
	if b == nil {
		return d
	}
	if b.DefaultStartFee > 0 {
		d.StartFee = b.DefaultStartFee
	}
	if b.DefaultStartMins > 0 {
		d.StartMins = b.DefaultStartMins
	}
	if b.DefaultHourlyFee > 0 {
		d.HourlyFee = b.DefaultHourlyFee
	}
	if b.DefaultDailyCap > 0 {
		d.DailyCap = b.DefaultDailyCap
	}
	if b.DefaultDeposit > 0 {
		d.Deposit = b.DefaultDeposit
	}
	return d
}
