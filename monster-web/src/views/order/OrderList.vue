<template>
  <div class="order-list-page">
    <van-nav-bar title="订单" />

    <van-tabs v-model:active="activeTab" @change="onTabChange">
      <van-tab title="全部" name="" />
      <van-tab title="进行中" name="租借中" />
      <van-tab title="已完成" name="completed" />
      <van-tab title="待赔付" name="pending_compensation" />
    </van-tabs>

    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="loadOrders"
      >
        <div
          v-for="item in list"
          :key="item.id"
          class="order-item"
          @click="$router.push(`/order/detail/${item.orderNo}`)"
        >
          <div class="order-top">
            <span class="order-station">{{ item.borrowStationName || '充电宝' }}</span>
            <van-tag :color="orderStatusColor(item.status)">
              {{ orderStatusText(item.status) }}
            </van-tag>
          </div>
          <div class="order-mid">
            <span class="order-no">{{ item.powerBankNo }}</span>
            <span class="order-amount">¥{{ formatFee(item.totalAmount) }}</span>
          </div>
          <div class="order-bottom">
            <span>{{ item.borrowTime }}</span>
            <span v-if="item.durationMinutes">{{ formatDuration(item.durationMinutes) }}</span>
          </div>
        </div>
      </van-list>
    </van-pull-refresh>

    <van-empty v-if="!loading && list.length === 0" description="暂无订单" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { listOrders } from '@/api/order'
import type { OrderInfo } from '@/api/order'
import { formatFee, formatDuration, orderStatusText, orderStatusColor } from '@/utils'

const activeTab = ref('')
const list = ref<OrderInfo[]>([])
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
let page = 1

async function loadOrders() {
  loading.value = true
  try {
    const res = await listOrders({
      page,
      page_size: 10,
      status: activeTab.value || undefined,
    })
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
    refreshing.value = false
  }
}

function onTabChange() {
  page = 1
  list.value = []
  finished.value = false
  loadOrders()
}

function onRefresh() {
  page = 1
  finished.value = false
  refreshing.value = true
  loadOrders()
}
</script>

<style scoped>
.order-list-page { min-height: 100vh; background: #f7f8fa; }
.order-item { padding: 14px 16px; background: #fff; border-bottom: 1px solid #f5f5f5; cursor: pointer; }
.order-item:active { background: #f9f9f9; }
.order-top { display: flex; justify-content: space-between; align-items: center; margin-bottom: 6px; }
.order-station { font-size: 15px; font-weight: 500; }
.order-mid { display: flex; justify-content: space-between; align-items: center; margin-bottom: 4px; }
.order-no { font-size: 13px; color: #646566; }
.order-amount { font-size: 16px; font-weight: 600; color: #ee0a24; }
.order-bottom { font-size: 12px; color: #969799; display: flex; gap: 12px; }
</style>
