<template>
  <div id="app">
    <el-container>
      <el-header>
        <h1 class="site_title"><router-link to='/'>TechLog</router-link></h1>
        <ul class="topmenu">
          <template v-if="isLoggedIn">
            <li class="user_name">
              <div class="nickname">{{state.me.nickname}}</div>
              <div class="name">{{state.me.name}}</div>
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
        このサービスはまだ実験的に作成しています。利用規約など気になる人はまだ使っちゃダメです！
        データの永続性なども保証していません。
        サービス提供者の<a :href="consts.TWITTER_URL" target="_blank">ms2sato</a>は一切の責任を負えませんので承知の上でご利用ください。
        ソースコードは<a :href="consts.REPOSITORY_URL" target="_blank">公開</a>されていますので、気になる方はご確認の上でご利用ください。
      </el-footer>
    </el-container>
  </div>
</template>

<script>
import core from './core'
import consts from './consts'
export default {
  name: 'App',
  data() {
    return {
      state: core.state,
      consts: consts
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

.el-header {
  display: flex;
  justify-content: space-between;

  .site_title {
    display: inline-block;
  }

  .site_title a {
    text-decoration: none;
  }

  .topmenu {
    list-style: none;
    display: flex;
    padding:0;
    justify-content: flex-end;
  }

  .user_name {
    word-break: break-all;
    padding: 0 5px;
    font-weight: bold;
    text-align: right;

    width: 200px;

    .nickname, .name {
      text-overflow: ellipsis;
      overflow: hidden;
      white-space: nowrap;
      word-break: break-all;
      font-size: 90%;
    }
  }
}

@media screen and (max-width: 480px) {
  .el-main {
    padding-left: 2px;
    padding-right: 2px;
  }

  .el-header {
    padding-left: 2px;
    padding-right: 2px;

    .user_name {
        width: 100px;
    }
  }
}

</style>
