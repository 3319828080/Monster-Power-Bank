// Generate seed SQL matching GORM AutoMigrate table structure
// Run: node scripts/gen_seed_sql.js > migrations/000004_seed_full.sql

const stations = [
  { id: 1,  name: '图书馆站',     address: '宿迁职业技术学院图书馆一楼大厅',   lat: 34.019900, lng: 118.290500, open_time: '07:00-22:00', desc: '紧邻图书馆借阅大厅，安静环境适合边学习边充电' },
  { id: 2,  name: '第一食堂站',   address: '宿迁职业技术学院第一食堂门口',     lat: 34.021200, lng: 118.289300, open_time: '06:30-21:30', desc: '食堂入口右侧，用餐高峰期人流量大，顺路借还方便' },
  { id: 9,  name: '万达广场站',   address: '宿迁市宿城区万达广场1号门入口',    lat: 34.025000, lng: 118.310000, open_time: '09:00-22:00', desc: '万达广场1楼服务台旁，宿迁核心商圈，逛街购物随时借还' },
  { id: 10, name: '宝龙城市广场站', address: '宿迁市宿城区宝龙城市广场1楼大厅', lat: 34.010000, lng: 118.305000, open_time: '09:00-22:00', desc: '宝龙广场中庭扶梯旁，商业核心区域，客流量大' },
  { id: 11, name: '第一人民医院站', address: '宿迁市第一人民医院门诊大厅',      lat: 34.033000, lng: 118.298000, open_time: '00:00-24:00', desc: '门诊大厅一楼导诊台对面，24小时全天候服务，满足医患应急充电需求' },
  { id: 12, name: '宿迁学院站',     address: '宿迁学院图书馆一楼',              lat: 34.002000, lng: 118.270000, open_time: '07:00-22:00', desc: '宿迁学院校园核心区域，服务广大师生日常充电需求' },
  { id: 13, name: '汽车客运站',     address: '宿迁汽车客运站候车大厅',          lat: 33.998000, lng: 118.280000, open_time: '06:00-20:00', desc: '候车大厅安检口左侧，出行途中手机没电不用愁' },
  { id: 14, name: '项王故里站',     address: '宿迁市宿城区项王故里景区入口',    lat: 33.995000, lng: 118.283000, open_time: '08:00-18:00', desc: '景区游客中心大厅，游览西楚霸王故里，随时保持电量满满' },
  { id: 15, name: '湖滨公园站',     address: '宿迁市湖滨新区湖滨公园游客中心',  lat: 34.035000, lng: 118.310000, open_time: '08:00-18:00', desc: '游客中心进门处，骆马湖畔休闲漫步也能随时充电' },
  { id: 16, name: '水韵城站',       address: '宿迁市宿城区水韵城购物中心B1层',  lat: 34.008000, lng: 118.298000, open_time: '09:00-22:00', desc: 'B1美食广场入口，餐饮娱乐一站式，等餐时即可充满电' },
  { id: 17, name: '金鹰购物中心站', address: '宿迁市宿城区金鹰国际购物中心1楼', lat: 34.015000, lng: 118.302000, open_time: '09:00-22:00', desc: '1楼化妆品区旁，高端商圈配套，购物体验更贴心' },
  { id: 18, name: '宿迁市中医院站', address: '宿迁市中医院门诊大厅一楼',        lat: 34.028000, lng: 118.272000, open_time: '00:00-24:00', desc: '中医门诊大厅左侧，24小时全天候守护健康与电量' },
  { id: 19, name: '宿城区政府站',   address: '宿迁市宿城区政府大楼一楼',        lat: 34.008000, lng: 118.285000, open_time: '08:00-18:00', desc: '政府大楼一楼大厅电梯旁，政务办事等待时免费充电' },
  { id: 20, name: '三台山森林公园站', address: '宿迁市三台山国家森林公园南门',  lat: 34.050000, lng: 118.310000, open_time: '08:00-17:00', desc: '景区南门游客服务中心，花海森林与满电同行，旅途更精彩' },
];

// [cabinet_id, station_id, cabinet_no, total_slots, occupied]
const cabinets = [
  [1,  1,  'CAB-001-A', 12, 8], [2,  1,  'CAB-001-B', 8,  5],
  [3,  2,  'CAB-002-A', 12, 6], [4,  2,  'CAB-002-B', 8,  4],
  [12, 9,  'CAB-009-A', 12, 7], [13, 9,  'CAB-009-B', 8,  4],
  [14, 10, 'CAB-010-A', 12, 5], [15, 10, 'CAB-010-B', 8,  5],
  [16, 11, 'CAB-011-A', 12, 6], [17, 11, 'CAB-011-B', 8,  6],
  [18, 12, 'CAB-012-A', 12, 8], [19, 12, 'CAB-012-B', 8,  5],
  [20, 13, 'CAB-013-A', 8,  4],
  [21, 14, 'CAB-014-A', 8,  5],
  [22, 15, 'CAB-015-A', 12, 8], [23, 15, 'CAB-015-B', 8,  6],
  [24, 16, 'CAB-016-A', 12, 6], [25, 16, 'CAB-016-B', 8,  5],
  [26, 17, 'CAB-017-A', 12, 7], [27, 17, 'CAB-017-B', 8,  4],
  [28, 18, 'CAB-018-A', 8,  5],
  [29, 19, 'CAB-019-A', 8,  4],
  [30, 20, 'CAB-020-A', 8,  6],
];

const batteries = [95, 87, 100, 72, 60, 98, 45, 100, 88, 93, 76, 100, 55, 90, 85, 73, 62, 97, 80, 94];
const letters = 'ABCDEFGHIJKL';

function pad2(n) { return String(n).padStart(2, '0'); }
function pad3(n) { return String(n).padStart(3, '0'); }

let sql = '';

// Clean existing seed data
sql += '-- 清理已有种子数据\n';
sql += 'DELETE FROM power_banks WHERE id >= 1;\n';
sql += 'DELETE FROM slots WHERE id >= 1;\n';
sql += 'DELETE FROM cabinets WHERE id >= 1;\n';
sql += 'DELETE FROM stations WHERE id >= 1;\n\n';

// Stations
sql += '-- 站点 (14个: 2校园 + 12周边)\n';
for (const s of stations) {
  sql += `INSERT INTO stations (id, name, address, latitude, longitude, status, open_time, description, images) VALUES
(${s.id}, '${s.name}', '${s.address}', ${s.lat}, ${s.lng}, 1, '${s.open_time}', '${s.desc}', '["https://picsum.photos/seed/station${s.id}/400/300"]');
`;
}
sql += '\n';

// Cabinets
sql += '-- 机柜 (22个)\n';
for (const c of cabinets) {
  sql += `INSERT INTO cabinets (id, station_id, cabinet_no, status) VALUES (${c[0]}, ${c[1]}, '${c[2]}', 1);
`;
}
sql += '\n';

// Slots and power banks
sql += '-- 仓位 & 充电宝\n';
let slotID = 1;
let pbID = 1;
for (const c of cabinets) {
  const [, stationID, , total, occupied] = c;
  for (let i = 0; i < total; i++) {
    const slotNo = letters[i] + pad2(i + 1);
    if (i < occupied) {
      const deviceNo = 'PB' + pad3(c[0]) + letters[i] + pad2(i + 1);
      const battery = batteries[(pbID - 1) % batteries.length];
      sql += `INSERT INTO slots (id, cabinet_id, station_id, slot_no, status, power_bank_id) VALUES (${slotID}, ${c[0]}, ${stationID}, '${slotNo}', 'occupied', ${pbID});
INSERT INTO power_banks (id, device_no, station_id, cabinet_id, slot_id, status, battery_level) VALUES (${pbID}, '${deviceNo}', ${stationID}, ${c[0]}, ${slotID}, 'available', ${battery});
`;
      pbID++;
    } else {
      sql += `INSERT INTO slots (id, cabinet_id, station_id, slot_no, status, power_bank_id) VALUES (${slotID}, ${c[0]}, ${stationID}, '${slotNo}', 'empty', 0);
`;
    }
    slotID++;
  }
}

console.log(sql);
