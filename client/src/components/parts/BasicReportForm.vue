<template>
  <div class="report_form">
    <div v-if="state.me.id">
      <div class="editor">
        <overlay :visible="loading"></overlay>
        <textarea v-model="state.newReport.content" class="textcontent" @keyup='updateMarkdown()'></textarea>
        <div class="preview markdown" v-html="markdown"></div>
      </div>
      <div class="errors" v-if="errors.length > 0">
        <div class="error" v-for='(item, key , index) in errors' v-bind:key="index">
          {{item}}
        </div>
      </div>
      <div class="actions">
        <overlay :visible="loading"></overlay>
        <el-button type="success" icon="el-icon-check" circle @click='postReport()'></el-button>
      </div>
    </div>
    <div v-if="!state.me.id">
      ログインすると使えます
    </div>
  </div>
</template>

<script>
import core from '../../core'
import marked from 'marked'
import Overlay from './Overlay'

const creatingStrategy = {
  initialize() {
    core.initNewReport()
    this.updateMarkdown()
  },
  postReport() {
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

const updatingStrategy = {
  initialize() {
    this.loading = true
    core
      .findReport4Update(this.$route.params.user_id, this.$route.params.id)
      .then(() => {
        this.updateMarkdown()
        this.loading = false
      })
  },
  postReport() {
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

export default {
  name: 'BasicReportForm',
  components: { Overlay },
  data() {
    return {
      markdown: '',
      errors: [],
      state: core.state,
      loading: false
    }
  },
  created() {
    console.log('created')
    this.initialize()
  },
  methods: {
    initialize() {},

    doPostReport() {},

    postReport() {
      this.errors = []
      this.doPostReport.call(this).catch(error => {
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

<style scoped>
.editor {
  display: flex;
  position: relative;
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
  position: relative;
}
</style>
