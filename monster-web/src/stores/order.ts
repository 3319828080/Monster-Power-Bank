import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { OrderInfo, RealTimeFee } from '@/api/order'
import { getCurrentOrder as apiGetCurrent, getRealTimeFee as apiGetFee } from '@/api/order'

export const useOrderStore = defineStore('order', () => {
  const currentOrder = ref<OrderInfo | null>(null)
  const realTimeFee = ref<RealTimeFee | null>(null)
  const selectedOrder = ref<OrderInfo | null>(null)
  const feeTimer = ref<number | null>(null)

  async function refreshCurrentOrder() {
    try {
      const res = await apiGetCurrent()
      currentOrder.value = res.order || null
    } catch {
      currentOrder.value = null
    }
  }

  async function startFeePolling(orderNo: string) {
    stopFeePolling()
    const poll = async () => {
      try {
        realTimeFee.value = await apiGetFee(orderNo)
      } catch { /* ignore */ }
    }
    await poll()
    feeTimer.value = window.setInterval(poll, 10000)
  }

  function stopFeePolling() {
    if (feeTimer.value !== null) {
      clearInterval(feeTimer.value)
      feeTimer.value = null
    }
    realTimeFee.value = null
  }

  function setSelectedOrder(order: OrderInfo) { selectedOrder.value = order }
  function clearSelectedOrder() { selectedOrder.value = null }

  return { currentOrder, realTimeFee, selectedOrder, refreshCurrentOrder, startFeePolling, stopFeePolling, setSelectedOrder, clearSelectedOrder }
})
