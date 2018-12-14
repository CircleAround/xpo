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
        .findReport4Update(this.$route.params.authorId, this.$route.params.id)
        .then(() => {
          this.updateMarkdown()
          this.loading = false
        })
    },
    doSaveReport() {
      return core.updateReport(this.$route.params).then(() => {
        core.alert('更新しました！', 'success')
      })
    }
  }
}
</script>
