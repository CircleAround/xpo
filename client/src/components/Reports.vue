<template>
  <div class='reports'>
    <ul>
      <li v-for='item in list'>
        {{item.content}}
      </li>
    </ul>
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
      alert("fail");
    });
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
h1,
h2 {
  font-weight: normal;
}

a {
  color: #42b983;
}
</style>
