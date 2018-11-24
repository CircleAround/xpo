<template>
  <div
    class="project_form"
  >
    <div v-if="state.me.id">
      <div class="editor">
      </div>
      <div class="errors" v-if="errors.length > 0">
        <div class="error" v-for='(item, key , index) in errors' v-bind:key="index">
          {{item}}
        </div>
      </div>
      <div class="actions">
        <overlay :visible="loading"></overlay>

        <v-btn fab dark color="primary" @click='saveReport()'>
          <v-icon dark>done</v-icon>
        </v-btn>
      </div>
    </div>
    <div v-if="!state.me.id">
      ログインすると使えます
    </div>
  </div>
</template>

<script>
import core from '../../core'
import Overlay from './Overlay'
import ErrorHandler from '../../app/ErrorHandler'

export default {
  name: 'BasicReportForm',
  components: { Overlay },
  data() {
    return {
      project: {},
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

    doSaveReport() {
      throw new Error('Unimplemented: doSaveReport')
    },

    saveReport() {
      this.loading = true
      this.errors = []
      this.doSaveReport()
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
