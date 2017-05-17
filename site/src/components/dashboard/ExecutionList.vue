<template>
  <v-list two-line subheader>
    <v-subheader inset>Commands</v-subheader>
    <v-list-item>
      <execution-list-item successful command="terraform plans" time="1 min ago in prod"></execution-list-item>
    </v-list-item>
    <v-list-item>
      <v-list-tile avatar>
        <v-list-tile-avatar>
          <v-icon large class="green--text text--darken-1">check_circle</v-icon>
        </v-list-tile-avatar>
        <v-list-tile-content>
          <v-list-tile-title>terraform plan</v-list-tile-title>
          <v-list-tile-sub-title>1 min ago in prod</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>
    </v-list-item>
    <v-list-item>
      <v-list-tile avatar>
        <v-list-tile-avatar>
          <v-icon large class="green--text text--darken-1">check_circle</v-icon>
        </v-list-tile-avatar>
        <v-list-tile-content>
          <v-list-tile-title>terraform apply</v-list-tile-title>
          <v-list-tile-sub-title>3 mins ago in nonprod</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>
    </v-list-item>
    <v-list-item>
      <v-list-tile avatar>
        <v-list-tile-avatar>
          <v-icon large class="red--text text--darken-1">error</v-icon>
        </v-list-tile-avatar>
        <v-list-tile-content>
          <v-list-tile-title>git pull origin master</v-list-tile-title>
          <v-list-tile-sub-title>12 mins ago</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>
    </v-list-item>
    <v-list-item>
      <v-list-tile avatar>
        <v-list-tile-avatar>
          <v-icon large class="green--text text--darken-1">check_circle</v-icon>
        </v-list-tile-avatar>
        <v-list-tile-content>
          <v-list-tile-title>terraform plan -o terraform.tfplan</v-list-tile-title>
          <v-list-tile-sub-title>8 mins ago in prod</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>
    </v-list-item>
  </v-list>
</template>

<script>
  import ExecutionListItem from './ExecutionListItem.vue'

  export default {
    name: 'execution-list',
    data () {
      return {}
    },
    methods: {
      fetchProjects () {
        this.$store.dispatch('LOAD_PROJECT_LIST')
      }
    },
    computed: {
      projects () {
        return this.$store.state.projects
      }
    },
    created () {
      this.fetchProjects()
    },
    components: {
      ExecutionListItem
    }
  }
</script>
