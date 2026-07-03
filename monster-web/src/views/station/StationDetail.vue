<template>
  <div class="station-page">
    <van-nav-bar :title="station?.name || '站点详情'" left-arrow @click-left="$router.back()" />

    <van-loading v-if="!station" class="page-loading" size="24px" />

    <template v-else>
      <!-- Station images -->
      <van-swipe :autoplay="3000" indicator-color="#07c160" class="station-swipe" v-if="imageList.length > 0">
        <van-swipe-item v-for="(img, idx) in imageList" :key="idx">
          <van-image :src="img" width="100%" height="200" fit="cover" />
        </van-swipe-item>
      </van-swipe>

      <!-- Station info -->
      <div class="station-header">
        <div class="station-name">{{ station.name }}</div>
        <div class="station-addr" @click="openNavigation">
          <van-icon name="location-o" />
          {{ station.address }}
          <van-icon name="arrow" class="nav-arrow" />
        </div>
        <div class="station-meta">
          <van-tag type="success">{{ station.openTime || '暂无' }}</van-tag>
          <span class="meta-divider">|</span>
          <span class="meta-text">{{ powerBanks.length }}个可借</span>
        </div>
        <div class="station-desc" v-if="station.description">
          <van-icon name="info-o" />
          {{ station.description }}
        </div>
      </div>

      <!-- Map preview -->
      <div class="mini-map" ref="miniMapRef" @click="openNavigation">
        <div class="map-placeholder" v-if="!miniMapLoaded">地图加载中...</div>
        <div class="map-nav-tip" v-else>点击地图导航到此处</div>
      </div>

      <!-- Pricing info -->
      <div class="info-card">
        <div class="card-title">计费规则</div>
        <div class="pricing-row">
          <span class="pricing-label">起步价</span>
          <span class="pricing-value">¥2.00 / 前30分钟</span>
        </div>
        <div class="pricing-row">
          <span class="pricing-label">每小时</span>
          <span class="pricing-value">¥2.00 / 小时</span>
        </div>
        <div class="pricing-row">
          <span class="pricing-label">每日封顶</span>
          <span class="pricing-value">¥20.00</span>
        </div>
        <div class="pricing-tip">充电宝归还后根据实际使用时长计费</div>
      </div>

      <!-- Power bank list -->
      <div class="pb-section">
        <div class="section-title">可选充电宝 ({{ powerBanks.length }})</div>

        <div class="pb-grid" v-if="powerBanks.length > 0">
          <div
            v-for="item in powerBanks"
            :key="item.slot.id"
            class="pb-card"
            :class="{ selected: selectedItem?.slot.id === item.slot.id }"
            @click="selectedItem = item"
          >
            <div class="pb-battery">
              <van-circle
                :rate="item.slot.batteryLevel"
                :text="item.slot.batteryLevel + '%'"
                :size="60"
                :stroke-width="6"
                :color="batteryColor(item.slot.batteryLevel)"
                layer-color="#ebedf0"
              />
            </div>
            <div class="pb-info">
              <div class="pb-no">{{ item.slot.powerBankNo }}</div>
              <div class="pb-meta">
                <span class="pb-cabinet">{{ item.cabinet.cabinetNo }}</span>
                <span class="pb-divider">·</span>
                <span class="pb-slot">{{ item.slot.slotNo }}</span>
              </div>
              <div class="pb-power">{{ item.slot.power || '22.5W' }}</div>
            </div>
            <van-icon v-if="selectedItem?.slot.id === item.slot.id" name="success" color="#07c160" size="20" class="pb-check" />
          </div>
        </div>
        <van-empty v-else description="暂无可用充电宝" />
      </div>

      <!-- Borrow button -->
      <div class="bottom-action" v-if="selectedItem">
        <div class="selected-info">
          <div class="selected-pb-no">{{ selectedItem.slot.powerBankNo }}</div>
          <div class="selected-detail">
            电量 {{ selectedItem.slot.batteryLevel }}% · {{ selectedItem.cabinet.cabinetNo }} {{ selectedItem.slot.slotNo }}
          </div>
        </div>
        <van-button round block type="primary" @click="handleBorrow" :loading="borrowing">
          立即借用
        </van-button>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showConfirmDialog, showToast } from 'vant'
import AMapLoader from '@amap/amap-jsapi-loader'
import { getStation, getCabinetList, preCheckBorrow, createOrder } from '@/api/order'
import type { StationInfo, CabinetInfo, SlotInfo } from '@/api/order'

const route = useRoute()
const router = useRouter()
const stationId = Number(route.params.id)

const station = ref<StationInfo | null>(null)
const cabinets = ref<CabinetInfo[]>([])
const selectedItem = ref<{ cabinet: CabinetInfo; slot: SlotInfo } | null>(null)
const borrowing = ref(false)
const miniMapRef = ref<HTMLDivElement>()
const miniMapLoaded = ref(false)
let miniMap: any = null

const imageList = computed(() => {
  if (!station.value?.images) return []
  try {
    return JSON.parse(station.value.images)
  } catch {
    return [station.value.images]
  }
})

const powerBanks = computed(() => {
  const result: { cabinet: CabinetInfo; slot: SlotInfo }[] = []
  for (const cab of cabinets.value) {
    for (const sl of cab.slots) {
      if (sl.status === '空闲') {
        result.push({ cabinet: cab, slot: sl })
      }
    }
  }
  return result
})

function batteryColor(level: number) {
  if (level >= 60) return '#07c160'
  if (level >= 30) return '#ff976a'
  return '#ee0a24'
}

function openNavigation() {
  if (!station.value) return
  const { longitude: lng, latitude: lat, name } = station.value
  window.open(`https://uri.amap.com/navigation?to=${lng},${lat},${encodeURIComponent(name)}&mode=walk&callnative=1`, '_blank')
}

async function handleBorrow() {
  if (!selectedItem.value) return
  const { cabinet, slot } = selectedItem.value
  try {
    await showConfirmDialog({ title: '确认借用', message: `借用 ${slot.powerBankNo}？\n将冻结押金，归还后按实际时长计费` })
  } catch {
    return
  }

  borrowing.value = true
  try {
    const preCheck = await preCheckBorrow({
      powerBankId: slot.powerBankId,
      stationId: stationId,
      cabinetId: cabinet.id,
      slotId: slot.id,
    })
    if (!preCheck.allowed) {
      const reason = preCheck.reason || '暂不可借用'
      if (reason.includes('余额不足') || reason.includes('押金')) {
        try {
          await showConfirmDialog({
            title: '余额不足',
            message: `${reason}\n\n请先充值押金后重试。\n押金: ¥${(preCheck.depositRequired / 100).toFixed(2)}`,
            confirmButtonText: '去充值',
          })
          router.push('/profile')
          return
        } catch { return }
      }
      showToast(reason)
      return
    }

    const order = await createOrder({
      powerBankId: slot.powerBankId,
      stationId: stationId,
      cabinetId: cabinet.id,
      slotId: slot.id,
    })
    showToast('借用成功')
    await showConfirmDialog({
      title: '租借成功',
      message: `订单号：${order.orderNo}\n\n起步价 ¥${(order.startFee / 100).toFixed(2)} / 前60分钟\n每小时 ¥${(order.hourlyFee / 100).toFixed(2)}\n每日封顶 ¥${(order.dailyCap / 100).toFixed(2)}\n\n已冻结押金 ¥${(order.depositFrozen / 100).toFixed(2)}\n\n归还时按实际使用时长计费`,
      confirmButtonText: '查看订单',
      showCancelButton: true,
      cancelButtonText: '关闭',
    })
    router.replace(`/order/current`)
  } catch {
    // toast handled by interceptor
  } finally {
    borrowing.value = false
  }
}

async function loadStation() {
  try {
    const res = await getStation(stationId)
    station.value = res.station
  } catch { /* */ }
}

async function loadMiniMap() {
  if (!station.value) return
  try {
    const AMap = await AMapLoader.load({
      key: '196b00dc9389d77ea035e0459a3f9925',
      version: '2.0',
    })
    miniMap = new AMap.Map(miniMapRef.value, {
      zoom: 16,
      center: [station.value.longitude, station.value.latitude],
      mapStyle: 'amap://styles/whitesmoke',
    })
    const marker = new AMap.Marker({
      position: [station.value.longitude, station.value.latitude],
      title: station.value.name,
    })
    miniMap.add(marker)
    miniMapLoaded.value = true
  } catch { /* */ }
}

onMounted(async () => {
  await loadStation()
  if (station.value) {
    loadMiniMap()
    try {
      const res = await getCabinetList(stationId)
      cabinets.value = res.cabinets || []
    } catch { /* */ }
  }
})

onUnmounted(() => { miniMap?.destroy() })
</script>

<style scoped>
.station-page { min-height: 100vh; background: #f7f8fa; padding-bottom: 80px; }
.page-loading { display: flex; justify-content: center; padding: 80px 0; }
.station-swipe { margin-bottom: 0; }
.station-header { background: #fff; padding: 16px; }
.station-name { font-size: 18px; font-weight: 600; margin-bottom: 6px; }
.station-addr { font-size: 13px; color: #969799; display: flex; align-items: center; gap: 4px; margin-bottom: 8px; cursor: pointer; }
.station-meta { display: flex; align-items: center; gap: 8px; margin-bottom: 8px; }
.meta-divider { color: #ebedf0; }
.meta-text { font-size: 13px; color: #07c160; font-weight: 500; }
.station-desc { font-size: 12px; color: #646566; background: #f7f8fa; padding: 8px 10px; border-radius: 6px; display: flex; align-items: flex-start; gap: 4px; line-height: 1.5; }
.station-addr .nav-arrow { margin-left: auto; color: #07c160; font-size: 14px; }

.mini-map { height: 150px; margin: 8px 16px; border-radius: 8px; overflow: hidden; position: relative; }
.map-placeholder { height: 100%; display: flex; align-items: center; justify-content: center; background: #f0f0f0; color: #969799; font-size: 13px; }
.map-nav-tip { position: absolute; bottom: 6px; right: 8px; background: rgba(0,0,0,0.55); color: #fff; font-size: 11px; padding: 2px 8px; border-radius: 10px; pointer-events: none; }

.info-card { background: #fff; margin: 8px 16px; border-radius: 8px; padding: 14px 16px; }
.card-title { font-size: 15px; font-weight: 600; margin-bottom: 10px; }
.pricing-row { display: flex; justify-content: space-between; padding: 6px 0; font-size: 13px; }
.pricing-label { color: #646566; }
.pricing-value { color: #323233; font-weight: 500; }
.pricing-tip { font-size: 11px; color: #c8c9cc; margin-top: 8px; }

/* Power bank grid */
.pb-section { padding: 0 16px; margin-top: 12px; }
.section-title { font-size: 16px; font-weight: 600; margin-bottom: 10px; }
.pb-grid { display: flex; flex-direction: column; gap: 8px; }
.pb-card {
  background: #fff; border-radius: 10px; padding: 14px;
  display: flex; align-items: center; gap: 14px;
  border: 2px solid transparent; cursor: pointer; position: relative;
}
.pb-card.selected { border-color: #07c160; background: #f0faf4; }
.pb-battery { flex-shrink: 0; }
.pb-info { flex: 1; min-width: 0; }
.pb-no { font-size: 15px; font-weight: 600; color: #323233; margin-bottom: 4px; }
.pb-meta { font-size: 12px; color: #969799; margin-bottom: 2px; }
.pb-divider { margin: 0 4px; }
.pb-power { font-size: 11px; color: #07c160; font-weight: 500; }
.pb-check { flex-shrink: 0; }

.selected-info { margin-bottom: 12px; text-align: center; }
.selected-pb-no { font-size: 16px; font-weight: 600; color: #323233; margin-bottom: 4px; }
.selected-detail { font-size: 12px; color: #969799; }

.bottom-action {
  position: fixed; bottom: 0; left: 0; right: 0;
  background: #fff; padding: 12px 16px;
  padding-bottom: calc(12px + env(safe-area-inset-bottom));
  box-shadow: 0 -2px 8px rgba(0,0,0,0.08);
}
</style>
