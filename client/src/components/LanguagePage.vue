<template>
  <div class='language-reports'>
    <h1>{{language}}</h1>
    <reports-panel :reports="list"></reports-panel>
  </div>
</template>

<script>
import core from '../core'
import ReportsPanel from './parts/ReportsPanel'
export default {
  name: 'LanguagePage',
  components: { ReportsPanel },
  data() {
    return {
      language: this.$route.params['language'],
      list: core.state.subList
    }
  },
  mounted() {
    core.searchReportsByLanguage(this.language)
  }
}
</script>
