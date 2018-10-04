<template>
  <div>
    <el-row>
      <el-col :sm="{span:6, offset: 4}" class="errors" v-if="errors.length > 0">
        <div class="error" v-for='(item, key , index) in errors' v-bind:key="index">
          {{item}}
        </div>
      </el-col>
    </el-row>
    <div>
      <el-row>
        <el-col  :sm="{span:6, offset: 4}">
          <el-form-item label="ユーザー名（半角英数小文字）" for="input-name"></el-form-item>
        </el-col>
        <el-col  :sm="10">
          <el-input placeholder="ユーザー名（半角英数小文字）" v-model="uname" id="input-name" minlength="3" maxlength="15"></el-input>
          <div class="errors" v-if="propErrors.name">
            <div class="error" v-for="(item, key, index) in propErrors.name" v-bind:key="index">{{item}}</div>
          </div>
        </el-col>
      </el-row>
      <el-row>
        <el-col  :sm="{span:6, offset: 4}">
          <el-form-item label="ニックネーム" for="input-nickname"></el-form-item>
        </el-col>
        <el-col  :sm="10">
          <el-input placeholder="ニックネーム" v-model="nname" id="input-nickname" minlength="3" maxlength="24"></el-input>
          <div class="errors" v-if="propErrors.nickname">
            <div class="error" v-for="(item, key, index) in propErrors.nickname" v-bind:key="index">{{item}}</div>
          </div>
        </el-col>
      </el-row>
    </div>
    <el-row class="actions">
      <el-col :sm="{offset:4, span: 16}">
        <el-button type="primary" icon="el-icon-check" @click='clickedSubmit'></el-button>
      </el-col>
    </el-row>
  </div>
</template>
<script>
export default {
  name: 'ProfileForm',
  props: {
    propErrors: Object,
    errors: Array,
    name: String,
    nickname: String
  },
  data() {
    return {
      uname: this.name,
      nname: this.nickname
    }
  },
  methods: {
    clickedSubmit: function() {
      this.$emit('clicked-submit', {
        name: this.uname,
        nickname: this.nname
      })
    }
  }
}
</script>
<style lang="scss" scoped>
.actions {
  text-align: right;
}
</style>
