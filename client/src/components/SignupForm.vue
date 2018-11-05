<template>
  <div>
    <v-alert
      :value="true"
      type="success"
    >
      初めてのご利用の方ですね？ユーザー名とニックネームを入れてから進みましょう。
    </v-alert>
    <v-form ref="form" class="signup_form" :model="form" :rules="rules">
      <profile-form :errors="errors" :propErrors="propErrors"
       :name="form.name" :nickname="form.nickname" v-on:clicked-submit="postXUser">
      </profile-form>
    </v-form>
  </div>
</template>

<script>
import ProfileForm from './parts/ProfileForm'
import core from '../core'
import ErrorHandler from '../app/ErrorHandler'

class FormErrorHandler extends ErrorHandler {
  messageKeyByType(type, property) {
    if (type === 'duplicatedObject') {
      return `error.messages.duplicatedUser`
    }
    if (type === 'valueNotUnique') {
      return `error.messages.duplicatedUserName`
    }
    return `error.messages.${type}`
  }
}

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
      this.form.name = params.name
      this.form.nickname = params.nickname
      this.errors = []
      this.propErrors = {}
      core
        .postXUser(this.form.name, this.form.nickname)
        .then(() => {
          core.alert(
            `${
              this.form.nickname
            }さん、サインアップできました！楽しんでください！`,
            'success'
          )
          this.$router.push('/')
        })
        .catch(error => {
          new FormErrorHandler(this.$i18n).eachInResponse(
            error.response.data,
            (msg, type, property) => {
              if (type === 'duplicatedUser') {
                core.alert(msg, 'error')
                return this.$router.push('/')
              }

              if (!property) {
                return this.errors.push(msg)
              }

              if (!this.propErrors[property]) {
                this.$set(this.propErrors, property, [])
              }

              this.propErrors[property].push(msg)
            }
          )
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
