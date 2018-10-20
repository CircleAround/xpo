<template>
  <div class='users-page'>
    <user-header :xuser="xuser"></user-header>
    <reports-panel :reports="list"></reports-panel>
  </div>
</template>

<script>
import core from '../core'
import UserHeader from './parts/UserHeader'
import ReportsPanel from './parts/ReportsPanel'
export default {
  name: 'reports',
  components: { ReportsPanel, UserHeader },
  data() {
    return {
      xuser: {},
      list: core.state.subList
    }
  },
  async mounted() {
    console.log('mounted')
    const params = this.$route.params
    try {
      const xresponse = await core.getXUserByName(params['author'])
      this.xuser = xresponse.data
      await core.searchByAuthor(this.xuser.id)
    } catch (e) {
      // nop
      console.error(e)
    }
  }
}
</script>
