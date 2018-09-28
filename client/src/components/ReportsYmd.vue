<template>
  <div class='newest-reports'>
    <reports-panel :reports="list"></reports-panel>
  </div>
</template>

<script>
import core from '../core'
import ReportsPanel from './parts/ReportsPanel'
export default {
  name: 'reports',
  components: { ReportsPanel },
  data() {
    return {
      list: core.state.subList
    }
  },
  created() {
    const params = this.$route.params
    core.searchReportsYmd(
      params.authorId,
      params.year,
      params.month,
      params.day
    )
  }
}
</script>
