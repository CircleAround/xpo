<template>
  <v-container fluid grid-list-lg pl-0 pr-0>
  <v-layout row wrap>
    <v-flex xs12 v-for='(item, key , index) in list' v-bind:key="index">
      <v-card class="box-card">
        <v-card-title primary-title class="clearfix card-header pt-2 pb-2">

          <router-link :to="{ name:'UserPage', params: { author: item.author } }">
            <div class="user_name headline">
              <div class="nickname">{{item.authorNickname}}</div>
              <div class="name">{{item.author}}</div>
            </div>
          </router-link>

          <div class="languages" v-if="item.languages">
            <v-chip v-for='(lng, k, i) in item.languages' v-bind:key="i">
              <router-link :to="{ name:'LanguagePage', params: { language: lng } }">
                {{ lng }}
              </router-link>
            </v-chip>
          </div>

          <v-spacer></v-spacer>

          <div class="date">
            <router-link :to="{ name:'UserReportsYmd', params: { authorId: item.authorId, year: item.reportedAt.format('YYYY'), month: item.reportedAt.format('M'), day: item.reportedAt.format('DD') } }">
              <div class="month-day">
                <div class="month">{{item.reportedAt.format('M')}}</div>
                <div class="separator">/</div>
                <div class="day">{{item.reportedAt.format('DD')}}</div>
              </div>
            </router-link>
          </div>

          <v-menu bottom left>
            <v-btn
              slot="activator"
              icon
            >
              <v-icon>more_vert</v-icon>
            </v-btn>

            <v-list>
              <v-list-tile>
                <router-link :to="{ name:'Report', params: { authorId: item.authorId, id: item.id } }">
                  <v-list-tile-title>Show</v-list-tile-title>
                </router-link>
              </v-list-tile>
              <v-list-tile>
                <router-link v-if="item.authorId == state.me.id" :to="{ name:'ReportEditForm', params: { authorId: item.authorId, id: item.id } }">
                  <v-list-tile-title>Edit</v-list-tile-title>
                </router-link>
              </v-list-tile>
            </v-list>
          </v-menu>
          </v-card-title>

          <v-container fluid>
            <v-layout row wrap>
              <v-flex xs12>
                <div v-html="item.markdown()" class="markdown"></div>
              </v-flex>
            </v-layout>
          </v-container>

          <v-container pt-2 fluid>
            <v-layout row wrap>
              <v-flex xs12>
                <div class="updated-at">
                  <router-link :to="{ name:'Report', params: { authorId: item.authorId, id: item.id } }">
                    {{item.updatedAt.format('YYYY[/]MM[/]DD HH[:]mm[:]ss')}}
                  </router-link>
                </div>
              </v-flex>
            </v-layout>
          </v-container>
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import core from '../../core'
export default {
  name: 'ReportsPanel',
  data() {
    return {
      state: core.state,
      list: []
    }
  },
  props: {
    reports: Array
  },
  watch: {
    reports(val) {
      this.list = val
    }
  },
  created() {
    console.log('created')
  }
}
</script>

<style lang="scss" scoped>
.user_name.headline {
  .nickname {
    font-size: 14px !important;
    line-height: 1.2 !important;
  }

  .name {
    font-size: 12px !important;
    line-height: 1.2 !important;
  }
}

.card-header {
  border-bottom: solid 1px #eee;
}

.date {
  text-align: center;

  a {
    text-decoration: none;
  }
}

.month-day {
  display: flex;
}

.year {
  font-size: 86%;
}

.updated-at {
  text-align: right;
}

.u  pdated-at a {
  text-decoration: none;
}
</style>
