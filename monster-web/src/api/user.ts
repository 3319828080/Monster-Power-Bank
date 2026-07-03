import request from './request'

export interface UserInfo {
  id: number
  open_id: string
  union_id: string
  nickname: string
  avatar: string
  phone: string
  gender: number
  status: number
  created_at: string
  updated_at: string
}

export interface LoginResponse {
  token: string
  user: UserInfo
}

export interface WalletInfo {
  id: number
  user_id: number
  balance: number
  frozen: number
  total_recharge: number
  total_consume: number
  status: number
  created_at: string
  updated_at: string
}

export interface TransactionInfo {
  id: number
  user_id: number
  wallet_id: number
  order_id: string
  trade_type: string
  amount: number
  balance_before: number
  balance_after: number
  remark: string
  created_at: string
}

// ---- Auth ----

export function sendSmsCode(phone: string) {
  return request.post('/user/send-sms-code', { phone })
}

export function login(phone: string, code: string) {
  return request.post<LoginResponse>('/user/login-phone', { phone, code })
}

export function loginWechat(code: string, nickname?: string, avatar?: string) {
  return request.post<LoginResponse>('/user/login', { code, nickname, avatar })
}

export function logout() {
  return request.post('/user/logout', {})
}

// ---- Profile ----

export function getProfile() {
  return request.get<{ user: UserInfo }>('/user/profile')
}

export function updateProfile(data: { nickname?: string; avatar?: string; gender?: number }) {
  return request.put<{ user: UserInfo }>('/user/profile', data)
}

export function bindPhone(phone: string, code: string) {
  return request.post('/user/bind-phone', { phone, code })
}

// ---- Wallet ----

export function getWallet() {
  return request.get<{ wallet: WalletInfo }>('/wallet')
}

export function listTransactions(params: { page: number; page_size: number }) {
  return request.get<{
    list: TransactionInfo[]
    total: number
    page: number
    page_size: number
  }>('/wallet/transactions', { params })
}

export function payDeposit(amount: number) {
  return request.post('/wallet/deposit/pay', { amount })
}

export function rechargeDeposit(amount: number) {
  return request.post('/wallet/deposit/recharge', { amount })
}

export function refundDeposit() {
  return request.post('/wallet/deposit/refund', {})
}

export interface DepositInfo {
  has_deposit: boolean
  deposit_amount: number
  deposit_frozen: number
  wallet_balance: number
  deposit_enough: boolean
  rules: string
}

export function getDepositInfo() {
  return request.get<DepositInfo>('/wallet/deposit-info')
}
