<template>
  <div class='newest-reports'>
    <user-header :xuser="xuser"></user-header>
    <div class="ymd">{{year}}/{{month}}/{{day}}</div>
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
    const params = this.$route.params
    return {
      xuser: {},
      list: core.state.subList,
      year: params.year,
      month: params.month,
      day: params.day
    }
  },
  async created() {
    const params = this.$route.params

    const xresponse = await core.getXUserByName(params.authorId)
    this.xuser = xresponse.data

    core.searchReportsYmd(
      params.authorId,
      params.year,
      params.month,
      params.day
    )
  }
}
</script>
<style scoped>
.ymd {
  text-align: center;
  font-weight: bold;
  margin-bottom: 4px;
}
</style>
