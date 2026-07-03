import request from './request'

export interface PaymentInfo {
  payment_no: string
  channel: string
  amount: number
  status: string
}

export function createPayment(data: {
  biz_type: string
  biz_no: string
  channel: string
  amount: number
  description?: string
}) {
  return request.post<{
    payment_no: string
    pay_params: Record<string, string>
    pay_url?: string
  }>('/payment/create', data)
}

export function markNotificationRead(ids: number[]) {
  return request.post('/notifications/read', { ids })
}
