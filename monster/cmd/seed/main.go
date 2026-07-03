package main

import (
	"fmt"
	"os"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:4ay1nkal3u8ed77y@tcp(115.190.199.25:3306)/monster?parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect failed: %v\n", err)
		os.Exit(1)
	}

	sqlBytes, err := os.ReadFile("migrations/000004_seed_full.sql")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read sql failed: %v\n", err)
		os.Exit(1)
	}

	lines := strings.Split(string(sqlBytes), "\n")
	var stmt strings.Builder
	n := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "--") {
			continue
		}
		stmt.WriteString(line)
		if strings.HasSuffix(trimmed, ";") {
			s := strings.TrimSpace(stmt.String())
			if s != "" {
				if err := db.Exec(s).Error; err != nil {
					preview := s
					if len(preview) > 100 {
						preview = preview[:100]
					}
					fmt.Fprintf(os.Stderr, "stmt %d err: %v\n  %s\n", n+1, err, preview)
					os.Exit(1)
				}
				n++
			}
			stmt.Reset()
		} else {
			stmt.WriteString(" ")
		}
	}

	fmt.Printf("Done. %d statements executed.\n", n)
	fmt.Println("14 stations, 22 cabinets, 220 slots, 129 power banks.")
}
