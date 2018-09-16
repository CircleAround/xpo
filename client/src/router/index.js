import Vue from 'vue'
import Router from 'vue-router'
import Reports from '@/components/Reports'
import ReportForm from '@/components/ReportForm'
import SignupForm from '@/components/SignupForm'
import About from '@/components/About'

Vue.use(Router)

const routerConfig = {
  routes: [
    {
      path: '/',
      name: 'Reports',
      component: Reports
    },
    {
      path: '/report',
      name: 'ReportForm',
      component: ReportForm
    },
    {
      path: '/signup',
      name: 'SignupForm',
      component: SignupForm
    },
    {
      path: '/about',
      name: 'About',
      component: About
    }
  ]
}

if (process.env.NODE_ENV === 'production') {
  routerConfig.mode = 'history'
}

export default new Router(routerConfig)
