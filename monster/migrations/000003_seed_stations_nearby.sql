-- 充电站种子数据 - 宿迁职业技术学院周边 (GCJ-02 center: 118.289660, 34.020516)
-- 周边1-5km范围站点，Run after 000002_seed_stations.sql

INSERT INTO stations (id, name, address, latitude, longitude, province, city, district, status, open_time, total_cabinets) VALUES
(9,  '万达广场站',       '宿迁市宿城区万达广场1号门入口',         34.025000, 118.310000, '江苏省', '宿迁市', '宿城区', 1, '09:00-22:00', 2),
(10, '宝龙城市广场站',   '宿迁市宿城区宝龙城市广场1楼大厅',     34.010000, 118.305000, '江苏省', '宿迁市', '宿城区', 1, '09:00-22:00', 2),
(11, '第一人民医院站',   '宿迁市第一人民医院门诊大厅',           34.033000, 118.298000, '江苏省', '宿迁市', '宿城区', 1, '00:00-24:00', 2),
(12, '宿迁学院站',       '宿迁学院图书馆一楼',                   34.002000, 118.270000, '江苏省', '宿迁市', '宿城区', 1, '07:00-22:00', 2),
(13, '汽车客运站',       '宿迁汽车客运站候车大厅',               33.998000, 118.280000, '江苏省', '宿迁市', '宿城区', 1, '06:00-20:00', 1),
(14, '项王故里站',       '宿迁市宿城区项王故里景区入口',         33.995000, 118.283000, '江苏省', '宿迁市', '宿城区', 1, '08:00-18:00', 1),
(15, '湖滨公园站',       '宿迁市湖滨新区湖滨公园游客中心',       34.035000, 118.310000, '江苏省', '宿迁市', '湖滨新区', 1, '08:00-18:00', 2),
(16, '水韵城站',         '宿迁市宿城区水韵城购物中心B1层',       34.008000, 118.298000, '江苏省', '宿迁市', '宿城区', 1, '09:00-22:00', 2),
(17, '金鹰购物中心站',   '宿迁市宿城区金鹰国际购物中心1楼',     34.015000, 118.302000, '江苏省', '宿迁市', '宿城区', 1, '09:00-22:00', 2),
(18, '宿迁市中医院站',   '宿迁市中医院门诊大厅一楼',             34.028000, 118.272000, '江苏省', '宿迁市', '宿城区', 1, '00:00-24:00', 1),
(19, '宿城区政府站',     '宿迁市宿城区政府大楼一楼',             34.008000, 118.285000, '江苏省', '宿迁市', '宿城区', 1, '08:00-18:00', 1),
(20, '三台山森林公园站', '宿迁市三台山国家森林公园南门游客中心', 34.050000, 118.310000, '江苏省', '宿迁市', '湖滨新区', 1, '08:00-17:00', 1);

INSERT INTO cabinets (id, station_id, cabinet_no, total_slots, available_slots, status) VALUES
-- 万达广场 (9)
(12, 9,  'CAB-009-A', 12, 7, 1),
(13, 9,  'CAB-009-B', 8,  4, 1),
-- 宝龙城市广场 (10)
(14, 10, 'CAB-010-A', 12, 5, 1),
(15, 10, 'CAB-010-B', 8,  5, 1),
-- 第一人民医院 (11)
(16, 11, 'CAB-011-A', 12, 6, 1),
(17, 11, 'CAB-011-B', 8,  6, 1),
-- 宿迁学院 (12)
(18, 12, 'CAB-012-A', 12, 8, 1),
(19, 12, 'CAB-012-B', 8,  5, 1),
-- 汽车客运站 (13)
(20, 13, 'CAB-013-A', 8,  4, 1),
-- 项王故里 (14)
(21, 14, 'CAB-014-A', 8,  5, 1),
-- 湖滨公园 (15)
(22, 15, 'CAB-015-A', 12, 8, 1),
(23, 15, 'CAB-015-B', 8,  6, 1),
-- 水韵城 (16)
(24, 16, 'CAB-016-A', 12, 6, 1),
(25, 16, 'CAB-016-B', 8,  5, 1),
-- 金鹰购物中心 (17)
(26, 17, 'CAB-017-A', 12, 7, 1),
(27, 17, 'CAB-017-B', 8,  4, 1),
-- 中医院 (18)
(28, 18, 'CAB-018-A', 8,  5, 1),
-- 宿城区政府 (19)
(29, 19, 'CAB-019-A', 8,  4, 1),
-- 三台山森林公园 (20)
(30, 20, 'CAB-020-A', 8,  6, 1);

-- 万达广场 CAB-012 (12 slots, 7 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(109, 12, 'A01', 1, 70), (110, 12, 'B02', 1, 71), (111, 12, 'C03', 1, 72), (112, 12, 'D04', 1, 73), (113, 12, 'E05', 1, 74), (114, 12, 'F06', 1, 75), (115, 12, 'G07', 1, 76), (116, 12, 'H08', 0, 0);
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(117, 12, 'I09', 0, 0), (118, 12, 'J10', 0, 0), (119, 12, 'K11', 0, 0), (120, 12, 'L12', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(70, 109, 12, 9, 'PB012A01', 95, '22.5W', 12, 1), (71, 110, 12, 9, 'PB012B02', 87, '20W', 8, 1), (72, 111, 12, 9, 'PB012C03', 100, '45W', 3, 1), (73, 112, 12, 9, 'PB012D04', 72, '22.5W', 20, 1), (74, 113, 12, 9, 'PB012E05', 60, '20W', 15, 1);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(75, 114, 12, 9, 'PB012F06', 98, '22.5W', 6, 1), (76, 115, 12, 9, 'PB012G07', 45, '20W', 25, 1);

-- 万达广场 CAB-013 (8 slots, 4 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(121, 13, 'A01', 1, 77), (122, 13, 'B02', 1, 78), (123, 13, 'C03', 1, 79), (124, 13, 'D04', 1, 80), (125, 13, 'E05', 0, 0), (126, 13, 'F06', 0, 0), (127, 13, 'G07', 0, 0), (128, 13, 'H08', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(77, 121, 13, 9, 'PB013A01', 100, '45W', 1, 1), (78, 122, 13, 9, 'PB013B02', 88, '22.5W', 10, 1), (79, 123, 13, 9, 'PB013C03', 93, '20W', 5, 1), (80, 124, 13, 9, 'PB013D04', 76, '22.5W', 18, 1);

-- 宝龙城市广场 CAB-014 (12 slots, 5 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(129, 14, 'A01', 1, 81), (130, 14, 'B02', 1, 82), (131, 14, 'C03', 1, 83), (132, 14, 'D04', 1, 84), (133, 14, 'E05', 1, 85), (134, 14, 'F06', 0, 0), (135, 14, 'G07', 0, 0), (136, 14, 'H08', 0, 0);
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(137, 14, 'I09', 0, 0), (138, 14, 'J10', 0, 0), (139, 14, 'K11', 0, 0), (140, 14, 'L12', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(81, 129, 14, 10, 'PB014A01', 100, '20W', 2, 1), (82, 130, 14, 10, 'PB014B02', 55, '45W', 22, 1), (83, 131, 14, 10, 'PB014C03', 90, '22.5W', 7, 1), (84, 132, 14, 10, 'PB014D04', 85, '20W', 11, 1), (85, 133, 14, 10, 'PB014E05', 73, '22.5W', 19, 1);

-- 宝龙城市广场 CAB-015 (8 slots, 5 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(141, 15, 'A01', 1, 86), (142, 15, 'B02', 1, 87), (143, 15, 'C03', 1, 88), (144, 15, 'D04', 1, 89), (145, 15, 'E05', 1, 90), (146, 15, 'F06', 0, 0), (147, 15, 'G07', 0, 0), (148, 15, 'H08', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(86, 141, 15, 10, 'PB015A01', 62, '20W', 23, 1), (87, 142, 15, 10, 'PB015B02', 97, '45W', 4, 1), (88, 143, 15, 10, 'PB015C03', 80, '22.5W', 14, 1), (89, 144, 15, 10, 'PB015D04', 94, '20W', 6, 1), (90, 145, 15, 10, 'PB015E05', 95, '22.5W', 12, 1);

-- 第一人民医院 CAB-016 (12 slots, 6 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(149, 16, 'A01', 1, 91), (150, 16, 'B02', 1, 92), (151, 16, 'C03', 1, 93), (152, 16, 'D04', 1, 94), (153, 16, 'E05', 1, 95), (154, 16, 'F06', 1, 96), (155, 16, 'G07', 0, 0), (156, 16, 'H08', 0, 0);
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(157, 16, 'I09', 0, 0), (158, 16, 'J10', 0, 0), (159, 16, 'K11', 0, 0), (160, 16, 'L12', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(91, 149, 16, 11, 'PB016A01', 87, '20W', 8, 1), (92, 150, 16, 11, 'PB016B02', 100, '45W', 3, 1), (93, 151, 16, 11, 'PB016C03', 72, '22.5W', 20, 1), (94, 152, 16, 11, 'PB016D04', 60, '20W', 15, 1), (95, 153, 16, 11, 'PB016E05', 98, '22.5W', 6, 1);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(96, 154, 16, 11, 'PB016F06', 45, '20W', 25, 1);

-- 第一人民医院 CAB-017 (8 slots, 6 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(161, 17, 'A01', 1, 97), (162, 17, 'B02', 1, 98), (163, 17, 'C03', 1, 99), (164, 17, 'D04', 1, 100), (165, 17, 'E05', 1, 101), (166, 17, 'F06', 1, 102), (167, 17, 'G07', 0, 0), (168, 17, 'H08', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(97, 161, 17, 11, 'PB017A01', 100, '45W', 1, 1), (98, 162, 17, 11, 'PB017B02', 88, '22.5W', 10, 1), (99, 163, 17, 11, 'PB017C03', 93, '20W', 5, 1), (100, 164, 17, 11, 'PB017D04', 76, '22.5W', 18, 1), (101, 165, 17, 11, 'PB017E05', 100, '20W', 2, 1);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(102, 166, 17, 11, 'PB017F06', 55, '45W', 22, 1);

-- 宿迁学院 CAB-018 (12 slots, 8 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(169, 18, 'A01', 1, 103), (170, 18, 'B02', 1, 104), (171, 18, 'C03', 1, 105), (172, 18, 'D04', 1, 106), (173, 18, 'E05', 1, 107), (174, 18, 'F06', 1, 108), (175, 18, 'G07', 1, 109), (176, 18, 'H08', 1, 110);
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(177, 18, 'I09', 0, 0), (178, 18, 'J10', 0, 0), (179, 18, 'K11', 0, 0), (180, 18, 'L12', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(103, 169, 18, 12, 'PB018A01', 90, '22.5W', 7, 1), (104, 170, 18, 12, 'PB018B02', 85, '20W', 11, 1), (105, 171, 18, 12, 'PB018C03', 73, '22.5W', 19, 1), (106, 172, 18, 12, 'PB018D04', 62, '20W', 23, 1), (107, 173, 18, 12, 'PB018E05', 97, '45W', 4, 1);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(108, 174, 18, 12, 'PB018F06', 80, '22.5W', 14, 1), (109, 175, 18, 12, 'PB018G07', 94, '20W', 6, 1), (110, 176, 18, 12, 'PB018H08', 95, '22.5W', 12, 1);

-- 宿迁学院 CAB-019 (8 slots, 5 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(181, 19, 'A01', 1, 111), (182, 19, 'B02', 1, 112), (183, 19, 'C03', 1, 113), (184, 19, 'D04', 1, 114), (185, 19, 'E05', 1, 115), (186, 19, 'F06', 0, 0), (187, 19, 'G07', 0, 0), (188, 19, 'H08', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(111, 181, 19, 12, 'PB019A01', 87, '20W', 8, 1), (112, 182, 19, 12, 'PB019B02', 100, '45W', 3, 1), (113, 183, 19, 12, 'PB019C03', 72, '22.5W', 20, 1), (114, 184, 19, 12, 'PB019D04', 60, '20W', 15, 1), (115, 185, 19, 12, 'PB019E05', 98, '22.5W', 6, 1);

-- 汽车客运站 CAB-020 (8 slots, 4 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(189, 20, 'A01', 1, 116), (190, 20, 'B02', 1, 117), (191, 20, 'C03', 1, 118), (192, 20, 'D04', 1, 119), (193, 20, 'E05', 0, 0), (194, 20, 'F06', 0, 0), (195, 20, 'G07', 0, 0), (196, 20, 'H08', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(116, 189, 20, 13, 'PB020A01', 45, '20W', 25, 1), (117, 190, 20, 13, 'PB020B02', 100, '45W', 1, 1), (118, 191, 20, 13, 'PB020C03', 88, '22.5W', 10, 1), (119, 192, 20, 13, 'PB020D04', 93, '20W', 5, 1);

-- 项王故里 CAB-021 (8 slots, 5 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(197, 21, 'A01', 1, 120), (198, 21, 'B02', 1, 121), (199, 21, 'C03', 1, 122), (200, 21, 'D04', 1, 123), (201, 21, 'E05', 1, 124), (202, 21, 'F06', 0, 0), (203, 21, 'G07', 0, 0), (204, 21, 'H08', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(120, 197, 21, 14, 'PB021A01', 76, '22.5W', 18, 1), (121, 198, 21, 14, 'PB021B02', 100, '20W', 2, 1), (122, 199, 21, 14, 'PB021C03', 55, '45W', 22, 1), (123, 200, 21, 14, 'PB021D04', 90, '22.5W', 7, 1), (124, 201, 21, 14, 'PB021E05', 85, '20W', 11, 1);

-- 湖滨公园 CAB-022 (12 slots, 8 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(205, 22, 'A01', 1, 125), (206, 22, 'B02', 1, 126), (207, 22, 'C03', 1, 127), (208, 22, 'D04', 1, 128), (209, 22, 'E05', 1, 129), (210, 22, 'F06', 1, 130), (211, 22, 'G07', 1, 131), (212, 22, 'H08', 1, 132);
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(213, 22, 'I09', 0, 0), (214, 22, 'J10', 0, 0), (215, 22, 'K11', 0, 0), (216, 22, 'L12', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(125, 205, 22, 15, 'PB022A01', 73, '22.5W', 19, 1), (126, 206, 22, 15, 'PB022B02', 62, '20W', 23, 1), (127, 207, 22, 15, 'PB022C03', 97, '45W', 4, 1), (128, 208, 22, 15, 'PB022D04', 80, '22.5W', 14, 1), (129, 209, 22, 15, 'PB022E05', 94, '20W', 6, 1);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(130, 210, 22, 15, 'PB022F06', 95, '22.5W', 12, 1), (131, 211, 22, 15, 'PB022G07', 87, '20W', 8, 1), (132, 212, 22, 15, 'PB022H08', 100, '45W', 3, 1);

-- 湖滨公园 CAB-023 (8 slots, 6 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(217, 23, 'A01', 1, 133), (218, 23, 'B02', 1, 134), (219, 23, 'C03', 1, 135), (220, 23, 'D04', 1, 136), (221, 23, 'E05', 1, 137), (222, 23, 'F06', 1, 138), (223, 23, 'G07', 0, 0), (224, 23, 'H08', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(133, 217, 23, 15, 'PB023A01', 72, '22.5W', 20, 1), (134, 218, 23, 15, 'PB023B02', 60, '20W', 15, 1), (135, 219, 23, 15, 'PB023C03', 98, '22.5W', 6, 1), (136, 220, 23, 15, 'PB023D04', 45, '20W', 25, 1), (137, 221, 23, 15, 'PB023E05', 100, '45W', 1, 1);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(138, 222, 23, 15, 'PB023F06', 88, '22.5W', 10, 1);

-- 水韵城 CAB-024 (12 slots, 6 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(225, 24, 'A01', 1, 139), (226, 24, 'B02', 1, 140), (227, 24, 'C03', 1, 141), (228, 24, 'D04', 1, 142), (229, 24, 'E05', 1, 143), (230, 24, 'F06', 1, 144), (231, 24, 'G07', 0, 0), (232, 24, 'H08', 0, 0);
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(233, 24, 'I09', 0, 0), (234, 24, 'J10', 0, 0), (235, 24, 'K11', 0, 0), (236, 24, 'L12', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(139, 225, 24, 16, 'PB024A01', 93, '20W', 5, 1), (140, 226, 24, 16, 'PB024B02', 76, '22.5W', 18, 1), (141, 227, 24, 16, 'PB024C03', 100, '20W', 2, 1), (142, 228, 24, 16, 'PB024D04', 55, '45W', 22, 1), (143, 229, 24, 16, 'PB024E05', 90, '22.5W', 7, 1);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(144, 230, 24, 16, 'PB024F06', 85, '20W', 11, 1);

-- 水韵城 CAB-025 (8 slots, 5 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(237, 25, 'A01', 1, 145), (238, 25, 'B02', 1, 146), (239, 25, 'C03', 1, 147), (240, 25, 'D04', 1, 148), (241, 25, 'E05', 1, 149), (242, 25, 'F06', 0, 0), (243, 25, 'G07', 0, 0), (244, 25, 'H08', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(145, 237, 25, 16, 'PB025A01', 73, '22.5W', 19, 1), (146, 238, 25, 16, 'PB025B02', 62, '20W', 23, 1), (147, 239, 25, 16, 'PB025C03', 97, '45W', 4, 1), (148, 240, 25, 16, 'PB025D04', 80, '22.5W', 14, 1), (149, 241, 25, 16, 'PB025E05', 94, '20W', 6, 1);

-- 金鹰购物中心 CAB-026 (12 slots, 7 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(245, 26, 'A01', 1, 150), (246, 26, 'B02', 1, 151), (247, 26, 'C03', 1, 152), (248, 26, 'D04', 1, 153), (249, 26, 'E05', 1, 154), (250, 26, 'F06', 1, 155), (251, 26, 'G07', 1, 156), (252, 26, 'H08', 0, 0);
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(253, 26, 'I09', 0, 0), (254, 26, 'J10', 0, 0), (255, 26, 'K11', 0, 0), (256, 26, 'L12', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(150, 245, 26, 17, 'PB026A01', 95, '22.5W', 12, 1), (151, 246, 26, 17, 'PB026B02', 87, '20W', 8, 1), (152, 247, 26, 17, 'PB026C03', 100, '45W', 3, 1), (153, 248, 26, 17, 'PB026D04', 72, '22.5W', 20, 1), (154, 249, 26, 17, 'PB026E05', 60, '20W', 15, 1);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(155, 250, 26, 17, 'PB026F06', 98, '22.5W', 6, 1), (156, 251, 26, 17, 'PB026G07', 45, '20W', 25, 1);

-- 金鹰购物中心 CAB-027 (8 slots, 4 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(257, 27, 'A01', 1, 157), (258, 27, 'B02', 1, 158), (259, 27, 'C03', 1, 159), (260, 27, 'D04', 1, 160), (261, 27, 'E05', 0, 0), (262, 27, 'F06', 0, 0), (263, 27, 'G07', 0, 0), (264, 27, 'H08', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(157, 257, 27, 17, 'PB027A01', 100, '45W', 1, 1), (158, 258, 27, 17, 'PB027B02', 88, '22.5W', 10, 1), (159, 259, 27, 17, 'PB027C03', 93, '20W', 5, 1), (160, 260, 27, 17, 'PB027D04', 76, '22.5W', 18, 1);

-- 宿迁市中医院 CAB-028 (8 slots, 5 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(265, 28, 'A01', 1, 161), (266, 28, 'B02', 1, 162), (267, 28, 'C03', 1, 163), (268, 28, 'D04', 1, 164), (269, 28, 'E05', 1, 165), (270, 28, 'F06', 0, 0), (271, 28, 'G07', 0, 0), (272, 28, 'H08', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(161, 265, 28, 18, 'PB028A01', 100, '20W', 2, 1), (162, 266, 28, 18, 'PB028B02', 55, '45W', 22, 1), (163, 267, 28, 18, 'PB028C03', 90, '22.5W', 7, 1), (164, 268, 28, 18, 'PB028D04', 85, '20W', 11, 1), (165, 269, 28, 18, 'PB028E05', 73, '22.5W', 19, 1);

-- 宿城区政府 CAB-029 (8 slots, 4 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(273, 29, 'A01', 1, 166), (274, 29, 'B02', 1, 167), (275, 29, 'C03', 1, 168), (276, 29, 'D04', 1, 169), (277, 29, 'E05', 0, 0), (278, 29, 'F06', 0, 0), (279, 29, 'G07', 0, 0), (280, 29, 'H08', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(166, 273, 29, 19, 'PB029A01', 62, '20W', 23, 1), (167, 274, 29, 19, 'PB029B02', 97, '45W', 4, 1), (168, 275, 29, 19, 'PB029C03', 80, '22.5W', 14, 1), (169, 276, 29, 19, 'PB029D04', 94, '20W', 6, 1);

-- 三台山森林公园 CAB-030 (8 slots, 6 occupied)
INSERT INTO cabinet_slots (id, cabinet_id, slot_no, status, power_bank_id) VALUES
(281, 30, 'A01', 1, 170), (282, 30, 'B02', 1, 171), (283, 30, 'C03', 1, 172), (284, 30, 'D04', 1, 173), (285, 30, 'E05', 1, 174), (286, 30, 'F06', 1, 175), (287, 30, 'G07', 0, 0), (288, 30, 'H08', 0, 0);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(170, 281, 30, 20, 'PB030A01', 95, '22.5W', 12, 1), (171, 282, 30, 20, 'PB030B02', 87, '20W', 8, 1), (172, 283, 30, 20, 'PB030C03', 100, '45W', 3, 1), (173, 284, 30, 20, 'PB030D04', 72, '22.5W', 20, 1), (174, 285, 30, 20, 'PB030E05', 60, '20W', 15, 1);
INSERT INTO power_banks (id, cabinet_slot_id, cabinet_id, station_id, power_bank_no, battery_level, power, total_charge_count, status) VALUES
(175, 286, 30, 20, 'PB030F06', 98, '22.5W', 6, 1);
