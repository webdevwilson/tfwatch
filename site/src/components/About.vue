<template>
  <v-list two-line subheader>
    <v-subheader inset>System Configuration</v-subheader>
    <v-list-item v-for="c in config" :key="config.id">
      <v-list-tile>
        <v-list-tile-content>
          <v-list-tile-title>{{c.name}}</v-list-tile-title>
          <v-list-tile-sub-title>{{c.value}}</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>
    </v-list-item>
  </v-list>
</template>

<script>
  import { mapGetters } from 'vuex'
  let computed = mapGetters({
    'config': 'configuration'
  })
  export default {
    name: 'about',
    data () {
      return {}
    },
    computed: computed,
    created () {
      this.$store.dispatch('GET_CONFIGURATION')
    }
  }
</script>
