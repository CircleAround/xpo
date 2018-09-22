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
      <profile-form :errors="errors" :propErrors="propErrors" :name="form.name" :nickname="form.nickname" v-on:clicked-submit="postXUser">
      </profile-form>
    </el-form>
  </div>
</template>

<script>
import ProfileForm from './parts/ProfileForm'
import core from '../core'
export default {
  name: 'signup_form',
  components: { ProfileForm },
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
    postXUser(params) {
      this.form.name = params.username
      this.form.nickname = params.nickname
      this.errors = []
      this.propErrors = {}
      core
        .postXUser(this.form.name, this.form.nickname)
        .then(() => {
          this.$message({
            showClose: true,
            message: `${
              this.form.nickname
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
.el-row {
  margin-bottom: 18px;
}
</style>
