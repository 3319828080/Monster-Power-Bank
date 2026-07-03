package data

import (
	"context"
	"monster/internal/conf"
	"monster/pkg/sms"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var ProviderSet = wire.NewSet(NewData, NewDB, NewRedis, NewRedisLock, NewUserRepo, NewWalletRepo, NewOrderRepo, NewPaymentRepo)

type Data struct {
	db  *gorm.DB
	rdb *redis.Client
	sms sms.Sender
	lock *RedisLock
}

type txKey struct{}

// DB returns the underlying *gorm.DB, using transactional context if available.
func (d *Data) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	if ok && tx != nil {
		return tx
	}
	return d.db.WithContext(ctx)
}

func NewData(c *conf.Data, db *gorm.DB, rdb *redis.Client, smsClient sms.Sender, logger log.Logger) (*Data, func(), error) {
	hl := log.NewHelper(logger)
	hl.Info("initializing data layer")

	// Ping Redis
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		hl.Warnf("redis ping failed: %v", err)
	} else {
		hl.Infof("redis connected successfully, addr: %s", c.Redis.Addr)
	}

	cleanup := func() {
		hl.Info("closing data resources")
		if rdb != nil {
			rdb.Close()
		}
	}
	return &Data{db: db, rdb: rdb, sms: smsClient, lock: NewRedisLock(rdb)}, cleanup, nil
}

func NewDB(c *conf.Data, l log.Logger) (*gorm.DB, error) {
	hl := log.NewHelper(l)
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	hl.Infof("database connected successfully, addr: %s", c.Database.Source)

	// Auto migrate tables
	if err := db.AutoMigrate(
		&user{},
		&wallet{},
		&walletTransaction{},
		&order{},
		&orderItem{},
		&powerBank{},
		&station{},
		&cabinet{},
		&slot{},
		&coupon{},
		&paymentModel{},
		&refundModel{},
		&notificationModel{},
	); err != nil {
		return nil, err
	}
	hl.Info("database auto migration completed")

	if err := SeedIfEmpty(db); err != nil {
		return nil, err
	}
	hl.Info("seed data check completed")

	return db, nil
}

func NewRedis(c *conf.Data) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password,
		DB:           int(c.Redis.Database),
		ReadTimeout:  time.Duration(c.Redis.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(c.Redis.WriteTimeout) * time.Millisecond,
	})
}
