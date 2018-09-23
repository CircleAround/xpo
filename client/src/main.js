// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import moment from 'vue-moment'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import App from './App'
import router from './router'
import core from './core'

Vue.config.productionTip = false

Vue.use(ElementUI)
Vue.use(moment)

// TODO: loading view...
core.initialize().then(() => {
  /* eslint-disable no-new */
  new Vue({
    el: '#app',
    router,
    components: { App },
    template: '<App/>'
  })
})
