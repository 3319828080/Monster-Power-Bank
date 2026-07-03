<template>
  <div class="current-order-page">
    <van-nav-bar title="当前订单" left-arrow @click-left="$router.back()" />

    <van-loading v-if="!order" class="page-loading" size="24px" />

    <template v-else>
      <!-- Status banner -->
      <div class="status-banner" :style="{ background: statusColor }">
        <div class="status-text">{{ statusText }}</div>
        <div class="status-sub">
          借用时间: {{ order.borrowTime }}
          <template v-if="fee"> | 已租 {{ Math.floor(fee.elapsedSeconds / 60) }} 分钟</template>
        </div>
      </div>

      <!-- Fee card -->
      <van-cell-group inset class="fee-card">
        <van-cell title="当前费用" :value="`¥${formatFee(fee?.currentFee || 0)}`" />
        <van-cell title="每小时" :value="`¥${formatFee(fee?.hourlyFee || order.hourlyFee)}`" />
        <van-cell title="每日封顶" :value="`¥${formatFee(fee?.dailyCap || 0)}`" />
      </van-cell-group>

      <!-- Order info -->
      <van-cell-group inset class="info-card">
        <van-cell title="订单编号" :value="order.orderNo" />
        <van-cell title="借出站点" :value="order.borrowStationName || '-'" />
        <van-cell title="充电宝编号" :value="order.powerBankNo" />
      </van-cell-group>

      <!-- Actions -->
      <div class="actions">
        <van-button block round icon="replay" @click="handleExtend">延长租借</van-button>
        <van-button block round type="primary" icon="location" @click="goReturn">归还充电宝</van-button>
        <van-button v-if="order?.status !== '租借中'" block round plain type="danger" @click="handleCancel" :loading="cancelling">取消订单</van-button>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'

import { useRouter } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import { useOrderStore } from '@/stores/order'
import { formatFee, orderStatusText, orderStatusColor } from '@/utils'
import { extendRental, getRealTimeFee, cancelOrder } from '@/api/order'

const router = useRouter()
const orderStore = useOrderStore()
const order = computed(() => orderStore.currentOrder)
const fee = computed(() => orderStore.realTimeFee)
const statusText = computed(() => order.value ? orderStatusText(order.value.status) : '')
const statusColor = computed(() => order.value ? orderStatusColor(order.value.status) : '#07c160')
const cancelling = ref(false)

async function goReturn() {
  if (!order.value) return
  router.push(`/return/${order.value.orderNo}`)
}

async function handleCancel() {
  if (!order.value) return
  try {
    await showConfirmDialog({
      title: '取消订单',
      message: '确定取消当前租借订单？',
      confirmButtonText: '确认取消',
      confirmButtonColor: '#ee0a24',
    })
    cancelling.value = true
    await cancelOrder(order.value.orderNo)
    showToast('订单已取消')
    orderStore.refreshCurrentOrder()
  } catch { /* user cancelled or error */ }
  finally { cancelling.value = false }
}

async function handleExtend() {
  if (!order.value) return
  try {
    await showConfirmDialog({
      title: '延长租借',
      message: '确认延长租借30分钟？',
      confirmButtonText: '确认延长',
    })
    await extendRental(order.value.orderNo, 30)
    showToast('延长成功')
  } catch { /* user cancelled or error */ }
}

onMounted(async () => {
  await orderStore.refreshCurrentOrder()
  if (order.value) {
    orderStore.startFeePolling(order.value.orderNo)
  }
})

onUnmounted(() => {
  orderStore.stopFeePolling()
})
</script>

<style scoped>
.current-order-page { min-height: 100vh; background: #f7f8fa; }
.page-loading { display: flex; justify-content: center; padding: 80px 0; }
.status-banner { padding: 24px 16px; color: #fff; }
.status-text { font-size: 20px; font-weight: 600; margin-bottom: 4px; }
.status-sub { font-size: 13px; opacity: 0.9; }
.fee-card, .info-card { margin-top: 12px; }
.actions { display: flex; flex-direction: column; gap: 12px; padding: 24px 16px; }
</style>
