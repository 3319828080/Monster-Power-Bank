<template>
  <div class="login-page">
    <div class="login-header">
      <img class="logo" src="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 100 100'%3E%3Ccircle cx='50' cy='50' r='45' fill='%2307c160'/%3E%3Cpath d='M35 55 L45 65 L65 40' stroke='white' stroke-width='6' fill='none' stroke-linecap='round' stroke-linejoin='round'/%3E%3C/svg%3E" alt="logo" />
      <h1>怪兽充电</h1>
      <p class="subtitle">随时随地，充电无忧</p>
    </div>

    <van-form @submit="handleLogin" class="login-form">
      <van-cell-group inset>
        <van-field
          v-model="phone"
          name="phone"
          label="手机号"
          placeholder="请输入手机号"
          type="tel"
          maxlength="11"
          :rules="[{ required: true, message: '请输入手机号' }, { pattern: /^1\d{10}$/, message: '手机号格式不正确' }]"
        />
        <van-field
          v-model="code"
          name="code"
          label="验证码"
          placeholder="请输入验证码"
          maxlength="6"
          :rules="[{ required: true, message: '请输入验证码' }]"
        >
          <template #button>
            <van-button
              size="small"
              type="primary"
              :disabled="sending || countdown > 0"
              @click="handleSendCode"
            >
              {{ countdown > 0 ? `${countdown}s` : '获取验证码' }}
            </van-button>
          </template>
        </van-field>
      </van-cell-group>

      <div style="margin: 16px">
        <van-button round block type="primary" native-type="submit" :loading="loading">
          登录 / 注册
        </van-button>
      </div>
    </van-form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showLoadingToast, closeToast } from 'vant'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()

const phone = ref('')
const code = ref('')
const loading = ref(false)
const sending = ref(false)
const countdown = ref(0)
let timer: number | null = null

async function handleSendCode() {
  if (!/^1\d{10}$/.test(phone.value)) {
    showToast('请输入正确的手机号')
    return
  }
  sending.value = true
  try {
    await auth.sendSmsCode(phone.value)
    showToast('验证码已发送')
    countdown.value = 60
    timer = window.setInterval(() => {
      countdown.value--
      if (countdown.value <= 0 && timer) {
        clearInterval(timer)
        timer = null
      }
    }, 1000)
  } catch {
    // toast handled by interceptor
  } finally {
    sending.value = false
  }
}

async function handleLogin() {
  loading.value = true
  try {
    await auth.login(phone.value, code.value)
    router.replace('/')
  } catch {
    // toast handled by interceptor
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: #f7f8fa;
  display: flex;
  flex-direction: column;
  align-items: center;
}
.login-header {
  text-align: center;
  padding: 60px 0 40px;
}
.logo {
  width: 72px;
  height: 72px;
  margin-bottom: 12px;
}
.login-header h1 {
  font-size: 28px;
  color: #323233;
  margin-bottom: 8px;
}
.subtitle {
  font-size: 14px;
  color: #969799;
}
.login-form {
  width: 100%;
  max-width: 400px;
}
</style>
