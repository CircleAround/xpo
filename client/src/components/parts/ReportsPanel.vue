<template>
  <div class='reports'>
    <el-card class="box-card" v-for='(item, key , index) in list' v-bind:key="index">
      <div slot="header" class="clearfix card-header">
        <div class="author-name">{{item.author}}</div>
        <div class="card-header-optoins">
          <div style="font-size: 90%" type="text">
            {{item.created_at.format('YYYY[/]MM[/]DD HH[:]mm[:]ss')}}
          </div>
          <el-dropdown class="card-menu">
            <el-button class="el-dropdown-link" icon="el-icon-arrow-down" circle></el-button>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item>
                <router-link :to="{ name:'ReportEditForm', params: { authorId: item.authorId, id: item.id } }">Edit</router-link>
              </el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </div>
      </div>
      <div v-html="item.markdown()" class="markdown"></div>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'ReportsPanel',
  data() {
    return {
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
  }
}
</script>

<style scoped>
.box-card {
  margin-bottom: 10px;
}

.el-card__header {
  padding: 12px 20px;
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
  flex-basis: 190px;

  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-menu {
  flex-basis: 30px;
}
</style>
