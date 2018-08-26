<template>
  <div id="app">
    <el-container>
      <el-header>
        X報
        <router-link to='/'>ホーム</router-link>
        <router-link to='/about'>About</router-link>
        <router-link to='/report'>報</router-link>
      </el-header>
      <el-main>
        <router-view/>
      </el-main>  
      <el-footer>
        footer
      </el-footer>  
    </el-container>
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
