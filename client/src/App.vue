<template>
  <div id="app">
    <el-container>
      <el-header>
        <h1 class="site_title"><router-link to='/'>TechLog</router-link></h1>
        <ul class="topmenu">
          <li>
            <router-link to='/about'><el-button icon="el-icon-info" circle></el-button></router-link>
          </li>
          <li>
            <router-link to='/report'><el-button type="primary" icon="el-icon-edit" circle></el-button></router-link>
          </li>
        </ul>

      </el-header>
      <el-main>
        <router-view/>
      </el-main>
      <el-footer>
        利用規約とか気になる人はまだ使っちゃダメです！まだそんなに責任持てるレベルまで仕上げていません。
        データの永続性とか保証しません。
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
  color: #2c3e50;
}

.el-container {
  margin: 0 auto;
  max-width: 960px;
}

.site_title {
  display: inline-block;
}

.site_title a {
  text-decoration: none;
}

.topmenu {
  list-style: none;
  display: flex;
  float: right;
}
</style>
