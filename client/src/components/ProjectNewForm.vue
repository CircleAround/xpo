<script>
import BasicProjectForm from './parts/BasicProjectForm'
import core from '../core'

export default {
  extends: BasicProjectForm,
  name: 'ProjectNewForm',
  methods: {
    initialize() {

    },
    doSaveReport() {
      // return core.postReport().then(() => {
      //   core.alert('投稿しました！', 'success')
      // })
    }
  }
}
</script>
