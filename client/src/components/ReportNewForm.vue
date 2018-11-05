<script>
import BasicReportForm from './parts/BasicReportForm'
import core from '../core'

export default {
  extends: BasicReportForm,
  name: 'NewReportForm',
  methods: {
    initialize() {
      core.initNewReport()
      this.updateMarkdown()
    },
    doPostReport() {
      return core.postReport().then(() => {
        core.alert('投稿しました！', 'success')
      })
    }
  }
}
</script>
