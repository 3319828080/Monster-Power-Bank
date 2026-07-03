export function formatFee(fen: number): string {
  return (fen / 100).toFixed(2)
}

export function formatDuration(minutes: number): string {
  if (minutes < 60) return `${minutes}分钟`
  const h = Math.floor(minutes / 60)
  const m = minutes % 60
  return m > 0 ? `${h}小时${m}分钟` : `${h}小时`
}

export function formatDistance(meters: number): string {
  if (meters < 1000) return `${Math.round(meters)}m`
  return `${(meters / 1000).toFixed(1)}km`
}

export function timeAgo(t: string): string {
  const now = Date.now()
  const then = new Date(t).getTime()
  const diff = Math.floor((now - then) / 1000)
  if (diff < 60) return '刚刚'
  if (diff < 3600) return `${Math.floor(diff / 60)}分钟前`
  if (diff < 86400) return `${Math.floor(diff / 3600)}小时前`
  return `${Math.floor(diff / 86400)}天前`
}

export function orderStatusText(status: string): string {
  const map: Record<string, string> = {
    pending: '待处理',
    borrowed: '租借中',
    '租借中': '租借中',
    completed: '已完成',
    cancelled: '已取消',
    pending_compensation: '待赔付',
    compensated: '已赔付',
  }
  return map[status] || status
}

export function orderStatusColor(status: string): string {
  const map: Record<string, string> = {
    pending: '#ff976a',
    borrowed: '#07c160',
    '租借中': '#07c160',
    completed: '#969799',
    cancelled: '#969799',
    pending_compensation: '#ee0a24',
    compensated: '#969799',
  }
  return map[status] || '#969799'
}

export function slotStatusText(status: string): string {
  const map: Record<string, string> = {
    '空闲': '可借',
    '已借出': '已借出',
    occupied: '可借',
    empty: '已借出',
  }
  return map[status] || status
}
