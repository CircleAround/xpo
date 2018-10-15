<template>
  <div class='reports'>
    <el-card class="box-card" v-for='(item, key , index) in list' v-bind:key="index">
      <div slot="header" class="clearfix card-header">
        <div class="user_name">
          <router-link :to="{ name:'UserPage', params: { author: item.author } }">
            <div class="nickname">{{item.authorNickname}}</div>
            <div class="name">{{item.author}}</div>
          </router-link>
        </div>
        <div class="card-header-optoins">
          <div class="date">
            <router-link :to="{ name:'ReportsYmd', params: { authorId: item.authorId, year: item.reportedAt.format('YYYY'), month: item.reportedAt.format('M'), day: item.reportedAt.format('DD') } }">
              <div class="month-day">
                <div class="month">{{item.reportedAt.format('M')}}</div>
                <div class="separator">/</div>
                <div class="day">{{item.reportedAt.format('DD')}}</div>
              </div>
            </router-link>
          </div>
          <el-dropdown class="card-menu">
            <el-button class="el-dropdown-link" icon="el-icon-arrow-down" circle></el-button>
            <el-dropdown-menu slot="dropdown" v-if="item.authorId == state.me.id">
              <router-link :to="{ name:'ReportEditForm', params: { authorId: item.authorId, id: item.id } }">
                <el-dropdown-item>Edit</el-dropdown-item>
              </router-link>
            </el-dropdown-menu>
          </el-dropdown>
        </div>
      </div>
      <div v-html="item.markdown()" class="markdown"></div>
      <div>
        <div class="updated-at">
          <router-link :to="{ name:'Report', params: { authorId: item.authorId, id: item.id } }">
            {{item.updatedAt.format('YYYY[/]MM[/]DD HH[:]mm[:]ss')}}
          </router-link>
        </div>
      </div>
    </el-card>
  </div>
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

<style scoped>
.box-card {
  margin-bottom: 10px;
}

.el-card__header {
  padding: 10px 20px;
}

.author-name {
  font-weight: bold;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-header-optoins {
  flex-basis: 90px;

  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-menu {
  flex-basis: 30px;
}

.date {
  text-align: center;
}

.date a {
  text-decoration: none;
}

.month-day {
  display: flex;
}

.year {
  font-size: 86%;
}

.updated-at {
  text-align: right;
  font-size: 90%;
}
</style>
