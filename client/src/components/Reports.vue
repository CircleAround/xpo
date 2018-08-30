<template>
  <div class='reports'>
    <el-card class="box-card" v-for='item in list'>
      <div slot="header" class="clearfix">
        <strong>{{item.author}}</strong>
        <div style="float: right; padding: 3px 0; font-size: 90%" type="text">{{item.created_at.format('YYYY[/]MM[/]DD HH[:]mm[:]ss')}}</div>
      </div>
      <div v-html="item.markdown()" class="text item"></div>
    </el-card>
  </div>
</template>

<script>
import core from "../core";
export default {
  name: "reports",
  data() {
    return {
      msg: "About!!!!!",
      list: core.status.list
    };
  },
  created() {
    console.log("created");
    core.retriveReports().catch(function(error) {
      console.log(error)
    });

    if(core.posted) {
      this.$message({
        showClose: true,
        message: '投稿しました！',
        type: 'success',
        center: true
      });
      core.posted = false
    }
  },
  mounted() {
    console.log("mounted");
  },
  updated() {
    console.log("updated");
  },
  methods: {
    push() {}
  }
};
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style scoped>
.reports {
  margin: 0 auto;
  max-width: 960px;
}

.box-card {
  margin-bottom: 10px;
}
</style>
