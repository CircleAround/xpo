<template>
  <div class='report_form'>
    <textarea v-model="report.content"></textarea>
    <button type="button" @click='updateMarkdown()'>更新</button>
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
    updateMarkdown() {
      this.markdown = "aamakred";
    }
  }
};
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style scoped>
h1,
h2 {
  font-weight: normal;
}

a {
  color: #42b983;
}
</style>
