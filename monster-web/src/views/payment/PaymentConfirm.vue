<template>
  <div class="payment-page">
    <van-nav-bar title="支付" left-arrow @click-left="$router.back()" />

    <div class="amount-card">
      <div class="amount-label">支付金额</div>
      <div class="amount-value">¥{{ formatFee(amount) }}</div>
    </div>

    <van-cell-group inset class="method-card">
      <van-cell title="支付方式" />
      <van-radio-group v-model="channel">
        <van-cell clickable @click="channel = 'wechat'">
          <template #title>
            <span style="display:flex;align-items:center;gap:8px">
              <span style="font-size:20px">💚</span> 微信支付
            </span>
          </template>
          <template #right-icon><van-radio name="wechat" /></template>
        </van-cell>
        <van-cell clickable @click="channel = 'alipay'">
          <template #title>
            <span style="display:flex;align-items:center;gap:8px">
              <span style="font-size:20px">💙</span> 支付宝
            </span>
          </template>
          <template #right-icon><van-radio name="alipay" /></template>
        </van-cell>
      </van-radio-group>
    </van-cell-group>

    <div class="bottom-action">
      <van-button round block type="primary" size="large" :loading="paying" @click="handlePay">
        确认支付 ¥{{ formatFee(amount) }}
      </van-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showToast } from 'vant'
import { createPayment } from '@/api/payment'
import { formatFee } from '@/utils'

const route = useRoute()
const router = useRouter()

const bizType = (route.query.biz_type as string) || ''
const bizNo = (route.query.biz_no as string) || ''
const amount = Number(route.query.amount || 0)
const channel = ref('wechat')
const paying = ref(false)

async function handlePay() {
  paying.value = true
  try {
    const res = await createPayment({
      biz_type: bizType,
      biz_no: bizNo,
      channel: channel.value,
      amount: amount,
      description: `支付${bizType === 'deposit' ? '押金' : bizType === 'compensation' ? '赔付' : '租金'}`,
    })
    // Redirect to payment channel
    if (res.pay_params?.redirect_url) {
      window.location.href = res.pay_params.redirect_url
    } else {
      showToast('支付发起成功')
      router.back()
    }
  } catch {
    // handled by interceptor
  } finally {
    paying.value = false
  }
}

onMounted(() => {
  if (!bizNo || amount <= 0) {
    showToast('参数错误')
    router.back()
  }
})
</script>

<style scoped>
.payment-page { min-height: 100vh; background: #f7f8fa; padding-bottom: 80px; }
.amount-card { text-align: center; padding: 40px 16px; background: #fff; }
.amount-label { font-size: 14px; color: #969799; margin-bottom: 8px; }
.amount-value { font-size: 36px; font-weight: 700; color: #323233; }
.method-card { margin-top: 12px; }
.bottom-action { position: fixed; bottom: 0; left: 0; right: 0; padding: 12px 16px; padding-bottom: calc(12px + env(safe-area-inset-bottom)); background: #fff; }
</style>
