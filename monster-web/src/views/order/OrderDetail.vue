<template>
  <div class="order-detail-page">
    <van-nav-bar title="订单详情" left-arrow @click-left="$router.back()" />

    <van-empty v-if="!order" description="未找到订单" />

    <template v-else>
      <!-- Status card -->
      <van-cell-group inset class="status-card">
        <van-cell>
          <template #title>
            <van-tag :color="orderStatusColor(order.status)" size="medium">
              {{ orderStatusText(order.status) }}
            </van-tag>
          </template>
        </van-cell>
        <van-cell title="订单编号" :value="order.orderNo" />
        <van-cell title="充电宝" :value="order.powerBankNo" />
      </van-cell-group>

      <!-- Location card -->
      <van-cell-group inset class="loc-card">
        <van-cell title="借出站点" :value="order.borrowStationName" />
        <van-cell v-if="order.borrowCabinetNo" title="借出机柜" :value="order.borrowCabinetNo" />
        <van-cell v-if="order.borrowSlotNo" title="借出仓位" :value="order.borrowSlotNo" />
        <van-cell v-if="order.returnStationName" title="归还站点" :value="order.returnStationName" />
        <van-cell v-if="order.returnCabinetNo" title="归还机柜" :value="order.returnCabinetNo" />
        <van-cell v-if="order.returnSlotNo" title="归还仓位" :value="order.returnSlotNo" />
      </van-cell-group>

      <!-- Time card -->
      <van-cell-group inset class="time-card">
        <van-cell title="借用时间" :value="formatTime(order.borrowTime)" />
        <van-cell v-if="order.returnTime" title="归还时间" :value="formatTime(order.returnTime)" />
        <van-cell v-if="order.durationMinutes" title="租借时长" :value="formatDuration(order.durationMinutes)" />
      </van-cell-group>

      <!-- Fee card -->
      <van-cell-group inset class="fee-card">
        <van-cell title="计价规则">
          <template #value>
            <span v-if="order.startFee && order.status !== 'completed'">
              起步¥{{ formatFee(order.startFee) }} / 每小时¥{{ formatFee(order.hourlyFee) }}
            </span>
          </template>
        </van-cell>
        <!-- Real-time fee for active orders -->
        <van-cell v-if="order.status === 'borrowed' && realTimeFee" title="当前费用">
          <template #value>
            <span class="live-fee">¥{{ formatFee(realTimeFee.currentFee) }}</span>
          </template>
          <template #label>
            <span v-if="realTimeFee">已租借 {{ formatElapsed(realTimeFee.elapsedSeconds) }}</span>
          </template>
        </van-cell>
        <van-cell title="总费用" :value="`¥${formatFee(order.totalAmount)}`" />
        <van-cell v-if="order.discountAmount > 0" title="优惠" :value="`-¥${formatFee(order.discountAmount)}`" />
        <van-cell v-if="order.deposit" title="押金" :value="`¥${formatFee(order.deposit)}`" />
        <van-cell title="实付" :value="`¥${formatFee(order.paidAmount)}`" title-class="paid-title" value-class="paid-value" />
      </van-cell-group>

      <!-- Actions -->
      <div class="actions" v-if="order.status === 'borrowed'">
        <van-button block round type="primary" @click="goReturn">去归还</van-button>
      </div>
      <div class="actions" v-if="order.status === 'pending_compensation'">
        <van-button block round type="danger" @click="handleCompensation">去赔付</van-button>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getOrderDetail, getRealTimeFee } from '@/api/order'
import { formatFee, formatDuration, orderStatusText, orderStatusColor } from '@/utils'
import type { OrderDetail, RealTimeFee } from '@/api/order'

const route = useRoute()
const router = useRouter()
const orderNo = route.params.orderNo as string
const order = ref<OrderDetail | null>(null)
const realTimeFee = ref<RealTimeFee | null>(null)
let feeTimer: number | null = null

function formatTime(t: string) {
  if (!t) return '-'
  return new Date(t).toLocaleString('zh-CN')
}

function formatElapsed(seconds: number) {
  const m = Math.floor(seconds / 60)
  if (m < 60) return `${m}分钟`
  const h = Math.floor(m / 60)
  return `${h}小时${m % 60}分钟`
}

function goReturn() {
  if (!order.value) return
  router.push(`/return/${order.value.orderNo}`)
}

function handleCompensation() {
  if (!order.value) return
  router.push(`/payment?biz_type=compensation&biz_no=${order.value.orderNo}&amount=${order.value.totalAmount}`)
}

async function startFeePolling() {
  const poll = async () => {
    try {
      realTimeFee.value = await getRealTimeFee(orderNo)
    } catch { /* */ }
  }
  await poll()
  feeTimer = window.setInterval(poll, 10000)
}

onMounted(async () => {
  try {
    const res = await getOrderDetail(orderNo)
    order.value = res.order
    if (order.value?.status === 'borrowed') {
      startFeePolling()
    }
  } catch { /* */ }
})

onUnmounted(() => {
  if (feeTimer !== null) clearInterval(feeTimer)
})
</script>

<style scoped>
.order-detail-page { min-height: 100vh; background: #f7f8fa; padding-bottom: 80px; }
.status-card, .loc-card, .time-card, .fee-card { margin-top: 12px; }
.paid-title { font-weight: 600; }
.paid-value { color: #ee0a24; font-weight: 600; }
.live-fee { color: #ee0a24; font-weight: 600; font-size: 16px; }
.actions { padding: 24px 16px; }
</style>
