<script>
import BasicReportForm from './parts/BasicReportForm'
import core from '../core'

export default {
  extends: BasicReportForm,
  name: 'EditReportForm',
  methods: {
    initialize() {
      this.loading = true
      core
        .findReport4Update(this.$route.params.author_id, this.$route.params.id)
        .then(() => {
          this.updateMarkdown()
          this.loading = false
        })
    },
    doPostReport() {
      return core.updateReport(this.$route.params).then(() => {
        this.$message({
          showClose: true,
          message: '更新しました！',
          type: 'success',
          center: true
        })
      })
    }
  }
}
</script>
