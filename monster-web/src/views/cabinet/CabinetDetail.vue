<template>
  <div class="cabinet-page">
    <van-nav-bar title="机柜详情" left-arrow @click-left="$router.back()" />

    <van-loading v-if="!cabinet" class="page-loading" size="24px" />

    <template v-else>
      <!-- Station info -->
      <div class="cabinet-header">
        <div class="station-name">{{ cabinet.stationName }}</div>
        <div class="station-addr">
          <van-icon name="location-o" />
          {{ cabinet.stationAddress }}
        </div>
      </div>

      <!-- Cabinet info card -->
      <van-cell-group inset class="info-card">
        <van-cell title="机柜编号" :value="cabinet.cabinetNo" />
        <van-cell title="空闲仓位">
          <template #value>
            <span :class="cabinet.availableSlots > 0 ? 'green' : 'red'">
              {{ cabinet.availableSlots }} / {{ cabinet.totalSlots }}
            </span>
          </template>
        </van-cell>
      </van-cell-group>

      <!-- Pricing card -->
      <van-cell-group inset class="pricing-card">
        <van-cell title="计价规则" />
        <van-cell title="起步价" :value="`¥${formatFee(cabinet.startFee)} / ${cabinet.startMins}分钟`" />
        <van-cell title="每小时" :value="`¥${formatFee(cabinet.hourlyFee)}`" />
        <van-cell title="每日封顶" :value="`¥${formatFee(cabinet.dailyCap)}`" />
        <van-cell title="押金" :value="`¥${formatFee(cabinet.deposit)}（归还后退还）`" />
      </van-cell-group>

      <!-- Station detail link -->
      <div class="action-section">
        <van-button block round type="primary" @click="$router.push(`/station/${cabinet.stationId}`)">
          查看站点详情
        </van-button>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { scanCabinet } from '@/api/order'
import type { ScanCabinetResponse } from '@/api/order'
import { formatFee } from '@/utils'

const route = useRoute()
const cabinetId = Number(route.params.id)
const cabinet = ref<ScanCabinetResponse | null>(null)

onMounted(async () => {
  try {
    cabinet.value = await scanCabinet(cabinetId)
  } catch { /* */ }
})
</script>

<style scoped>
.cabinet-page { min-height: 100vh; background: #f7f8fa; }
.page-loading { display: flex; justify-content: center; padding: 80px 0; }
.cabinet-header { background: linear-gradient(135deg, #1989fa, #07c160); padding: 24px 16px; color: #fff; }
.station-name { font-size: 20px; font-weight: 600; margin-bottom: 8px; }
.station-addr { font-size: 13px; display: flex; align-items: center; gap: 4px; opacity: 0.9; }
.info-card, .pricing-card { margin-top: 12px; }
.green { color: #07c160; font-weight: 600; }
.red { color: #ee0a24; font-weight: 600; }
.action-section { padding: 24px 16px; }
</style>
