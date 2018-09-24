<template>
  <div>
    <el-form ref="form" class="profile_edit_form" :model="form" :rules="rules">
      <profile-form :errors="errors" :propErrors="propErrors"
        :name="form.name" :nickname="form.nickname"
        v-on:clicked-submit="updateXUser"
      >
      </profile-form>
    </el-form>
  </div>
</template>

<script>
import ProfileForm from './parts/ProfileForm'
import core from '../core'
export default {
  name: 'profile_edit_form',
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
  mounted() {
    this.form.name = core.state.me.name
    this.form.nickname = core.state.me.nickname
  },
  methods: {
    updateXUser(params) {
      this.form.name = params.name
      this.form.nickname = params.nickname
      this.errors = []
      this.propErrors = {}
      core
        .updateXUser(this.form.name, this.form.nickname)
        .then(() => {
          this.$message({
            showClose: true,
            message: 'プロフィールを更新しました',
            type: 'success',
            center: true
          })
          this.$router.push('/')
        })
        .catch(error => {
          core.eachResponseErrors(error, (msg, type, property) => {
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
