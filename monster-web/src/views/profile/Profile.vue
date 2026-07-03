<template>
  <div class="profile-page">
    <van-nav-bar title="我的" />

    <!-- User card -->
    <div class="user-card">
      <van-image round width="56" height="56" :src="avatar" />
      <div class="user-info">
        <div class="user-name">{{ nickname || '未设置昵称' }}</div>
        <div class="user-phone">{{ phone || '未绑定手机' }}</div>
      </div>
      <van-icon name="arrow" color="#c8c9cc" @click="showEditProfile = true" />
    </div>

    <!-- Wallet + Deposit cards -->
    <div class="asset-cards">
      <div class="asset-card wallet-card" @click="showWallet = true">
        <div class="asset-label">钱包余额</div>
        <div class="asset-value">¥{{ formatFee(walletBalance) }}</div>
        <div class="asset-sub">查看钱包详情</div>
      </div>
      <div class="asset-card deposit-card">
        <div class="asset-label">押金状态</div>
        <div class="asset-value" :class="depositFrozen > 0 ? 'frozen' : ''">
          {{ depositFrozen > 0 ? `¥${formatFee(depositFrozen)}` : '无需押金' }}
        </div>
        <div class="asset-sub">{{ depositFrozen > 0 ? '归还后退还' : '余额充足免押金' }}</div>
      </div>
    </div>

    <!-- Menu -->
    <van-cell-group inset class="menu-group">
      <van-cell title="当前订单" icon="records" is-link to="/order/current">
        <template #extra><van-badge v-if="hasActiveOrder" dot /></template>
      </van-cell>
      <van-cell title="历史订单" icon="notes" is-link to="/orders" />
      <van-cell title="资金流水" icon="balance-list" is-link @click="showTransactions = true" />
      <van-cell title="押金规则" icon="info-o" is-link @click="showDepositRules = true" />
    </van-cell-group>

    <van-cell-group inset class="menu-group">
      <van-cell title="绑定手机" icon="phone" is-link @click="showBindPhone = true" />
      <van-cell title="编辑资料" icon="edit" is-link @click="showEditProfile = true" />
    </van-cell-group>

    <div class="logout-btn">
      <van-button round block plain type="danger" @click="handleLogout" :loading="loggingOut">退出登录</van-button>
    </div>

    <!-- Wallet popup -->
    <van-popup v-model:show="showWallet" round position="bottom" :style="{ height: '50vh' }">
      <div class="popup-content">
        <h3>钱包</h3>
        <van-cell-group inset>
          <van-cell title="可用余额" :value="`¥${formatFee(walletBalance)}`" />
          <van-cell title="冻结押金" :value="`¥${formatFee(depositFrozen)}`" />
          <van-cell title="累计充值" :value="`¥${formatFee(totalRecharge)}`" />
          <van-cell title="累计消费" :value="`¥${formatFee(totalConsume)}`" />
        </van-cell-group>
        <div class="recharge-section">
          <div class="recharge-title">快速充值</div>
          <div class="recharge-amounts">
            <span class="amount-tag" v-for="a in [5000,10000,20000,50000]" :key="a"
              :class="{ active: rechargeAmount === a }"
              @click="rechargeAmount = a">¥{{ a / 100 }}</span>
          </div>
          <van-button block round type="primary" :loading="recharging" @click="doRecharge" style="margin-top: 12px">
            充值 ¥{{ formatFee(rechargeAmount) }}
          </van-button>
        </div>
      </div>
    </van-popup>

    <!-- Transactions popup -->
    <van-popup v-model:show="showTransactions" round position="bottom" :style="{ height: '60vh' }">
      <div class="popup-content">
        <h3>资金流水</h3>
        <van-list v-model:loading="txLoading" :finished="txFinished" finished-text="没有更多了" @load="loadTransactions">
          <van-cell v-for="tx in transactions" :key="tx.id"
            :title="tx.remark || tx.trade_type"
            :label="tx.created_at"
            :value="tx.amount > 0 ? `+¥${formatFee(tx.amount)}` : `-¥${formatFee(-tx.amount)}`"
            :value-style="{ color: tx.amount > 0 ? '#07c160' : '#ee0a24' }" />
        </van-list>
      </div>
    </van-popup>

    <!-- Deposit rules popup -->
    <van-popup v-model:show="showDepositRules" round position="bottom" :style="{ height: '40vh' }">
      <div class="popup-content">
        <h3>押金规则</h3>
        <div class="rules-text">{{ depositRules || '押金在归还充电宝且无欠费后自动退还，预计1-3个工作日到账。' }}</div>
      </div>
    </van-popup>

    <!-- Edit profile popup -->
    <van-popup v-model:show="showEditProfile" round position="bottom" :style="{ height: '40vh' }">
      <div class="popup-content">
        <h3>编辑资料</h3>
        <van-field v-model="editNickname" label="昵称" placeholder="请输入昵称" />
        <van-field v-model="editAvatar" label="头像URL" placeholder="请输入头像链接" />
        <van-field v-model="editGender" label="性别" placeholder="0-未知 1-男 2-女" type="number" />
        <van-button block round type="primary" @click="saveProfile" :loading="savingProfile" style="margin-top: 16px">保存</van-button>
      </div>
    </van-popup>

    <!-- Bind phone popup -->
    <van-popup v-model:show="showBindPhone" round position="bottom" :style="{ height: '40vh' }">
      <div class="popup-content">
        <h3>绑定手机</h3>
        <van-field v-model="bindPhoneNumber" label="手机号" placeholder="请输入手机号" />
        <van-field v-model="bindSmsCode" label="验证码" placeholder="请输入验证码">
          <template #button>
            <van-button size="small" type="primary" @click="sendBindSms" :loading="sendingBindSms">发送验证码</van-button>
          </template>
        </van-field>
        <van-button block round type="primary" @click="doBindPhone" :loading="bindingPhone" style="margin-top: 16px">确认绑定</van-button>
      </div>
    </van-popup>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import { useAuthStore } from '@/stores/auth'
import { getProfile, updateProfile, logout, bindPhone, sendSmsCode, getWallet, listTransactions, getDepositInfo, rechargeDeposit } from '@/api/user'
import type { TransactionInfo } from '@/api/user'
import { formatFee } from '@/utils'

const router = useRouter()
const auth = useAuthStore()

const nickname = ref('')
const phone = ref('')
const avatar = ref('https://api.dicebear.com/7.x/thumbs/svg?seed=default')
const loggingOut = ref(false)
const hasActiveOrder = ref(false)

// Wallet
const showWallet = ref(false)
const walletBalance = ref(0)
const depositFrozen = ref(0)
const totalRecharge = ref(0)
const totalConsume = ref(0)
const rechargeAmount = ref(10000)
const recharging = ref(false)

async function doRecharge() {
  recharging.value = true
  try {
    await rechargeDeposit(rechargeAmount.value)
    showToast('充值成功')
    fetchWallet()
  } catch { /* */ }
  finally { recharging.value = false }
}

// Deposit
const showDepositRules = ref(false)
const depositRules = ref('')

// Transactions
const showTransactions = ref(false)
const transactions = ref<TransactionInfo[]>([])
const txLoading = ref(false)
const txFinished = ref(false)
let txPage = 1

// Edit profile
const showEditProfile = ref(false)
const editNickname = ref('')
const editAvatar = ref('')
const editGender = ref(0)
const savingProfile = ref(false)

// Bind phone
const showBindPhone = ref(false)
const bindPhoneNumber = ref('')
const bindSmsCode = ref('')
const sendingBindSms = ref(false)
const bindingPhone = ref(false)

async function fetchProfile() {
  try {
    const res = await getProfile()
    const u = res.user
    nickname.value = u.nickname
    phone.value = u.phone
    if (u.avatar) avatar.value = u.avatar
  } catch { /* */ }
}

async function fetchWallet() {
  try {
    const res = await getWallet()
    const w = res.wallet
    walletBalance.value = w.balance
    depositFrozen.value = w.frozen
    totalRecharge.value = w.total_recharge
    totalConsume.value = w.total_consume
  } catch { /* */ }
}

async function fetchDepositInfo() {
  try {
    const info = await getDepositInfo()
    depositRules.value = info.rules
    depositFrozen.value = info.deposit_frozen
  } catch { /* */ }
}

async function loadTransactions() {
  txLoading.value = true
  try {
    const res = await listTransactions({ page: txPage, page_size: 20 })
    if (txPage === 1) {
      transactions.value = res.list
    } else {
      transactions.value.push(...res.list)
    }
    txFinished.value = transactions.value.length >= res.total
    txPage++
  } catch {
    txFinished.value = true
  } finally {
    txLoading.value = false
  }
}

async function saveProfile() {
  savingProfile.value = true
  try {
    const res = await updateProfile({
      nickname: editNickname.value || undefined,
      avatar: editAvatar.value || undefined,
      gender: Number(editGender.value) || 0,
    })
    const u = res.user
    nickname.value = u.nickname
    avatar.value = u.avatar || avatar.value
    showToast('保存成功')
    showEditProfile.value = false
  } catch {
    showToast('保存失败')
  } finally {
    savingProfile.value = false
  }
}

async function sendBindSms() {
  if (!bindPhoneNumber.value) { showToast('请输入手机号'); return }
  sendingBindSms.value = true
  try {
    await sendSmsCode(bindPhoneNumber.value)
    showToast('验证码已发送')
  } catch { /* */ }
  finally { sendingBindSms.value = false }
}

async function doBindPhone() {
  if (!bindPhoneNumber.value || !bindSmsCode.value) { showToast('请填写完整'); return }
  bindingPhone.value = true
  try {
    await bindPhone(bindPhoneNumber.value, bindSmsCode.value)
    phone.value = bindPhoneNumber.value
    showToast('绑定成功')
    showBindPhone.value = false
  } catch { /* */ }
  finally { bindingPhone.value = false }
}

async function handleLogout() {
  try {
    await showConfirmDialog({ title: '退出登录', message: '确定退出登录？' })
  } catch { return }
  loggingOut.value = true
  await auth.logout()
  loggingOut.value = false
  router.replace('/login')
}

onMounted(() => {
  fetchProfile()
  fetchWallet()
  fetchDepositInfo()
})
</script>

<style scoped>
.profile-page { min-height: 100vh; background: #f7f8fa; padding-bottom: 24px; }
.user-card { display: flex; align-items: center; gap: 12px; padding: 20px 16px; background: #fff; margin-bottom: 12px; }
.user-info { flex: 1; }
.user-name { font-size: 18px; font-weight: 600; margin-bottom: 4px; }
.user-phone { font-size: 14px; color: #969799; }

/* Asset cards */
.asset-cards { display: flex; gap: 10px; padding: 0 16px; margin-bottom: 12px; }
.asset-card { flex: 1; background: #fff; border-radius: 10px; padding: 16px; cursor: pointer; }
.wallet-card { background: linear-gradient(135deg, #07c160, #05a34f); color: #fff; }
.deposit-card { background: linear-gradient(135deg, #1989fa, #0d6ed6); color: #fff; }
.asset-label { font-size: 13px; opacity: 0.85; margin-bottom: 6px; }
.asset-value { font-size: 22px; font-weight: 700; margin-bottom: 4px; }
.asset-value.frozen { color: #ffeaa7; }
.asset-sub { font-size: 11px; opacity: 0.75; }

.menu-group { margin-top: 12px; }
.logout-btn { padding: 40px 16px; }
.popup-content { padding: 24px 16px; }
.popup-content h3 { font-size: 18px; font-weight: 600; margin-bottom: 16px; text-align: center; }
.recharge-section { margin-top: 16px; }
.recharge-title { font-size: 15px; font-weight: 600; margin-bottom: 10px; text-align: center; }
.recharge-amounts { display: flex; gap: 10px; justify-content: center; }
.amount-tag { padding: 8px 16px; background: #f7f8fa; border-radius: 20px; font-size: 14px; font-weight: 500; cursor: pointer; border: 1px solid transparent; }
.amount-tag.active { background: #f0faf4; color: #07c160; border-color: #07c160; }
.rules-text { font-size: 14px; color: #646566; line-height: 1.8; white-space: pre-line; }
</style>
