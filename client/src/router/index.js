import Vue from 'vue'
import Router from 'vue-router'
import Report from '@/components/Report'
import Reports from '@/components/Reports'
import UserReportsYmd from '@/components/UserReportsYmd'
import ReportEditForm from '@/components/ReportEditForm'
import ReportNewForm from '@/components/ReportNewForm'
import ProjectNewForm from '@/components/ProjectNewForm'
import SignupForm from '@/components/SignupForm'
import ProfileEditForm from '@/components/ProfileEditForm'
import About from '@/components/About'
import UserPage from '@/components/UserPage'
import LanguagePage from '@/components/LanguagePage'

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
      path: '/reports/:authorId/:id/edit',
      name: 'ReportEditForm',
      component: ReportEditForm
    },
    {
      path: '/reports/:authorId/:id',
      name: 'Report',
      component: Report
    },
    {
      path: '/reports/:authorId/_/:year/:month/:day',
      name: 'UserReportsYmd',
      component: UserReportsYmd
    },
    {
      path: '/projects/new',
      name: 'ProjectNewForm',
      component: ProjectNewForm
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
    },
    {
      path: '/languages/:language',
      name: 'LanguagePage',
      component: LanguagePage
    },
    {
      path: '/:author',
      name: 'UserPage',
      component: UserPage
    }
  ]
}

if (process.env.NODE_ENV === 'production') {
  routerConfig.mode = 'history'
}

export default new Router(routerConfig)
