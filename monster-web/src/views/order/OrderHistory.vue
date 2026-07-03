<template>
  <div class="order-history-page">
    <van-nav-bar title="历史订单" left-arrow @click-left="$router.back()" />

    <van-tabs v-model:active="activeTab" @change="loadOrders">
      <van-tab title="全部" name="" />
      <van-tab title="已完成" name="completed" />
      <van-tab title="待赔付" name="pending_compensation" />
    </van-tabs>

    <van-list
      v-model:loading="loading"
      :finished="finished"
      finished-text="没有更多了"
      @load="loadOrders"
    >
      <van-cell
        v-for="item in list"
        :key="item.id"
        :title="item.powerBankNo || '充电宝'"
        :label="item.borrowTime"
        :value="`¥${formatFee(item.totalAmount)}`"
        is-link
        @click="goDetail(item)"
      >
        <template #extra>
          <van-tag :color="orderStatusColor(item.status)">
            {{ orderStatusText(item.status) }}
          </van-tag>
        </template>
      </van-cell>
    </van-list>

    <van-empty v-if="!loading && list.length === 0" description="暂无订单" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { listOrders } from '@/api/order'
import type { OrderInfo } from '@/api/order'
import { useOrderStore } from '@/stores/order'
import { formatFee, orderStatusText, orderStatusColor } from '@/utils'

const router = useRouter()
const orderStore = useOrderStore()
const activeTab = ref('')
const list = ref<OrderInfo[]>([])
const loading = ref(false)
const finished = ref(false)
let page = 1

function goDetail(order: OrderInfo) {
  orderStore.setSelectedOrder(order)
  router.push(`/order/detail/${order.orderNo}`)
}

async function loadOrders() {
  loading.value = true
  try {
    const status = activeTab.value || undefined
    const res = await listOrders({ page, page_size: 10, status })
    if (page === 1) {
      list.value = res.list
    } else {
      list.value.push(...res.list)
    }
    finished.value = list.value.length >= res.total
    page++
  } catch {
    finished.value = true
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.order-history-page { min-height: 100vh; background: #f7f8fa; }
</style>
