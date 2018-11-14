<template>
  <v-item-group multiple>
    <v-item
      v-for="lng in languages.slice(0, 10)"
      :key="lng.name"
    >
      <v-chip
        slot-scope="{ active, toggle }"
        :selected="active"
      >
        <router-link :to="{ name:'LanguagePage', params: { language: lng.name } }">
            {{ lng.name }}
        </router-link>
      </v-chip>
    </v-item>
  </v-item-group>
</template>

<script>
import core from '../../core'

export default {
  name: 'TopLanguagesPanel',
  data() {
    return {
      languages: core.state.languages
    }
  },
  created() {
    core.getAllLanguages()
  },
  computed: {
    topLanguages() {
      return this.languages.slice(0, 10)
    }
  }
}
</script>
