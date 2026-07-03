package data

import (
	"gorm.io/gorm"
)

// SeedIfEmpty populates stations/cabinets/slots/power_banks when tables are empty.
func SeedIfEmpty(db *gorm.DB) error {
	// Migrate existing data from English to Chinese status values
	db.Exec("UPDATE power_banks SET status = '空闲' WHERE status = 'available'")
	db.Exec("UPDATE power_banks SET status = '租借中' WHERE status = 'borrowed'")
	db.Exec("UPDATE slots SET status = '空闲' WHERE status = 'occupied'")
	db.Exec("UPDATE slots SET status = '已借出' WHERE status = 'empty'")
	db.Exec("UPDATE orders SET status = '租借中' WHERE status = 'borrowed'")

	// Repair data integrity
	repairData(db)

	var stationCount int64
	if err := db.Model(&station{}).Count(&stationCount).Error; err != nil {
		return err
	}
	var pbCount int64
	if err := db.Model(&powerBank{}).Count(&pbCount).Error; err != nil {
		return err
	}
	// If both stations and power_banks already have data, skip.
	if stationCount > 0 && pbCount > 0 {
		return nil
	}

	// One or more tables are empty — clean up partial seed data and re-seed.
	db.Exec("DELETE FROM power_banks WHERE id >= 1")
	db.Exec("DELETE FROM slots WHERE id >= 1")
	db.Exec("DELETE FROM cabinets WHERE id >= 1")
	db.Exec("DELETE FROM stations WHERE id >= 1")

	// ---- Stations (14 total: 2 campus + 12 nearby) ----
	stations := []station{
		{ID: 1, Name: "图书馆站", Address: "宿迁职业技术学院图书馆一楼大厅", Latitude: 34.019900, Longitude: 118.290500, Status: 1,
			OpenTime: "07:00-22:00", Description: "紧邻图书馆借阅大厅，安静环境适合边学习边充电", Images: `["https://picsum.photos/seed/station1/400/300"]`},
		{ID: 2, Name: "第一食堂站", Address: "宿迁职业技术学院第一食堂门口", Latitude: 34.021200, Longitude: 118.289300, Status: 1,
			OpenTime: "06:30-21:30", Description: "食堂入口右侧，用餐高峰期人流量大，顺路借还方便", Images: `["https://picsum.photos/seed/station2/400/300"]`},
		{ID: 9, Name: "万达广场站", Address: "宿迁市宿城区万达广场1号门入口", Latitude: 34.025000, Longitude: 118.310000, Status: 1,
			OpenTime: "09:00-22:00", Description: "万达广场1楼服务台旁，宿迁核心商圈，逛街购物随时借还", Images: `["https://picsum.photos/seed/station9/400/300"]`},
		{ID: 10, Name: "宝龙城市广场站", Address: "宿迁市宿城区宝龙城市广场1楼大厅", Latitude: 34.010000, Longitude: 118.305000, Status: 1,
			OpenTime: "09:00-22:00", Description: "宝龙广场中庭扶梯旁，商业核心区域，客流量大", Images: `["https://picsum.photos/seed/station10/400/300"]`},
		{ID: 11, Name: "第一人民医院站", Address: "宿迁市第一人民医院门诊大厅", Latitude: 34.033000, Longitude: 118.298000, Status: 1,
			OpenTime: "00:00-24:00", Description: "门诊大厅一楼导诊台对面，24小时全天候服务，满足医患应急充电需求", Images: `["https://picsum.photos/seed/station11/400/300"]`},
		{ID: 12, Name: "宿迁学院站", Address: "宿迁学院图书馆一楼", Latitude: 34.002000, Longitude: 118.270000, Status: 1,
			OpenTime: "07:00-22:00", Description: "宿迁学院校园核心区域，服务广大师生日常充电需求", Images: `["https://picsum.photos/seed/station12/400/300"]`},
		{ID: 13, Name: "汽车客运站", Address: "宿迁汽车客运站候车大厅", Latitude: 33.998000, Longitude: 118.280000, Status: 1,
			OpenTime: "06:00-20:00", Description: "候车大厅安检口左侧，出行途中手机没电不用愁", Images: `["https://picsum.photos/seed/station13/400/300"]`},
		{ID: 14, Name: "项王故里站", Address: "宿迁市宿城区项王故里景区入口", Latitude: 33.995000, Longitude: 118.283000, Status: 1,
			OpenTime: "08:00-18:00", Description: "景区游客中心大厅，游览西楚霸王故里，随时保持电量满满", Images: `["https://picsum.photos/seed/station14/400/300"]`},
		{ID: 15, Name: "湖滨公园站", Address: "宿迁市湖滨新区湖滨公园游客中心", Latitude: 34.035000, Longitude: 118.310000, Status: 1,
			OpenTime: "08:00-18:00", Description: "游客中心进门处，骆马湖畔休闲漫步也能随时充电", Images: `["https://picsum.photos/seed/station15/400/300"]`},
		{ID: 16, Name: "水韵城站", Address: "宿迁市宿城区水韵城购物中心B1层", Latitude: 34.008000, Longitude: 118.298000, Status: 1,
			OpenTime: "09:00-22:00", Description: "B1美食广场入口，餐饮娱乐一站式，等餐时即可充满电", Images: `["https://picsum.photos/seed/station16/400/300"]`},
		{ID: 17, Name: "金鹰购物中心站", Address: "宿迁市宿城区金鹰国际购物中心1楼", Latitude: 34.015000, Longitude: 118.302000, Status: 1,
			OpenTime: "09:00-22:00", Description: "1楼化妆品区旁，高端商圈配套，购物体验更贴心", Images: `["https://picsum.photos/seed/station17/400/300"]`},
		{ID: 18, Name: "宿迁市中医院站", Address: "宿迁市中医院门诊大厅一楼", Latitude: 34.028000, Longitude: 118.272000, Status: 1,
			OpenTime: "00:00-24:00", Description: "中医门诊大厅左侧，24小时全天候守护健康与电量", Images: `["https://picsum.photos/seed/station18/400/300"]`},
		{ID: 19, Name: "宿城区政府站", Address: "宿迁市宿城区政府大楼一楼", Latitude: 34.008000, Longitude: 118.285000, Status: 1,
			OpenTime: "08:00-18:00", Description: "政府大楼一楼大厅电梯旁，政务办事等待时免费充电", Images: `["https://picsum.photos/seed/station19/400/300"]`},
		{ID: 20, Name: "三台山森林公园站", Address: "宿迁市三台山国家森林公园南门", Latitude: 34.050000, Longitude: 118.310000, Status: 1,
			OpenTime: "08:00-17:00", Description: "景区南门游客服务中心，花海森林与满电同行，旅途更精彩", Images: `["https://picsum.photos/seed/station20/400/300"]`},
	}

	// Cabinet definitions: [cabinet_id, station_id, cabinet_no, total_slots, occupied]
	type cabDef struct {
		ID        int64
		StationID int64
		CabinetNo string
		Total     int
		Occupied  int
	}
	cabinets := []cabDef{
		{1, 1, "CAB-001-A", 12, 8}, {2, 1, "CAB-001-B", 8, 5},
		{3, 2, "CAB-002-A", 12, 6}, {4, 2, "CAB-002-B", 8, 4},
		{12, 9, "CAB-009-A", 12, 7}, {13, 9, "CAB-009-B", 8, 4},
		{14, 10, "CAB-010-A", 12, 5}, {15, 10, "CAB-010-B", 8, 5},
		{16, 11, "CAB-011-A", 12, 6}, {17, 11, "CAB-011-B", 8, 6},
		{18, 12, "CAB-012-A", 12, 8}, {19, 12, "CAB-012-B", 8, 5},
		{20, 13, "CAB-013-A", 8, 4},
		{21, 14, "CAB-014-A", 8, 5},
		{22, 15, "CAB-015-A", 12, 8}, {23, 15, "CAB-015-B", 8, 6},
		{24, 16, "CAB-016-A", 12, 6}, {25, 16, "CAB-016-B", 8, 5},
		{26, 17, "CAB-017-A", 12, 7}, {27, 17, "CAB-017-B", 8, 4},
		{28, 18, "CAB-018-A", 8, 5},
		{29, 19, "CAB-019-A", 8, 4},
		{30, 20, "CAB-020-A", 8, 6},
	}

	batteries := []int32{95, 87, 100, 72, 60, 98, 45, 100, 88, 93, 76, 100, 55, 90, 85, 73, 62, 97, 80, 94}
	letters := "ABCDEFGHIJKL"

	slotID := int64(1)
	pbID := int64(1)

	for _, st := range stations {
		if err := db.Create(&st).Error; err != nil {
			return err
		}
	}

	for _, c := range cabinets {
		if err := db.Create(&cabinet{
			ID: c.ID, StationID: c.StationID, CabinetNo: c.CabinetNo, Status: 1,
		}).Error; err != nil {
			return err
		}

		for i := 0; i < c.Total; i++ {
			s := slot{
				ID: slotID, CabinetID: c.ID, StationID: c.StationID,
				SlotNo: string(letters[i]) + pad2(i+1),
			}
			if i < c.Occupied {
				s.Status = "空闲"
				s.PowerBankID = pbID

				pb := powerBank{
					ID: pbID, DeviceNo: "PB" + pad3(int(c.ID)) + string(letters[i]) + pad2(i+1),
					StationID: c.StationID, CabinetID: c.ID, SlotID: slotID,
					Status: "空闲", BatteryLevel: batteries[(pbID-1)%int64(len(batteries))],
				}
				if err := db.Create(&pb).Error; err != nil {
					return err
				}
				pbID++
			} else {
				s.Status = "已借出"
			}
			if err := db.Create(&s).Error; err != nil {
				return err
			}
			slotID++
		}
	}
	return nil
}

// repairData fixes data inconsistencies that can arise from code changes or partial failures.
func repairData(db *gorm.DB) {
	// 1. Fix 租借中 power banks still showing at a station (should have no location)
	db.Exec(`UPDATE power_banks SET station_id = 0, cabinet_id = 0, slot_id = 0
		WHERE status = '租借中' AND station_id != 0`)

	// 2. Fix 空闲 slots whose power_bank_id points to a non-existent power bank
	db.Exec(`UPDATE slots SET power_bank_id = 0
		WHERE status = '空闲' AND power_bank_id > 0
		AND power_bank_id NOT IN (SELECT id FROM power_banks)`)

	// 3. Fix 空闲 power banks whose location doesn't match any slot
	//    Update power bank location to match the slot that claims it
	db.Exec(`UPDATE power_banks pb
		JOIN slots sl ON sl.power_bank_id = pb.id
		SET pb.station_id = sl.station_id, pb.cabinet_id = sl.cabinet_id, pb.slot_id = sl.id
		WHERE pb.status = '空闲'
		AND (pb.station_id != sl.station_id OR pb.cabinet_id != sl.cabinet_id OR pb.slot_id != sl.id)`)

	// 4. Fix 空闲 slots whose power_bank_id doesn't match the power bank that claims to be there
	//    Update slot to reference the power bank that's actually at this location
	db.Exec(`UPDATE slots sl
		JOIN power_banks pb ON pb.station_id = sl.station_id AND pb.cabinet_id = sl.cabinet_id AND pb.slot_id = sl.id
		SET sl.power_bank_id = pb.id, sl.status = '空闲'
		WHERE pb.status = '空闲' AND sl.power_bank_id != pb.id`)

	// 5. For 空闲 power banks with no slot claiming them, find a homeless slot
	db.Exec(`UPDATE power_banks pb
		JOIN slots sl ON sl.station_id = pb.station_id AND sl.cabinet_id = pb.cabinet_id AND sl.status = '已借出'
		SET pb.slot_id = sl.id, sl.power_bank_id = pb.id, sl.status = '空闲'
		WHERE pb.status = '空闲' AND pb.id NOT IN (
			SELECT power_bank_id FROM slots WHERE power_bank_id > 0
		)`)
}

func pad2(n int) string { s := "00" + itoa(n); return s[len(s)-2:] }
func pad3(n int) string { s := "000" + itoa(n); return s[len(s)-3:] }
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	r := ""
	for n > 0 {
		r = string(rune('0'+n%10)) + r
		n /= 10
	}
	return r
}
