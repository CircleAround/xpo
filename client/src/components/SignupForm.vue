<template>
  <div class="signup_form">
    <div class="errors" v-if="errors.length > 0">
      <div class="error" v-for='(item, key , index) in errors' v-bind:key="index">
        {{item}}
      </div>
    </div>
    <div>
      <el-input placeholder="ユーザー名（半角英数小文字）" v-model="name"></el-input>
      <div class="errors" v-if="propErrors.name">
        <div class="error" v-for="(item, key, index) in propErrors.name" v-bind:key="index">{{item}}</div>
      </div>

      <el-input placeholder="ニックネーム" v-model="nickname"></el-input>
      <div class="errors" v-if="propErrors.nickname">
        <div class="error" v-for="(item, key, index) in propErrors.nickname" v-bind:key="index">{{item}}</div>
      </div>
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
      propErrors: {},
      errors: [],
      name: '',
      nickname: ''
    }
  },
  methods: {
    postXUser() {
      this.errors = []
      this.propErrors = {}
      core
        .postXUser(this.name, this.nickname)
        .then(() => {
          this.$message({
            showClose: true,
            message: 'サインアップしました！',
            type: 'success',
            center: true
          })
          this.$router.push('/')
        })
        .catch(error => {
          core.eachResponseErrors(error, (msg, type, property) => {
            if (type === 'DuplicatedObjectError') {
              this.$message({
                showClose: true,
                message: 'すでにユーザー登録済みです',
                type: 'error',
                center: true
              })

              return this.$router.push('/')
            }

            if (type === 'ValueNotUniqueError') {
              if (!this.propErrors[property]) {
                this.$set(this.propErrors, property, [])
              }
              return this.propErrors[property].push(
                '既に取得されてしまったユーザー名です。別の名前にしましょう。'
              )
            }

            if (!property) {
              return this.errors.push(msg)
            }

            if (!this.propErrors[property]) {
              this.$set(this.propErrors, property, [])
            }

            this.propErrors[property].push(msg)
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

.error {
  color: red;
}
</style>
