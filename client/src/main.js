// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import VueI18n from 'vue-i18n'
import moment from 'vue-moment'
import App from './App'
import router from './router'
import core, { setServices } from './core'
import messages from './i18n'
import DefaultServiceFactory from './app/DefaultServicesFactory'
import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
import 'material-design-icons-iconfont/dist/material-design-icons.css'

Vue.use(Vuetify)

Vue.use(VueI18n)

Vue.config.productionTip = false

Vue.use(moment)

// @see https://jp.vuejs.org/v2/guide/custom-directive.html
Vue.directive('focus', {
  inserted: function(el) {
    el.focus()
  }
})

Vue.use(VueI18n)
var i18n = new VueI18n({
  locale: 'ja',
  fallbackLocale: 'ja',
  messages: messages
})

setServices(new DefaultServiceFactory().create())

// TODO: loading view...
core.initialize().then(() => {
  /* eslint-disable no-new */
  new Vue({
    el: '#app',
    router,
    components: { App },
    template: '<App/>',
    i18n: i18n
  })
})
