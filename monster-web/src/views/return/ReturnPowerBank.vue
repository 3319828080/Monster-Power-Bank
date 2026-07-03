<template>
  <div class="return-page">
    <van-nav-bar title="归还充电宝" left-arrow @click-left="$router.back()" />

    <div class="section">
      <div class="section-title">输入机柜编号查找</div>
      <form class="search-row" @submit.prevent="handleSearch">
        <van-field v-model="cabinetNo" placeholder="输入机柜编号，如 CAB-001" :disabled="searching" />
        <van-button type="primary" size="small" :loading="searching" native-type="submit">查找</van-button>
      </form>
      <div v-if="searchResult" class="search-result">
        <div class="cabinet-item selected" style="margin-top:0">
          <div class="cabinet-left">
            <div class="cabinet-no">{{ searchResult.cabinetNo }}</div>
            <div class="cabinet-station">{{ searchResult.stationName }}</div>
            <div class="cabinet-addr">{{ searchResult.stationAddress }}</div>
          </div>
          <div class="cabinet-right">
            <div class="cabinet-slots">
              <span class="slots-empty">{{ searchResult.emptySlotCount }}</span>
              <span class="slots-sep">/</span>
              <span>{{ searchResult.totalSlots }}</span>
              <span class="slots-label">空位</span>
            </div>
          </div>
          <van-icon name="success" color="#07c160" size="20" class="cabinet-check" />
        </div>
        <van-button
          round block type="primary"
          :loading="returningBySearch"
          @click="handleSearchReturn"
          style="margin-top:10px"
        >
          归还至该机柜
        </van-button>
      </div>
      <div v-if="searchNotFound" class="search-404">未找到该机柜或机柜已满</div>
    </div>

    <div class="section">
      <div class="section-title">选择附近可归还机柜</div>

      <van-loading v-if="loading" class="page-loading" size="24px" />

      <template v-else>
        <div
          v-for="item in cabinets"
          :key="item.cabinetId"
          class="cabinet-item"
          :class="{ selected: selectedCabinet?.cabinetId === item.cabinetId }"
          @click="selectedCabinet = item"
        >
          <div class="cabinet-left">
            <div class="cabinet-no">{{ item.cabinetNo }}</div>
            <div class="cabinet-station">{{ item.stationName }}</div>
            <div class="cabinet-addr">{{ item.stationAddress }}</div>
          </div>
          <div class="cabinet-right">
            <div class="cabinet-dist">{{ formatDist(item.distance) }}</div>
            <div class="cabinet-slots">
              <span class="slots-empty">{{ item.emptySlotCount }}</span>
              <span class="slots-sep">/</span>
              <span>{{ item.totalSlots }}</span>
              <span class="slots-label">空位</span>
            </div>
          </div>
          <van-icon v-if="selectedCabinet?.cabinetId === item.cabinetId" name="success" color="#07c160" size="20" class="cabinet-check" />
        </div>

        <van-empty v-if="cabinets.length === 0" description="附近暂无空闲机柜">
          <van-button round type="primary" size="small" @click="fetchCabinets">重新搜索</van-button>
        </van-empty>
      </template>
    </div>

    <div class="bottom-action" v-if="selectedCabinet">
      <div class="selected-info">
        归还至 {{ selectedCabinet.cabinetNo }}（{{ selectedCabinet.stationName }}）
      </div>
      <van-button
        round
        block
        type="primary"
        :loading="returning"
        @click="handleReturn"
      >
        立即归还
      </van-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showConfirmDialog, showToast } from 'vant'
import { listReturnCabinets, confirmReturn, searchCabinet } from '@/api/order'
import type { ReturnCabinetInfo } from '@/api/order'

const route = useRoute()
const router = useRouter()
const orderNo = route.params.orderNo as string

const cabinets = ref<ReturnCabinetInfo[]>([])
const selectedCabinet = ref<ReturnCabinetInfo | null>(null)
const loading = ref(true)
const returning = ref(false)
const cabinetNo = ref('')
const searching = ref(false)
const searchResult = ref<ReturnCabinetInfo | null>(null)
const searchNotFound = ref(false)
const returningBySearch = ref(false)

function formatDist(m: number) {
  if (m < 1000) return Math.round(m) + 'm'
  return (m / 1000).toFixed(1) + 'km'
}

async function fetchCabinets() {
  loading.value = true
  let lat = 34.020516
  let lng = 118.289660
  try {
    const pos = await new Promise<GeolocationPosition>((resolve, reject) => {
      navigator.geolocation.getCurrentPosition(resolve, reject, { timeout: 5000 })
    })
    lat = pos.coords.latitude
    lng = pos.coords.longitude
  } catch { /* use default */ }

  try {
    const res = await listReturnCabinets({
      latitude: lat,
      longitude: lng,
      radiusMeters: 50000,
    })
    cabinets.value = res.list || []
  } catch { /* */ } finally {
    loading.value = false
  }
}

async function handleReturn() {
  if (!selectedCabinet.value) return
  const c = selectedCabinet.value
  try {
    await showConfirmDialog({
      title: '确认归还',
      message: `将充电宝归还至 ${c.cabinetNo}？\n${c.stationName}\n\n归还后按实际使用时长自动结算`,
    })
  } catch {
    return
  }

  returning.value = true
  try {
    const res = await confirmReturn({
      orderNo,
      cabinetId: c.cabinetId,
    })
    showToast(`归还成功，费用 ¥${(res.paidAmount / 100).toFixed(2)}`)
    router.replace('/orders')
  } catch {
    // handled by interceptor
  } finally {
    returning.value = false
  }
}

async function handleSearch() {
  const no = cabinetNo.value.trim()
  if (!no) return
  searchNotFound.value = false
  searchResult.value = null
  searching.value = true
  try {
    const res = await searchCabinet(no)
    if (res.cabinet && res.cabinet.emptySlotCount > 0) {
      searchResult.value = res.cabinet
    } else {
      searchNotFound.value = true
    }
  } catch {
    searchNotFound.value = true
  } finally {
    searching.value = false
  }
}

async function handleSearchReturn() {
  if (!searchResult.value) return
  const c = searchResult.value
  try {
    await showConfirmDialog({
      title: '确认归还',
      message: `将充电宝归还至 ${c.cabinetNo}？\n${c.stationName}\n\n归还后按实际使用时长自动结算`,
    })
  } catch {
    return
  }
  returningBySearch.value = true
  try {
    const res = await confirmReturn({
      orderNo,
      cabinetId: c.cabinetId,
    })
    showToast(`归还成功，费用 ¥${(res.paidAmount / 100).toFixed(2)}`)
    router.replace('/orders')
  } catch {
    // handled by interceptor
  } finally {
    returningBySearch.value = false
  }
}

onMounted(() => {
  fetchCabinets()
})
</script>

<style scoped>
.return-page { min-height: 100vh; background: #f7f8fa; padding-bottom: 80px; }
.page-loading { display: flex; justify-content: center; padding: 80px 0; }
.section { padding: 12px 16px; }
.section-title { font-size: 15px; font-weight: 600; margin-bottom: 10px; }

.cabinet-item {
  background: #fff; border-radius: 10px; padding: 14px;
  display: flex; align-items: center; gap: 12px;
  border: 2px solid transparent; cursor: pointer;
  margin-bottom: 8px; position: relative;
}
.cabinet-item.selected { border-color: #07c160; background: #f0faf4; }
.cabinet-left { flex: 1; min-width: 0; }
.cabinet-no { font-size: 15px; font-weight: 600; color: #323233; margin-bottom: 4px; }
.cabinet-station { font-size: 13px; color: #646566; margin-bottom: 2px; }
.cabinet-addr { font-size: 11px; color: #969799; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.cabinet-right { text-align: right; flex-shrink: 0; }
.cabinet-dist { font-size: 13px; color: #07c160; font-weight: 500; margin-bottom: 4px; }
.cabinet-slots { font-size: 12px; color: #969799; }
.slots-empty { color: #07c160; font-weight: 600; font-size: 14px; }
.slots-sep { margin: 0 1px; }
.slots-label { margin-left: 2px; font-size: 11px; }
.cabinet-check { flex-shrink: 0; }

.bottom-action {
  position: fixed; bottom: 0; left: 0; right: 0;
  background: #fff; padding: 12px 16px;
  padding-bottom: calc(12px + env(safe-area-inset-bottom));
  box-shadow: 0 -2px 8px rgba(0,0,0,0.08);
}
.selected-info { margin-bottom: 12px; text-align: center; font-size: 13px; color: #646566; }
.search-row { display: flex; gap: 8px; align-items: center; }
.search-row :deep(.van-field) { flex: 1; padding: 8px 12px; background: #fff; border-radius: 8px; }
.search-result { margin-top: 10px; }
.search-404 { text-align: center; color: #ee0a24; font-size: 13px; margin-top: 8px; }
</style>
