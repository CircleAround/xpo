<template>
  <div class='report_form'>
    <textarea v-model="report.content" class="textcontent" @keyup='updateMarkdown()'></textarea>
    <el-button type="success" icon="el-icon-check" circle @click='postReport()'></el-button>
    <div v-html="markdown"></div>
  </div>
</template>

<script>
import core from "../core"
import marked from 'marked'
export default {
  name: "report_form",
  data() {
    return {
      report: {
        content: "test",
        "content-type": "text/x-markdown"
      },
      markdown: ""
    };
  },
  methods: {
    postReport() {
      core.postReport(this.report).then((response)=>{
        
      })
    },
    updateMarkdown() {
      this.markdown = marked(this.report.content);
    }
  }
};
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style scoped>
  .textcontent{
    width: 100%;
    height: 10em;
  }
</style>
