<template>
  <form-layout>
    <template slot="inputs">
      <v-text-field label="プロジェクト名（半角英数小文字）"
        v-model="project.name" id="input-name" minlength="3" counter="15" maxlength="15"
        :error-messages="propErrors.name" error-count="3"
      ></v-text-field>
      <v-textarea label="詳細"
        v-model="project.description" id="input-description" minlength="3" counter="140" maxlength="140"
        :error-messages="propErrors.nickname" error-count="3"
      ></v-textarea>
    </template>

    <template slot="actions">
      <overlay :visible="loading"></overlay>

      <v-btn fab dark color="primary" @click='save()'>
        <v-icon dark>done</v-icon>
      </v-btn>
    </template>
  </form-layout>
</template>

<script>
import core from '../../core'
import Overlay from './Overlay'
import ErrorHandler from '../../app/ErrorHandler'
import FormLayout from './FormLayout'

export default {
  extends: FormLayout,
  name: 'BasicReportForm',
  components: { Overlay, FormLayout },
  data() {
    return {
      project: {},
      propErrors: {},
      errors: [],
      state: core.state,
      loading: false
    }
  },
  created: async function() {
    this.initialize()
  },
  methods: {
    initialize() {
      throw new Error('Unimplemented: initialize')
    },

    doSave() {
      throw new Error('Unimplemented: doSave')
    },

    save() {
      this.loading = true
      this.errors = []
      this.doSave()
        .catch(error => {
          new ErrorHandler(this.$i18n).eachInResponse(
            error.response.data,
            (msg, type, property) => {
              this.errors.push(msg)
            }
          )
        })
        .finally(() => {
          this.loading = false
        })
    }
  }
}
</script>

<style lang="scss" scoped>
@import '@/scss/mixin.scss';
.actions {
  padding: 5px 0;
  text-align: right;
  position: relative;
}
</style>
