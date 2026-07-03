package main

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:4ay1nkal3u8ed77y@tcp(115.190.199.25:3306)/monster?parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Fprintln(os.Stderr, "connect error:", err)
		os.Exit(1)
	}

	// 1. Fix wallet: restore balance and frozen
	db.Exec("UPDATE wallets SET balance = 10000, frozen = 9900 WHERE id = 12 AND user_id = 15")
	fmt.Println("1. Wallet fixed: Balance=10000, Frozen=9900")

	// 2. Clean up erroneous order items from failed return attempts
	db.Exec("DELETE FROM order_items WHERE order_id = 4 AND fee_type = 'rental_fee'")
	fmt.Println("2. Erroneous order_items cleaned up")

	// 3. Verify state
	fmt.Println("\n=== Verification ===")

	type wallet struct {
		ID      int64
		UserID  int64
		Balance int64
		Frozen  int64
	}
	var w wallet
	db.Table("wallets").Where("id = 12").First(&w)
	fmt.Printf("Wallet: Balance=%d Frozen=%d\n", w.Balance, w.Frozen)

	type order struct {
		ID      int64
		Status  string
	}
	var o order
	db.Table("orders").Where("id = 4").First(&o)
	fmt.Printf("Order 4: Status=%s\n", o.Status)

	type powerBank struct {
		ID     int64
		Status string
	}
	var pb powerBank
	db.Table("power_banks").Where("id = 14").First(&pb)
	fmt.Printf("PB 14: Status=%s\n", pb.Status)

	type slot struct {
		ID     int64
		Status string
	}
	var sl slot
	db.Table("slots").Where("id = 21").First(&sl)
	fmt.Printf("Slot 21: Status=%s\n", sl.Status)

	var itemCount int64
	db.Table("order_items").Where("order_id = 4").Count(&itemCount)
	fmt.Printf("Order items for order 4: %d\n", itemCount)

	fmt.Println("\nData fixed. Ready for code fix.")
}
