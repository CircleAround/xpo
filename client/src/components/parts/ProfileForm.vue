<template>
  <v-container>
    <v-layout row wrap>
      <v-flex xs12 offset-sm3 sm6>

        <div class="errors" v-if="errors.length > 0">
          <div class="error" v-for='(item, key , index) in errors' v-bind:key="index">
            {{item}}
          </div>
        </div>

        <div>
          <v-text-field label="ユーザー名（半角英数小文字）"
            v-model="uname" id="input-name" minlength="3" counter="15" maxlength="15"
            :error-messages="propErrors.name" error-count="3"
          ></v-text-field>
          <v-text-field label="ニックネーム"
            v-model="nname" id="input-nickname" minlength="3" counter="24" maxlength="24"
            :error-messages="propErrors.nickname" error-count="3"
          ></v-text-field>
        </div>
        <div class="actions">
          <v-btn
            fab dark color="primary" 
            @click="clickedSubmit"
          >
            <v-icon dark>done</v-icon>
          </v-btn>
        </div>
      </v-flex>
    </v-layout>
  </v-container>
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
      nname: this.nickname,
      rules: {}
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
  padding-top: 1em;
  text-align: right;
}
</style>
