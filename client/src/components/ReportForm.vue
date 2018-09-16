<template>
  <div class="report_form">
    <div v-if="state.me.id">
      <div class="editor">
        <textarea v-model="state.newReport.content" class="textcontent" @keyup='updateMarkdown()'></textarea>
        <div class="preview markdown" v-html="markdown"></div>
      </div>
      <div class="errors" v-if="errors.length > 0">
        <div class="error" v-for='(item, key , index) in errors' v-bind:key="index">
          {{item}}
        </div>
      </div>
      <div class="actions">
        <el-button type="success" icon="el-icon-check" circle @click='postReport()'></el-button>
      </div>
    </div>
    <div v-if="!state.me.id">
      ログインすると使えます
    </div>
  </div>
</template>

<script>
import core from '../core'
import marked from 'marked'
export default {
  name: 'report_form',
  data() {
    return {
      markdown: '',
      errors: [],
      state: core.state
    }
  },
  created() {
    console.log('created')
    console.log(core)
    console.log(core.state)
  },
  methods: {
    postReport() {
      this.errors = []
      core
        .postReport()
        .then(() => {
          this.$message({
            showClose: true,
            message: '投稿しました！',
            type: 'success',
            center: true
          })
        })
        .catch(error => {
          core.eachResponseErrors(error, (msg, type, property) => {
            this.errors.push(msg)
          })
        })
    },
    updateMarkdown() {
      this.markdown = marked(this.state.newReport.content)
    }
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style scoped>
.editor {
  display: flex;
}

.textcontent {
  flex: 1;
  width: 100%;
  height: 10em;

  font-size: 16px;
  border: solid 1px #ccc;
}

.preview {
  flex: 1;
  margin: 0 10px;
  padding: 10px;
  background-color: #f2f2f2;
}

.actions {
  padding: 5px;
  text-align: right;
}
</style>
