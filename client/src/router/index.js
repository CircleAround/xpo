import Vue from 'vue'
import Router from 'vue-router'
import Reports from '@/components/Reports'
import ReportEditForm from '@/components/ReportEditForm'
import ReportNewForm from '@/components/ReportNewForm'
import SignupForm from '@/components/SignupForm'
import ProfileEditForm from '@/components/ProfileEditForm'
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
      path: '/reports/new',
      name: 'ReportNewForm',
      component: ReportNewForm
    },
    {
      path: '/reports/:author_id/:id/edit',
      name: 'ReportEditForm',
      component: ReportEditForm
    },
    {
      path: '/signup',
      name: 'SignupForm',
      component: SignupForm
    },
    {
      path: '/users/me/edit',
      name: 'ProfileEditForm',
      component: ProfileEditForm
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
