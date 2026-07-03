-- 充电站种子数据 - 宿迁职业技术学院校园内 (GCJ-02 center: 118.289660, 34.020516)
-- 所有站点在中心200m范围内

-- ============================================================
-- 1. stations - 充电站（校园内各楼宇）
-- ============================================================
INSERT INTO stations (id, name, address, latitude, longitude, province, city, district, status, open_time, total_cabinets) VALUES
(1, '图书馆站',     '宿迁职业技术学院图书馆一楼大厅',     34.019900, 118.290500, '江苏省', '宿迁市', '宿城区', 1, '07:00-22:00', 2),
(2, '第一食堂站',   '宿迁职业技术学院第一食堂门口',       34.021200, 118.289300, '江苏省', '宿迁市', '宿城区', 1, '06:00-22:00', 2),
(3, '教学楼A区站',  '宿迁职业技术学院教学楼A区一楼',      34.020800, 118.290200, '江苏省', '宿迁市', '宿城区', 1, '07:00-21:00', 1),
(4, '宿舍楼1号站',  '宿迁职业技术学院宿舍楼1号楼',        34.019300, 118.288500, '江苏省', '宿迁市', '宿城区', 1, '00:00-24:00', 2),
(5, '行政楼站',     '宿迁职业技术学院行政楼一楼大厅',     34.019700, 118.288000, '江苏省', '宿迁市', '宿城区', 1, '08:00-18:00', 1),
(6, '实训楼站',     '宿迁职业技术学院实训楼一楼',         34.021400, 118.290700, '江苏省', '宿迁市', '宿城区', 1, '07:00-21:00', 1),
(7, '体育馆站',     '宿迁职业技术学院体育馆入口',         34.021700, 118.288800, '江苏省', '宿迁市', '宿城区', 1, '08:00-21:00', 1),
(8, '校门南站',     '宿迁职业技术学院南门传达室旁',       34.019100, 118.289800, '江苏省', '宿迁市', '宿城区', 1, '00:00-24:00', 1);

-- ============================================================
-- 2. cabinets - 充电柜
-- ============================================================
INSERT INTO cabinets (id, station_id, cabinet_no, total_slots, available_slots, status) VALUES
(1, 1, 'CAB-001-A', 12, 8,  1),
(2, 1, 'CAB-001-B', 8,  5,  1),
(3, 2, 'CAB-002-A', 12, 6,  1),
(4, 2, 'CAB-002-B', 8,  4,  1),
(5, 3, 'CAB-003-A', 12, 9,  1),
(6, 4, 'CAB-004-A', 12, 5,  1),
(7, 4, 'CAB-004-B', 8,  6,  1),
(8, 5, 'CAB-005-A', 8,  7,  1),
(9, 6, 'CAB-006-A', 8,  5,  1),
(10,7, 'CAB-007-A', 8,  6,  1),
(11,8, 'CAB-008-A', 12, 8,  1);

-- ============================================================
-- 3. cabinet_slots + power_banks
-- ============================================================

-- 图书馆 CAB-001-A (12 slots, 8 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(1,  1, 'A01', 1, 1),  (2,  1, 'A02', 1, 2),  (3,  1, 'A03', 1, 3),
(4,  1, 'A04', 1, 4),  (5,  1, 'A05', 1, 5),  (6,  1, 'A06', 1, 6),
(7,  1, 'A07', 1, 7),  (8,  1, 'A08', 1, 8),  (9,  1, 'A09', 0, 0),
(10, 1, 'A10', 0, 0),  (11, 1, 'A11', 0, 0),  (12, 1, 'A12', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(1,1,1,1,'PB001A01',95, '22.5W',12,1), (2,2,1,1,'PB001A02',87, '22.5W',8,1),
(3,3,1,1,'PB001A03',100,'22.5W',3,1), (4,4,1,1,'PB001A04',72, '22.5W',20,1),
(5,5,1,1,'PB001A05',60, '22.5W',15,1), (6,6,1,1,'PB001A06',98, '22.5W',6,1),
(7,7,1,1,'PB001A07',45, '22.5W',25,1), (8,8,1,1,'PB001A08',100,'22.5W',1,1);

-- 图书馆 CAB-001-B (8 slots, 5 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(13, 2, 'B01', 1, 9),  (14, 2, 'B02', 1, 10), (15, 2, 'B03', 1, 11),
(16, 2, 'B04', 1, 12), (17, 2, 'B05', 1, 13), (18, 2, 'B06', 0, 0),
(19, 2, 'B07', 0, 0),  (20, 2, 'B08', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(9,13,2,1,'PB001B01',88,'20W',10,1),  (10,14,2,1,'PB001B02',93,'20W',5,1),
(11,15,2,1,'PB001B03',76,'20W',18,1), (12,16,2,1,'PB001B04',100,'20W',2,1),
(13,17,2,1,'PB001B05',55,'20W',22,1);

-- 第一食堂 CAB-002-A (12 slots, 6 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(21, 3, 'A01', 1, 14), (22, 3, 'A02', 1, 15), (23, 3, 'A03', 1, 16),
(24, 3, 'A04', 1, 17), (25, 3, 'A05', 1, 18), (26, 3, 'A06', 1, 19),
(27, 3, 'A07', 0, 0),  (28, 3, 'A08', 0, 0),  (29, 3, 'A09', 0, 0),
(30, 3, 'A10', 0, 0),  (31, 3, 'A11', 0, 0),  (32, 3, 'A12', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(14,21,3,2,'PB002A01',90,'22.5W',7,1),  (15,22,3,2,'PB002A02',100,'22.5W',1,1),
(16,23,3,2,'PB002A03',73,'22.5W',19,1), (17,24,3,2,'PB002A04',85,'22.5W',11,1),
(18,25,3,2,'PB002A05',62,'22.5W',23,1), (19,26,3,2,'PB002A06',97,'22.5W',4,1);

-- 第一食堂 CAB-002-B (8 slots, 4 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(33, 4, 'B01', 1, 20), (34, 4, 'B02', 1, 21), (35, 4, 'B03', 1, 22),
(36, 4, 'B04', 1, 23), (37, 4, 'B05', 0, 0),  (38, 4, 'B06', 0, 0),
(39, 4, 'B07', 0, 0),  (40, 4, 'B08', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(20,33,4,2,'PB002B01',78,'20W',14,1), (21,34,4,2,'PB002B02',92,'20W',6,1),
(22,35,4,2,'PB002B03',66,'20W',21,1), (23,36,4,2,'PB002B04',100,'20W',0,1);

-- 教学楼A区 CAB-003-A (12 slots, 9 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(41, 5, 'A01', 1, 24), (42, 5, 'A02', 1, 25), (43, 5, 'A03', 1, 26),
(44, 5, 'A04', 1, 27), (45, 5, 'A05', 1, 28), (46, 5, 'A06', 1, 29),
(47, 5, 'A07', 1, 30), (48, 5, 'A08', 1, 31), (49, 5, 'A09', 1, 32),
(50, 5, 'A10', 0, 0),  (51, 5, 'A11', 0, 0),  (52, 5, 'A12', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(24,41,5,3,'PB003A01',100,'45W',2,1), (25,42,5,3,'PB003A02',82,'45W',15,1),
(26,43,5,3,'PB003A03',91,'45W',8,1),  (27,44,5,3,'PB003A04',74,'45W',18,1),
(28,45,5,3,'PB003A05',100,'45W',0,1), (29,46,5,3,'PB003A06',63,'45W',22,1),
(30,47,5,3,'PB003A07',89,'45W',11,1), (31,48,5,3,'PB003A08',96,'45W',5,1),
(32,49,5,3,'PB003A09',77,'45W',17,1);

-- 宿舍楼1号 CAB-004-A (12 slots, 5 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(53, 6, 'A01', 1, 33), (54, 6, 'A02', 1, 34), (55, 6, 'A03', 1, 35),
(56, 6, 'A04', 1, 36), (57, 6, 'A05', 1, 37), (58, 6, 'A06', 0, 0),
(59, 6, 'A07', 0, 0),  (60, 6, 'A08', 0, 0),  (61, 6, 'A09', 0, 0),
(62, 6, 'A10', 0, 0),  (63, 6, 'A11', 0, 0),  (64, 6, 'A12', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(33,53,6,4,'PB004A01',85,'22.5W',13,1), (34,54,6,4,'PB004A02',100,'22.5W',1,1),
(35,55,6,4,'PB004A03',69,'22.5W',20,1), (36,56,6,4,'PB004A04',93,'22.5W',6,1),
(37,57,6,4,'PB004A05',58,'22.5W',24,1);

-- 宿舍楼1号 CAB-004-B (8 slots, 6 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(65, 7, 'B01', 1, 38), (66, 7, 'B02', 1, 39), (67, 7, 'B03', 1, 40),
(68, 7, 'B04', 1, 41), (69, 7, 'B05', 1, 42), (70, 7, 'B06', 1, 43),
(71, 7, 'B07', 0, 0),  (72, 7, 'B08', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(38,65,7,4,'PB004B01',76,'20W',16,1), (39,66,7,4,'PB004B02',91,'20W',9,1),
(40,67,7,4,'PB004B03',64,'20W',22,1), (41,68,7,4,'PB004B04',100,'20W',0,1),
(42,69,7,4,'PB004B05',83,'20W',12,1), (43,70,7,4,'PB004B06',95,'20W',5,1);

-- 行政楼 CAB-005-A (8 slots, 7 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(73, 8, 'A01', 1, 44), (74, 8, 'A02', 1, 45), (75, 8, 'A03', 1, 46),
(76, 8, 'A04', 1, 47), (77, 8, 'A05', 1, 48), (78, 8, 'A06', 1, 49),
(79, 8, 'A07', 1, 50), (80, 8, 'A08', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(44,73,8,5,'PB005A01',80,'22.5W',14,1), (45,74,8,5,'PB005A02',94,'22.5W',7,1),
(46,75,8,5,'PB005A03',71,'22.5W',19,1), (47,76,8,5,'PB005A04',100,'22.5W',1,1),
(48,77,8,5,'PB005A05',58,'22.5W',25,1), (49,78,8,5,'PB005A06',87,'22.5W',11,1),
(50,79,8,5,'PB005A07',96,'22.5W',4,1);

-- 实训楼 CAB-006-A (8 slots, 5 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(81, 9, 'A01', 1, 51), (82, 9, 'A02', 1, 52), (83, 9, 'A03', 1, 53),
(84, 9, 'A04', 1, 54), (85, 9, 'A05', 1, 55), (86, 9, 'A06', 0, 0),
(87, 9, 'A07', 0, 0),  (88, 9, 'A08', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(51,81,9,6,'PB006A01',92,'45W',8,1),  (52,82,9,6,'PB006A02',75,'45W',17,1),
(53,83,9,6,'PB006A03',100,'45W',0,1), (54,84,9,6,'PB006A04',63,'45W',23,1),
(55,85,9,6,'PB006A05',88,'45W',12,1);

-- 体育馆 CAB-007-A (8 slots, 6 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(89, 10, 'A01', 1, 56), (90, 10, 'A02', 1, 57), (91, 10, 'A03', 1, 58),
(92, 10, 'A04', 1, 59), (93, 10, 'A05', 1, 60), (94, 10, 'A06', 1, 61),
(95, 10, 'A07', 0, 0),  (96, 10, 'A08', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(56,89,10,7,'PB007A01',86,'22.5W',10,1), (57,90,10,7,'PB007A02',99,'22.5W',3,1),
(58,91,10,7,'PB007A03',70,'22.5W',21,1), (59,92,10,7,'PB007A04',95,'22.5W',5,1),
(60,93,10,7,'PB007A05',61,'22.5W',24,1), (61,94,10,7,'PB007A06',90,'22.5W',9,1);

-- 校门南站 CAB-008-A (12 slots, 8 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(97,  11, 'A01', 1, 62), (98,  11, 'A02', 1, 63), (99,  11, 'A03', 1, 64),
(100, 11, 'A04', 1, 65), (101, 11, 'A05', 1, 66), (102, 11, 'A06', 1, 67),
(103, 11, 'A07', 1, 68), (104, 11, 'A08', 1, 69), (105, 11, 'A09', 0, 0),
(106, 11, 'A10', 0, 0),  (107, 11, 'A11', 0, 0),  (108, 11, 'A12', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(62,97,11,8,'PB008A01',100,'22.5W',4,1),  (63,98,11,8,'PB008A02',81,'22.5W',15,1),
(64,99,11,8,'PB008A03',93,'22.5W',7,1),   (65,100,11,8,'PB008A04',67,'22.5W',20,1),
(66,101,11,8,'PB008A05',89,'22.5W',11,1),  (67,102,11,8,'PB008A06',74,'22.5W',18,1),
(68,103,11,8,'PB008A07',98,'22.5W',2,1),   (69,104,11,8,'PB008A08',60,'22.5W',26,1);
