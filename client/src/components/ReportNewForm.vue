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
        this.$message({
          showClose: true,
          message: '投稿しました！',
          type: 'success',
          center: true
        })
      })
    }
  }
}
</script>
