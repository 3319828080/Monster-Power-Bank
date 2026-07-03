import request from './request'

export interface StationInfo {
  id: number
  name: string
  address: string
  latitude: number
  longitude: number
  distance: number
  availableBanks: number
  openTime?: string
  description?: string
  images?: string
}

export interface OrderInfo {
  id: number
  orderNo: string
  userId: number
  powerBankId: number
  powerBankNo: string
  borrowStationId: number
  borrowStationName: string
  returnStationId: number
  returnStationName: string
  borrowTime: string
  returnTime: string
  durationMinutes: number
  startFee: number
  hourlyFee: number
  totalAmount: number
  paidAmount: number
  discountAmount: number
  dailyCap: number
  status: string
  remark: string
  createdAt: string
  updatedAt: string
}

export interface PreCheckResult {
  allowed: boolean
  reason: string
  depositRequired: number
  startFee: number
  hourlyFee: number
  dailyCap: number
}

export interface RealTimeFee {
  orderNo: string
  elapsedSeconds: number
  currentFee: number
  dailyCap: number
  hourlyFee: number
}

export function preCheckBorrow(data: {
  powerBankId: number
  stationId: number
  cabinetId: number
  slotId: number
}) {
  return request.post<PreCheckResult>('/order/pre-check', data)
}

export function createOrder(data: {
  powerBankId: number
  stationId: number
  cabinetId: number
  slotId: number
}) {
  return request.post<{
    orderNo: string
    startFee: number
    hourlyFee: number
    dailyCap: number
    depositFrozen: number
    borrowTime: string
  }>('/order/create', data)
}

export function cancelOrder(orderNo: string) {
  return request.post('/order/cancel', { orderNo })
}

export function returnPowerBank(data: {
  orderNo: string
  stationId: number
  cabinetId: number
  slotId: number
  couponId?: string
}) {
  return request.post<{
    order_no: string
    duration_minutes: number
    total_amount: number
    paid_amount: number
    discount_amount: number
    balance_deducted: number
    deposit_deducted: number
    return_time: string
  }>('/order/return', data)
}

export function extendRental(orderNo: string, extraMinutes: number) {
  return request.post('/order/extend', { orderNo, extraMinutes })
}

export function getCurrentOrder() {
  return request.get<{ order: OrderInfo | null }>('/order/current')
}

export function getRealTimeFee(orderNo: string) {
  return request.get<RealTimeFee>('/order/fee', { params: { order_no: orderNo } })
}

export function listOrders(params: { page: number; page_size: number; status?: string }) {
  return request.get<{
    list: OrderInfo[]
    total: number
    page: number
    page_size: number
  }>('/order/list', { params })
}

export function createCompensationPayment(orderNo: string) {
  return request.post('/order/compensation', { orderNo })
}

export function getStation(id: number) {
  return request.get<{ station: StationInfo }>(`/station/${id}`)
}

export function listNearbyStations(params: {
  latitude: number
  longitude: number
  radius_meters?: number
}) {
  return request.get<{ list: StationInfo[] }>('/station/nearby', { params })
}

// ---- Cabinet ----

export interface SlotInfo {
  id: number
  slotNo: string
  status: string
  powerBankId: number
  powerBankNo: string
  batteryLevel: number
  power: string
}

export interface CabinetInfo {
  id: number
  cabinetNo: string
  status: number
  totalSlots: number
  availableSlots: number
  occupiedSlots: number
  slots: SlotInfo[]
}

export function getCabinetList(stationId: number) {
  return request.get<{ cabinets: CabinetInfo[] }>('/station/cabinet-list', { params: { station_id: stationId } })
}

// ---- Order Detail (with cabinet/slot info) ----

export interface OrderDetail {
  id: number
  orderNo: string
  userId: number
  powerBankId: number
  powerBankNo: string
  borrowStationId: number
  borrowStationName: string
  borrowCabinetId: number
  borrowCabinetNo: string
  borrowSlotId: number
  borrowSlotNo: string
  returnStationId: number
  returnStationName: string
  returnCabinetId: number
  returnCabinetNo: string
  returnSlotId: number
  returnSlotNo: string
  borrowTime: string
  returnTime: string
  durationMinutes: number
  startFee: number
  hourlyFee: number
  dailyCap: number
  totalAmount: number
  paidAmount: number
  discountAmount: number
  deposit: number
  status: string
  remark: string
  createdAt: string
  updatedAt: string
}

export function getOrderDetail(orderNo: string) {
  return request.get<{ order: OrderDetail }>('/order/detail', { params: { order_no: orderNo } })
}

// ---- Delay Fee ----

export interface DelayFeeInfo {
  orderNo: string
  extraMinutes: number
  additionalFee: number
  estimatedTotal: number
  additionalDeposit: number
  hourlyFee: number
  dailyCap: number
}

export function getDelayFee(data: { orderNo: string; extraMinutes: number }) {
  return request.post<DelayFeeInfo>('/order/delay-fee', data)
}

// ---- Scan Cabinet ----

export interface ScanCabinetResponse {
  stationId: number
  stationName: string
  stationAddress: string
  cabinetId: number
  cabinetNo: string
  availableSlots: number
  totalSlots: number
  startFee: number
  startMins: number
  hourlyFee: number
  dailyCap: number
  deposit: number
}

export function scanCabinet(cabinetId: number) {
  return request.post<ScanCabinetResponse>('/cabinet/scan', { cabinetId })
}

// ---- Return Cabinet List ----

export interface ReturnCabinetInfo {
  cabinetId: number
  cabinetNo: string
  stationId: number
  stationName: string
  stationAddress: string
  latitude: number
  longitude: number
  distance: number
  emptySlotCount: number
  totalSlots: number
}

export function listReturnCabinets(params: {
  latitude: number
  longitude: number
  radiusMeters?: number
}) {
  return request.get<{ list: ReturnCabinetInfo[] }>('/cabinet/return-list', { params })
}

export function searchCabinet(cabinetNo: string) {
  return request.get<{ cabinet: ReturnCabinetInfo }>('/cabinet/search', { params: { cabinetNo } })
}

// ---- Confirm Return (simplified) ----

export function confirmReturn(data: {
  orderNo: string
  cabinetId: number
  couponId?: string
}) {
  return request.post<{
    orderNo: string
    durationMinutes: number
    totalAmount: number
    paidAmount: number
    discountAmount: number
    balanceDeducted: number
    depositDeducted: number
    returnTime: string
  }>('/order/return', data)
}
