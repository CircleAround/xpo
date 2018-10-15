<template>
  <div class='report'>
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
      list: []
    }
  },
  created() {
    const params = this.$route.params
    core.findReport4Update(params.authorId, params.id).then(() => {
      this.list.push(core.state.targetReport)
    })
  }
}
</script>
