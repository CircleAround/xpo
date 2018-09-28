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
      list: core.state.list
    }
  },
  mounted() {
    console.log('mounted')
    if (core.state.list.length === 0) {
      core.retriveReports()
    } else {
      core.forceUpdateMainList() // for update report list
    }
  }
}
</script>
