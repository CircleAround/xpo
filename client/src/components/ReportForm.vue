<template>
  <div class='report_form'>
    <textarea v-model="report.content"></textarea>
    <el-button type="success" icon="el-icon-check" circle @click='postReport()'></el-button>
    <div>{{markdown}}</div>
  </div>
</template>

<script>
import core from "../core";
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
        alert('POOST!')
      })
    },
    updateMarkdown() {
      this.markdown = "aamakred";
    }
  }
};
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style scoped>
</style>
