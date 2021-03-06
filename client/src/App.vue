<template>
  <v-app id="app">
    <v-navigation-drawer
      v-if="isLoggedIn"
      v-model="drawer"
      fixed
      app
    >
      <v-list>
        <v-list-tile :to="{ name:'UserPage', params: { author: state.me.name } }">
          <v-list-tile-action>
            <v-icon>home</v-icon>
          </v-list-tile-action>
          <v-list-tile-title>MyPage</v-list-tile-title>
        </v-list-tile>
        <v-list-tile to='/users/me/edit'>
          <v-list-tile-action>
            <v-icon>person</v-icon>
          </v-list-tile-action>
          <v-list-tile-title>EditProfile</v-list-tile-title>
        </v-list-tile>
        <v-list-group
            :v-model="openedOption"
            :prepend-icon="openedOption ? 'keyboard_arrow_up' : 'keyboard_arrow_down'"
            append-icon=""
          >
            <v-list-tile slot="activator">
              <v-list-tile-content>
                <v-list-tile-title>Options</v-list-tile-title>
              </v-list-tile-content>
            </v-list-tile>
            <v-list-tile to='/about'>
              <v-list-tile-action>
                <v-icon>info</v-icon>
              </v-list-tile-action>
              <v-list-tile-title>About</v-list-tile-title>
            </v-list-tile>
            <v-list-tile v-bind:href="state.me.logoutUrl">
              <v-list-tile-action>
                <v-icon>power_settings_new</v-icon>
              </v-list-tile-action>
              <v-list-tile-title>Logout</v-list-tile-title>
            </v-list-tile>
        </v-list-group>
      </v-list>
    </v-navigation-drawer>

    <v-toolbar color="white" class="toolbar" app>
      <v-toolbar-side-icon @click.native="drawer = !drawer" v-if="isLoggedIn"></v-toolbar-side-icon>
      <v-toolbar-title>
        <h1 class="site_title"><router-link to='/'>TechLog</router-link></h1>
      </v-toolbar-title>

      <v-spacer></v-spacer>

      <v-toolbar-items>
        <template v-if="isLoggedIn">
          <v-btn flat :to="{ name:'UserPage', params: { author: state.me.name } }">
            <div class="user_name">
              <div class="nickname">{{state.me.nickname}}</div>
              <div class="name">{{state.me.name}}</div>
            </div>
          </v-btn>
          <v-btn flat color="primary" to='/reports/new'>
            <v-icon dark>edit</v-icon>
          </v-btn>
        </template>
        <template v-else>
          <v-btn v-bind:href="state.me.loginUrl" flat>
            Googleアカウントでログインして投稿する
          </v-btn>
          <v-btn flat to='/about'>
            <v-icon>info</v-icon>
          </v-btn>
        </template>
      </v-toolbar-items>
    </v-toolbar>

    <v-content>
      <v-container fluid>
        <v-alert
          :value="true"
          :dismissible="true"
          :type="alert.type"
          v-for="(alert, key , index) in state.alerts" v-bind:key="index"
        >
          {{alert.message}}
        </v-alert>

        <router-view></router-view>
      </v-container>
    </v-content>

    <v-content>
      <v-container fluid>
        <v-layout row wrap>
          <v-flex xs12>
            このサービスはまだ実験的に作成しています。利用規約など気になる人はまだ使っちゃダメです！
            データの永続性なども保証していません。
            サービス提供者の<a :href="consts.TWITTER_URL" target="_blank">ms2sato</a>は一切の責任を負えませんので承知の上でご利用ください。
            ソースコードは<a :href="consts.REPOSITORY_URL" target="_blank">公開</a>されていますので、気になる方はご確認の上でご利用ください。
          </v-flex>
        </v-layout>
      </v-container>
    </v-content>
  </v-app>
</template>

<script>
import core from './core'
import consts from './consts'
export default {
  name: 'App',
  data() {
    return {
      drawer: null,
      openedOption: false,
      state: core.state,
      consts: consts
    }
  },
  computed: {
    isLoggedIn() {
      return core.isLoggedIn()
    }
  },
  created() {
    console.log('app created')
  }
}
</script>

<style lang="scss">
@import '@/scss/main.scss';

.application {
  font-family: 游ゴシック体, 'Yu Gothic', YuGothic, 'ヒラギノ角ゴシック Pro',
    'Hiragino Kaku Gothic Pro', メイリオ, Meiryo, Osaka, 'ＭＳ Ｐゴシック',
    'MS PGothic', sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: $text-color;
}

.toolbar {
  .site_title {
    display: inline-block;
    font-family: 'Avenir', Helvetica, Arial, sans-serif;
  }

  .site_title a {
    text-decoration: none;
  }

  .user_name {
    text-align: right;
    max-width: 200px;
    text-transform: none;

    a,
    a:visited {
      color: $text-color;
    }
  }
}

.user_name {
  word-break: break-all;
  padding: 0 5px;

  a {
    text-decoration: none;
  }

  .nickname,
  .name {
    text-overflow: ellipsis;
    overflow: hidden;
    white-space: nowrap;
    word-break: break-all;
    font-size: 90%;
  }

  .nickname {
    font-weight: bold;
  }
}
</style>
