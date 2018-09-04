<template>
  <div class="signup_form">
    <div class="errors" v-if="errors.length > 0">
      <div class="error" v-for='(item, key , index) in errors' v-bind:key="index">
        {{item}}
      </div>
    </div>
    <div>
      <el-input placeholder="ユーザー名（半角英数小文字）" v-model="name"></el-input>
      <el-input placeholder="ニックネーム" v-model="nickname"></el-input>
    </div>
    <div class="actions">
      <el-button type="success" icon="el-icon-check" circle @click='postXUser()'></el-button>
    </div>
  </div>
</template>

<script>
import core from '../core'
export default {
  name: 'signup_form',
  data() {
    return {
      errors: [],
      name: '',
      nickname: ''
    }
  },
  methods: {
    postXUser() {
      this.errors = []
      core.postXUser(this.name, this.nickname).catch(error => {
        core.getMessagesOfValidationError(error).forEach(msg => {
          this.errors.push(msg)
        })
      })
    }
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style scoped>
.editor {
  display: flex;
}

.textcontent {
  flex: 1;
  width: 100%;
  height: 10em;

  font-size: 16px;
  border: solid 1px #ccc;
}

.preview {
  flex: 1;
  margin: 0 10px;
  padding: 10px;
  background-color: #f2f2f2;
}

.actions {
  padding: 5px;
  text-align: right;
}

.errors {
  color: red;
}
</style>
