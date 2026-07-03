import { createRouter, createWebHashHistory } from 'vue-router'
import { getToken } from '@/utils/auth'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/login/Login.vue'),
      meta: { title: '登录' },
    },
    // Tab pages
    {
      path: '/',
      name: 'Home',
      component: () => import('@/views/home/Home.vue'),
      meta: { title: '怪兽充电', requiresAuth: true },
    },
    {
      path: '/orders',
      name: 'OrderList',
      component: () => import('@/views/order/OrderList.vue'),
      meta: { title: '订单', requiresAuth: true },
    },
    {
      path: '/profile',
      name: 'Profile',
      component: () => import('@/views/profile/Profile.vue'),
      meta: { title: '我的', requiresAuth: true },
    },
    // Sub pages (pushed on top of tabs)
    {
      path: '/station/:id',
      name: 'StationDetail',
      component: () => import('@/views/station/StationDetail.vue'),
      meta: { title: '站点详情', requiresAuth: true },
    },
    {
      path: '/cabinet/:id',
      name: 'CabinetDetail',
      component: () => import('@/views/cabinet/CabinetDetail.vue'),
      meta: { title: '机柜详情', requiresAuth: true },
    },
    {
      path: '/order/current',
      name: 'CurrentOrder',
      component: () => import('@/views/order/CurrentOrder.vue'),
      meta: { title: '当前订单', requiresAuth: true },
    },
    {
      path: '/order/detail/:orderNo',
      name: 'OrderDetail',
      component: () => import('@/views/order/OrderDetail.vue'),
      meta: { title: '订单详情', requiresAuth: true },
    },
    {
      path: '/return/:orderNo',
      name: 'ReturnPowerBank',
      component: () => import('@/views/return/ReturnPowerBank.vue'),
      meta: { title: '归还充电宝', requiresAuth: true },
    },
    {
      path: '/payment',
      name: 'Payment',
      component: () => import('@/views/payment/PaymentConfirm.vue'),
      meta: { title: '支付', requiresAuth: true },
    },
  ],
})

router.beforeEach((to, _from, next) => {
  document.title = (to.meta.title as string) || '怪兽充电'
  if (to.meta.requiresAuth && !getToken()) {
    next('/login')
  } else {
    next()
  }
})

export default router
