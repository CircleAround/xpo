<template>
  <div>
    <el-form ref="form" class="signup_form" :model="form" :rules="rules">
      <el-row>
        <el-col :sm="{offset:4, span: 16}">
          <el-alert
            title="初めてのご利用の方ですね？ユーザー名とニックネームを入れてから進みましょう。"
            type="success">
          </el-alert>
        </el-col>
      </el-row>
      <div class="errors" v-if="errors.length > 0">
        <div class="error" v-for='(item, key , index) in errors' v-bind:key="index">
          {{item}}
        </div>
      </div>
      <div>
        <el-row>
          <el-col  :sm="{span:6, offset: 4}">
            <el-form-item label="ユーザー名（半角英数小文字）" for="input-name"></el-form-item>
          </el-col>
          <el-col  :sm="10">
            <el-input placeholder="ユーザー名（半角英数小文字）" v-model="form.name" id="input-name"></el-input>
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
            <el-input placeholder="ニックネーム" v-model="form.nickname" id="input-nickname"></el-input>
            <div class="errors" v-if="propErrors.nickname">
              <div class="error" v-for="(item, key, index) in propErrors.nickname" v-bind:key="index">{{item}}</div>
            </div>
          </el-col>
        </el-row>
      </div>
      <el-row class="actions">
        <el-col :sm="{offset:4, span: 16}">
          <el-button type="primary" icon="el-icon-check" @click='postXUser()'></el-button>
        </el-col>
      </el-row>
    </el-form>
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
      form: {
        name: '',
        nickname: ''
      },
      rules: {}
    }
  },
  methods: {
    postXUser() {
      this.errors = []
      this.propErrors = {}
      core
        .postXUser(this.form.name, this.form.nickname)
        .then(() => {
          this.$message({
            showClose: true,
            message: `${
              this.nickname
            }さん、サインアップできました！楽しんでください！`,
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

<style lang="scss" scoped>
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
  text-align: right;
}

.el-row {
  margin-bottom: 18px;
}
</style>
