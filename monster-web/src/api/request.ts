import axios, { type AxiosRequestConfig } from 'axios'
import { showToast } from 'vant'
import { getToken, removeToken } from '@/utils/auth'
import router from '@/router'

const instance = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
})

instance.interceptors.request.use((config) => {
  const token = getToken()
  if (token && config.headers) {
    config.headers['Authorization'] = `Bearer ${token}`
  }
  return config
})

instance.interceptors.response.use(
  (res) => res.data,
  (error) => {
    const msg =
      error.response?.data?.message ||
      error.response?.data?.reason ||
      error.message ||
      '请求失败'
    if (error.response?.status === 401) {
      removeToken()
      router.replace('/login')
      showToast('登录已过期，请重新登录')
    } else {
      showToast(msg)
    }
    return Promise.reject(error)
  },
)

// Typed wrapper returning Promise<T> directly instead of AxiosResponse<T>
const request = {
  get<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
    return instance.get(url, config) as Promise<T>
  },
  post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    return instance.post(url, data, config) as Promise<T>
  },
  put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    return instance.put(url, data, config) as Promise<T>
  },
}

export default request
