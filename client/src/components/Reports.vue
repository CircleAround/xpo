<template>
  <div>
    <top-languages-panel></top-languages-panel>
    <reports-panel :reports="list"></reports-panel>
  </div>
</template>

<script>
import core from '../core'
import ReportsPanel from './parts/ReportsPanel'
import TopLanguagesPanel from './parts/TopLanguagesPanel'
export default {
  name: 'reports',
  components: { ReportsPanel, TopLanguagesPanel },
  data() {
    return {
      languages: core.state.languageNames,
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
