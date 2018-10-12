<template>
  <div class="report_form">
    <div v-if="state.me.id">
      <div class="editor">
        <overlay :visible="loading"></overlay>
        <textarea v-model="state.newReport.content" v-focus class="textcontent" @keydown.meta.enter="postReport()" @keyup='updateMarkdown()' :placeholder="$t('ui.placeholder.markdown')"></textarea>
        <div class="preview markdown" v-html="markdown"></div>
      </div>
      <div class="errors" v-if="errors.length > 0">
        <div class="error" v-for='(item, key , index) in errors' v-bind:key="index">
          {{item}}
        </div>
      </div>
      <div class="actions">
        <overlay :visible="loading"></overlay>
        <el-popover
          placement="top-start"
          title="Help"
          width="200"
          trigger="click"
        >
          <p>
            {{$t('ui.help.markdown')}}
          </p>
          <div>
            <div>Mac: {{$t('ui.help.shortcutkey.post.mac')}}</div>
            <div>Windows: {{$t('ui.help.shortcutkey.post.win')}}</div>
          </div>
          <el-button icon="el-icon-question" circle slot="reference"></el-button>
        </el-popover>
        <el-button type="success" icon="el-icon-check" circle @click='postReport()' title="Mac: Command + Enter, Win: Control + Enter"></el-button>
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
import ErrorHandler from '../../app/ErrorHandler'

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
    this.initialize()
  },
  methods: {
    initialize() {
      throw new Error('Unimplemented: initialize')
    },

    doPostReport() {
      throw new Error('Unimplemented: doPostReport')
    },

    postReport() {
      this.loading = true
      this.errors = []
      this.doPostReport()
        .catch(error => {
          new ErrorHandler(this.$i18n).eachInResponse(
            error.response.data,
            (msg, type, property) => {
              this.errors.push(msg)
            }
          )
        })
        .finally(() => {
          this.loading = false
        })
    },
    updateMarkdown() {
      this.markdown = marked(this.state.newReport.content)
    }
  }
}
</script>

<style lang="scss" scoped>
@import '@/scss/mixin.scss';
.editor {
  display: flex;
  position: relative;

  textarea {
    @include placeholder() {
      color: #aaa;
    }
  }
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
