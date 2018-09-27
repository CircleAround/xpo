<template>
  <div id="app">
    <el-container>
      <el-header>
        <h1 class="site_title"><router-link to='/'>TechLog</router-link></h1>
        <ul class="topmenu">
          <template v-if="isLoggedIn">
            <li class="user_name">
              {{state.me.nickname}}[{{state.me.name}}]
            </li>
            <li>
              <router-link to='/reports/new'><el-button type="primary" icon="el-icon-edit" circle></el-button></router-link>
            </li>
            <li>
              <el-dropdown>
                <el-button class="el-dropdown-link" icon="el-icon-arrow-down" circle></el-button>
                <el-dropdown-menu slot="dropdown">
                  <router-link to='/users/me/edit'>
                    <el-dropdown-item >
                      EditProfile
                    </el-dropdown-item>
                  </router-link>
                  <router-link to='/about'>
                    <el-dropdown-item >
                      About
                    </el-dropdown-item>
                  </router-link>
                  <a v-bind:href="state.me.logoutUrl">
                    <el-dropdown-item divided>Logout</el-dropdown-item>
                  </a>
                </el-dropdown-menu>
              </el-dropdown>
            </li>
          </template>
          <template v-else>
            <li>
              <a v-bind:href="state.me.loginUrl">Googleアカウントでログイン</a>して投稿する
              <router-link to='/about'><el-button icon="el-icon-info" circle></el-button></router-link>
            </li>
          </template>
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
import core from './core'
export default {
  name: 'App',
  data() {
    return {
      state: core.state
    }
  },
  computed: {
    isLoggedIn() {
      return core.isLoggedIn()
    }
  }
}
</script>

<style lang="scss">
@import '@/scss/main.scss';

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

.el-dropdown-menu a {
  text-decoration: none;
}

.el-dropdown-menu a:visited {
  text-decoration: none;
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

.user_name {
  word-break: break-all;

  display: flex;
  align-items: center;
  padding: 0 5px;
  font-weight: bold;
}
</style>
