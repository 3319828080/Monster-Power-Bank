<template>
  <div class="home-page">
    <!-- Location header -->
    <van-sticky>
      <div class="location-bar" @click="showLocationPicker = true">
        <van-icon name="location-o" size="18" />
        <span class="location-text">{{ currentAddress || '获取位置中...' }}</span>
        <van-icon name="arrow-down" size="12" color="#969799" />
      </div>
    </van-sticky>

    <!-- Map area -->
    <div class="map-container" ref="mapRef">
      <div class="map-placeholder" v-if="!mapLoaded">
        <van-loading size="24px" /> 加载地图中...
      </div>
      <!-- Center pin -->
      <div class="center-pin" v-if="mapLoaded">
        <div class="pin-icon"></div>
        <div class="pin-shadow"></div>
      </div>
      <!-- Locate button -->
      <div class="locate-btn" v-if="mapLoaded" @click="locateUser" :class="{ spinning: locating }">
        <van-icon name="aim" size="20" />
      </div>
      <!-- Brand watermark -->
      <div class="monster-watermark" v-if="mapLoaded">怪兽充电宝</div>
    </div>

    <!-- Bottom panel -->
    <div class="bottom-panel">
      <div class="panel-header">
        <span class="panel-title">附近站点</span>
        <van-icon name="replay" @click="refreshStations" :class="{ spinning: loading }" />
      </div>

      <van-skeleton title :row="2" v-if="loading && stations.length === 0" />
      <van-empty description="附近暂无站点" v-else-if="!loading && stations.length === 0" />

      <div class="station-list" v-else>
        <div
          v-for="s in stations"
          :key="s.id"
          class="station-card"
          @click="$router.push(`/station/${s.id}`)"
        >
          <div class="station-left">
            <div class="station-name">{{ s.name }}</div>
            <div class="station-addr">{{ s.address }}</div>
          </div>
          <div class="station-right">
            <div class="station-avail">{{ s.availableBanks }}<span class="avail-unit">个可用</span></div>
            <div class="station-dist">{{ formatDistance(s.distance) }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Scan button -->
    <div class="scan-btn" @click="startScan">
      <van-icon name="scan" size="24" />
    </div>

    <!-- Current order floating button -->
    <div v-if="hasActiveOrder" class="float-btn" @click="$router.push('/order/current')">
      <van-button round type="primary" size="small">当前订单</van-button>
    </div>

    <!-- Scan overlay -->
    <van-overlay :show="scanning" @click="stopScan" />
    <div class="scan-dialog" v-if="scanning">
      <div class="scan-header">
        <span>扫描机柜二维码</span>
        <van-icon name="cross" size="20" @click="stopScan" />
      </div>
      <div class="scan-view">
        <div id="qr-scanner" class="scan-video" />
        <div class="scan-frame"></div>
      </div>
      <div class="scan-footer">
        <van-button block round @click="manualInput">手动输入机柜号</van-button>
      </div>
    </div>

    <!-- Location picker -->
    <van-popup v-model:show="showLocationPicker" round position="bottom" :style="{ height: '70vh' }">
      <div class="picker-content">
        <div class="picker-header">
          <span>选择位置</span>
          <van-icon name="cross" size="20" @click="showLocationPicker = false" />
        </div>
        <van-search v-model="searchKeyword" placeholder="搜索地址" @search="searchAddress" />
        <van-cell-group>
          <van-cell
            v-for="item in searchResults"
            :key="item.id"
            :title="item.name"
            :label="item.address"
            @click="selectLocation(item)"
          />
        </van-cell-group>
        <van-empty v-if="searchKeyword && searchResults.length === 0" description="未找到结果" />
        <div class="picker-current" v-if="!searchKeyword">
          <van-button block round type="primary" @click="locateUser; showLocationPicker = false">
            <van-icon name="aim" /> 重新定位当前位置
          </van-button>
        </div>
      </div>
    </van-popup>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showDialog } from 'vant'
import AMapLoader from '@amap/amap-jsapi-loader'
import { Html5Qrcode } from 'html5-qrcode'
import { listNearbyStations, getCurrentOrder, scanCabinet } from '@/api/order'
import type { StationInfo } from '@/api/order'
import { formatDistance } from '@/utils'

const router = useRouter()

const mapRef = ref<HTMLDivElement>()
const mapLoaded = ref(false)
const loading = ref(false)
const locating = ref(false)
const stations = ref<StationInfo[]>([])
const hasActiveOrder = ref(false)
const currentAddress = ref('宿迁职业技术学院')
const scanning = ref(false)
const showLocationPicker = ref(false)
const searchKeyword = ref('')
const searchResults = ref<any[]>([])

let map: any = null
let markers: any[] = []
let userMarker: any = null
let AMap: any = null
let userLat = 0
let userLng = 0
let html5Qrcode: Html5Qrcode | null = null

// Fallback station data — 2 campus stations + surrounding area
const fallbackStations: StationInfo[] = [
  { id: 1,  name: '图书馆站',     address: '宿迁职业技术学院图书馆一楼大厅',  latitude: 34.019900, longitude: 118.290500, distance: 0, availableBanks: 4 },
  { id: 2,  name: '第一食堂站',   address: '宿迁职业技术学院第一食堂门口',    latitude: 34.021200, longitude: 118.289300, distance: 0, availableBanks: 2 },
  { id: 9,  name: '万达广场站',   address: '宿迁市宿城区万达广场1号门入口',  latitude: 34.025000, longitude: 118.310000, distance: 0, availableBanks: 5 },
  { id: 10, name: '宝龙城市广场站', address: '宿迁市宿城区宝龙城市广场1楼大厅', latitude: 34.010000, longitude: 118.305000, distance: 0, availableBanks: 3 },
  { id: 11, name: '第一人民医院站', address: '宿迁市第一人民医院门诊大厅',     latitude: 34.033000, longitude: 118.298000, distance: 0, availableBanks: 4 },
  { id: 12, name: '宿迁学院站',   address: '宿迁学院图书馆一楼',               latitude: 34.002000, longitude: 118.270000, distance: 0, availableBanks: 7 },
  { id: 13, name: '汽车客运站',   address: '宿迁汽车客运站候车大厅',           latitude: 33.998000, longitude: 118.280000, distance: 0, availableBanks: 4 },
  { id: 14, name: '项王故里站',   address: '宿迁市宿城区项王故里景区入口',     latitude: 33.995000, longitude: 118.283000, distance: 0, availableBanks: 3 },
  { id: 15, name: '湖滨公园站',   address: '宿迁市湖滨新区湖滨公园游客中心',   latitude: 34.035000, longitude: 118.310000, distance: 0, availableBanks: 6 },
  { id: 16, name: '水韵城站',     address: '宿迁市宿城区水韵城购物中心B1层',   latitude: 34.008000, longitude: 118.298000, distance: 0, availableBanks: 4 },
  { id: 17, name: '金鹰购物中心站', address: '宿迁市宿城区金鹰国际购物中心1楼', latitude: 34.015000, longitude: 118.302000, distance: 0, availableBanks: 5 },
  { id: 18, name: '宿迁市中医院站', address: '宿迁市中医院门诊大厅一楼',        latitude: 34.028000, longitude: 118.272000, distance: 0, availableBanks: 3 },
  { id: 19, name: '宿城区政府站', address: '宿迁市宿城区政府大楼一楼',           latitude: 34.008000, longitude: 118.285000, distance: 0, availableBanks: 4 },
  { id: 20, name: '三台山森林公园站', address: '宿迁市三台山国家森林公园南门',  latitude: 34.050000, longitude: 118.310000, distance: 0, availableBanks: 2 },
]

const CENTER_LNG = 118.289660
const CENTER_LAT = 34.020516

function toRad(d: number) { return d * Math.PI / 180 }
function haversineKm(lat1: number, lng1: number, lat2: number, lng2: number) {
  const R = 6371
  const dLat = toRad(lat2 - lat1), dLng = toRad(lng2 - lng1)
  const a = Math.sin(dLat/2)**2 + Math.cos(toRad(lat1))*Math.cos(toRad(lat2))*Math.sin(dLng/2)**2
  return R * 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a))
}

async function loadMap() {
  try {
    (window as any)._AMapSecurityConfig = { securityJsCode: '在这里填入你的安全密钥' }
    AMap = await AMapLoader.load({
      key: '196b00dc9389d77ea035e0459a3f9925',
      version: '2.0',
      plugins: ['AMap.Geolocation', 'AMap.Geocoder'],
    })
    map = new AMap.Map(mapRef.value, {
      zoom: 15,
      center: [CENTER_LNG, CENTER_LAT],
    })
    mapLoaded.value = true
  } catch {
    showToast('地图加载失败')
  }
}

// Custom marker: "怪兽充电宝" branded bubble
function createMonsterMarker(s: StationInfo, isUser: boolean) {
  if (!map || !AMap) return null
  if (isUser) {
    return new AMap.Marker({
      position: [s.longitude, s.latitude],
      zIndex: 200,
      content: '<div class="user-pin"><div class="user-dot"></div><div class="user-ring"></div></div>',
      offset: new AMap.Pixel(-12, -12),
    })
  }
  const color = s.availableBanks > 0 ? '#07c160' : '#969799'
  const html =
    '<div class="monster-marker" style="position:relative;display:flex;flex-direction:column;align-items:center;cursor:pointer">' +
      '<div style="position:relative">' +
        '<div style="background:' + color + ';color:#fff;width:52px;height:28px;border-radius:14px;display:flex;align-items:center;justify-content:center;font-size:11px;font-weight:700;box-shadow:0 2px 8px rgba(7,193,96,.4);white-space:nowrap;padding:0 8px">' +
          '<span style="margin-right:2px">' + (s.availableBanks > 0 ? s.availableBanks + '个' : '满') + '</span>' +
        '</div>' +
        '<div style="width:0;height:0;border-left:6px solid transparent;border-right:6px solid transparent;border-top:6px solid ' + color + ';margin:0 auto"></div>' +
      '</div>' +
      '<div style="background:rgba(0,0,0,.65);color:#fff;font-size:10px;padding:2px 8px;border-radius:8px;margin-top:2px;white-space:nowrap;max-width:100px;overflow:hidden;text-overflow:ellipsis">' + s.name + '</div>' +
    '</div>'
  const marker = new AMap.Marker({
    position: [s.longitude, s.latitude],
    content: html,
    offset: new AMap.Pixel(-30, -48),
    zIndex: 100,
  })
  marker.on('click', () => router.push('/station/' + s.id))
  return marker
}

function addMarkers(list: StationInfo[]) {
  if (!map || !AMap) return
  markers.forEach((m: any) => map.remove(m))
  markers = []
  // User marker
  if (userMarker) map.add(userMarker)
  list.forEach((s) => {
    const m = createMonsterMarker(s, false)
    if (m) { map.add(m); markers.push(m) }
  })
}

async function loadNearbyStations() {
  if (!map) return
  const center = map.getCenter()
  if (!center) return
  const lat = center.getLat()
  const lng = center.getLng()

  try {
    const res = await listNearbyStations({ latitude: lat, longitude: lng, radius_meters: 50000 })
    if (res.list && res.list.length > 0) {
      stations.value = res.list
      addMarkers(res.list)
      return
    }
  } catch {
    console.log('[附近站点] API请求失败，使用本地数据')
  }

  // Fallback: use local data with distance calculation
  const sorted = fallbackStations.map((s) => ({
    ...s,
    distance: haversineKm(lat, lng, s.latitude, s.longitude) * 1000,
  })).sort((a, b) => a.distance - b.distance).filter((s) => s.distance <= 50000)
  stations.value = sorted
  addMarkers(sorted)
}

async function refreshStations() {
  loading.value = true
  await loadNearbyStations()
  loading.value = false
}

async function reverseGeocode(lat: number, lng: number) {
  if (!AMap?.Geocoder) return
  try {
    const geocoder = new AMap.Geocoder({})
    const result = await new Promise<any>((resolve, reject) => {
      geocoder.getAddress([lng, lat], (status: string, res: any) => {
        if (status === 'complete') resolve(res)
        else reject(new Error(status))
      })
    })
    if (result?.regeocode?.formattedAddress) {
      currentAddress.value = result.regeocode.formattedAddress
    }
  } catch { /* */ }
}

async function locateUser() {
  locating.value = true
  if (!mapLoaded.value) await loadMap()
  if (!map) { locating.value = false; return }
  map.resize?.()

  let lat: number | null = null
  let lng: number | null = null

  if (navigator.geolocation) {
    try {
      const pos = await new Promise<GeolocationPosition>((resolve, reject) => {
        navigator.geolocation.getCurrentPosition(resolve, reject, {
          enableHighAccuracy: true, timeout: 10000,
        })
      })
      lat = pos.coords.latitude
      lng = pos.coords.longitude
    } catch { /* */ }
  }

  if (lat == null && AMap?.Geolocation) {
    try {
      const geo = new AMap.Geolocation({ enableHighAccuracy: true, timeout: 10000 })
      const result = await new Promise<any>((resolve, reject) => {
        geo.getCurrentPosition((status: string, res: any) => {
          if (status === 'complete') resolve(res)
          else reject(new Error(status))
        })
      })
      if (result?.position) {
        lat = result.position.lat ?? result.position.getLat?.()
        lng = result.position.lng ?? result.position.getLng?.()
      }
    } catch { /* */ }
  }

  if (lat !== null && lng !== null) {
    userLat = lat; userLng = lng
    map.setCenter([lng, lat])
    reverseGeocode(lat, lng)
    userMarker = createMonsterMarker({ id: 0, name: '', address: '', latitude: lat, longitude: lng, distance: 0, availableBanks: 0 }, true)
    showToast('定位成功')
  } else {
    showToast('无法获取位置，显示默认区域')
  }

  locating.value = false
  await loadNearbyStations()
}

// ---- Scan ----
async function startScan() {
  scanning.value = true
  // Wait for DOM to render the scanner div
  await new Promise((r) => setTimeout(r, 100))

  try {
    html5Qrcode = new Html5Qrcode('qr-scanner')
    await html5Qrcode.start(
      { facingMode: 'environment' },
      { fps: 10, qrbox: { width: 200, height: 200 } },
      (decodedText: string) => {
        stopScan()
        const match = decodedText.match(/cabinet[=:](\d+)/i) || decodedText.match(/(\d+)/)
        if (match) {
          handleScanResult(Number(match[1]))
        } else {
          showToast('未识别到机柜号')
        }
      },
      () => {},
    )
  } catch (err: any) {
    const msg = typeof err === 'string' ? err : err?.message || ''
    if (msg.includes('NotAllowed') || msg.includes('Permission')) {
      showToast('请允许摄像头权限后重试')
    } else {
      showToast('摄像头启动失败，请手动输入')
    }
    stopScan()
  }
}

async function stopScan() {
  scanning.value = false
  try {
    if (html5Qrcode?.isScanning) {
      await html5Qrcode.stop()
    }
    html5Qrcode?.clear()
  } catch { /* */ }
  html5Qrcode = null
}

async function handleScanResult(cabinetId: number) {
  try {
    await scanCabinet(cabinetId)
    router.push('/cabinet/' + cabinetId)
  } catch (err: any) {
    showToast(err?.response?.data?.message || '机柜不存在')
  }
}

function manualInput() {
  stopScan()
  showDialog({
    title: '输入机柜编号',
    message: '请输入机柜ID或编号',
    confirmButtonText: '确认',
  }).then((r: any) => {
    if (r?.value) {
      const id = parseInt(r.value, 10)
      if (id > 0) handleScanResult(id)
    }
  })
}

function searchAddress() { showToast('搜索功能开发中') }
function selectLocation(_item: any) { showLocationPicker.value = false }

onMounted(async () => {
  await loadMap()
  try {
    const res = await getCurrentOrder()
    hasActiveOrder.value = !!res.order
  } catch { /* */ }
  await loadNearbyStations()
})

onUnmounted(() => {
  markers.forEach((m: any) => map?.remove(m))
  map?.destroy()
  stopScan()
})
</script>

<style scoped>
.home-page { height: 100vh; display: flex; flex-direction: column; background: #f7f8fa; }

/* Location bar */
.location-bar {
  display: flex; align-items: center; gap: 6px;
  padding: 10px 16px; background: #fff; cursor: pointer;
  border-bottom: 1px solid #f5f5f5;
}
.location-text { flex: 1; font-size: 14px; font-weight: 500; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

/* Map */
.map-container { flex: 1; min-height: 300px; position: relative; }
.map-placeholder { height: 100%; display: flex; align-items: center; justify-content: center; background: #f7f8fa; color: #969799; gap: 8px; }

/* Center pin */
.center-pin { position: absolute; top: 50%; left: 50%; z-index: 10; pointer-events: none; transform: translate(-50%, -100%); }
.pin-icon { width: 32px; height: 32px; background: #ff4444; border-radius: 50% 50% 50% 0; transform: rotate(-45deg); border: 3px solid #fff; box-shadow: 0 2px 6px rgba(0,0,0,.3); }
.pin-shadow { width: 8px; height: 8px; background: rgba(0,0,0,.25); border-radius: 50%; margin: 4px auto 0; }

/* Locate button */
.locate-btn {
  position: absolute; bottom: 20px; right: 16px; z-index: 10;
  width: 40px; height: 40px; border-radius: 50%;
  background: #fff; box-shadow: 0 2px 8px rgba(0,0,0,.15);
  display: flex; align-items: center; justify-content: center;
  cursor: pointer; color: #1989fa;
}
.locate-btn.spinning { color: #969799; }
.locate-btn:active { background: #f5f5f5; }

/* Bottom panel */
.bottom-panel { max-height: 42vh; overflow-y: auto; background: #fff; border-radius: 12px 12px 0 0; position: relative; z-index: 1; }
.panel-header { display: flex; justify-content: space-between; align-items: center; padding: 16px; border-bottom: 1px solid #f5f5f5; position: sticky; top: 0; background: #fff; }
.panel-title { font-size: 16px; font-weight: 600; }
.station-list { padding: 8px 16px 24px; }
.station-card { display: flex; align-items: center; padding: 14px 0; border-bottom: 1px solid #f5f5f5; cursor: pointer; }
.station-card:last-child { border-bottom: none; }
.station-left { flex: 1; }
.station-name { font-size: 15px; color: #323233; font-weight: 500; margin-bottom: 4px; }
.station-addr { font-size: 12px; color: #969799; }
.station-right { text-align: right; }
.station-avail { font-size: 15px; font-weight: 600; color: #07c160; }
.avail-unit { font-size: 11px; font-weight: 400; color: #969799; }
.station-dist { font-size: 12px; color: #969799; margin-top: 2px; }

/* Scan button */
.scan-btn {
  position: fixed; bottom: 100px; right: 16px; z-index: 100;
  width: 48px; height: 48px; border-radius: 50%;
  background: linear-gradient(135deg, #1989fa, #07c160);
  color: #fff; box-shadow: 0 4px 12px rgba(7,193,96,.35);
  display: flex; align-items: center; justify-content: center;
  cursor: pointer;
}
.scan-btn:active { opacity: 0.9; }

/* Float button */
.float-btn { position: fixed; bottom: 80px; left: 16px; right: 16px; z-index: 100; }

/* Scan dialog */
.scan-dialog { position: fixed; top: 0; left: 0; right: 0; bottom: 0; z-index: 200; display: flex; flex-direction: column; }
.scan-header { display: flex; justify-content: space-between; align-items: center; padding: 16px; background: #000; color: #fff; font-size: 16px; }
.scan-view { flex: 1; position: relative; background: #000; overflow: hidden; }
.scan-video { width: 100%; height: 100%; object-fit: cover; }
.scan-frame { position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); width: 200px; height: 200px; border: 2px solid #07c160; border-radius: 12px; box-shadow: 0 0 0 2000px rgba(0,0,0,.5); }
.scan-footer { padding: 16px; background: #000; }
.scan-footer .van-button { --van-button-default-background: #333; --van-button-default-color: #fff; }

/* Location picker */
.picker-content { padding: 16px; }
.picker-header { display: flex; justify-content: space-between; align-items: center; font-size: 16px; font-weight: 600; margin-bottom: 12px; }
.picker-current { padding: 24px 0; }

/* Animations */
.spinning { animation: spin 1s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>

<!-- AMap custom marker styles (must be unscoped for AMap DOM) -->
<style>
.monster-marker { pointer-events: auto; }
.monster-marker:hover > div:first-child > div:first-child {
  transform: scale(1.1); transition: transform .15s;
}
.user-pin { position: relative; }
.user-dot {
  width: 14px; height: 14px; background: #1989fa;
  border: 3px solid #fff; border-radius: 50%;
  box-shadow: 0 2px 8px rgba(25,137,250,.5);
}
.user-ring {
  position: absolute; top: -6px; left: -6px;
  width: 26px; height: 26px; border: 2px solid rgba(25,137,250,.4);
  border-radius: 50%; animation: pulse-ring 2s ease-out infinite;
}
@keyframes pulse-ring {
  0% { transform: scale(.8); opacity: 1; }
  100% { transform: scale(2.5); opacity: 0; }
}

/* 怪兽充电宝 品牌水印 */
.monster-watermark {
  position: absolute; bottom: 8px; left: 50%; z-index: 20;
  transform: translateX(-50%); pointer-events: none;
  background: rgba(7,193,96,.9); color: #fff;
  padding: 3px 12px; border-radius: 10px;
  font-size: 11px; font-weight: 600; white-space: nowrap;
}
</style>
