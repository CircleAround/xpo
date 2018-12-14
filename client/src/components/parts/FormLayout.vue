<template>
  <div class="form">
    <div v-if="state.me.id">
      <v-layout row wrap>
        <v-flex xs12 offset-sm3 sm6>
          <div class="errors" v-if="errors.length > 0">
            <div class="error" v-for="(item, key , index) in errors" v-bind:key="index">{{item}}</div>
          </div>

          <div slot="inputs" class="inputs"></div>

          <div slot="actions" class="actions"></div>
        </v-flex>
      </v-layout>
    </div>
  </div>
</template>

<script>
import core from '../../core'
import ErrorHandler from '../../app/ErrorHandler'

export default {
  name: 'FormLayout',
  data() {
    return {
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

    newErrorHandler() {
      return new ErrorHandler(this.$i18n)
    },

    beforeSave() {
      this.loading = true
      this.errors = []
      this.propErrors = {}
    },

    save() {
      this.beforeSave()
      this.doSave()
        .catch(error => {
          this.newErrorHandler().eachInResponse(
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
