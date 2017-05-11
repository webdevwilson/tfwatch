<template>
  <v-list two-line subheader>
    <v-subheader inset>Terraform Components</v-subheader>
    <v-list-item v-for="prj in projects" :key="prj.guid">
      <v-list-tile avatar router v-bind:to="{ name: 'Project', params: { guid: prj.guid }}">
        <v-list-tile-avatar>
          <v-icon v-if="prj.status == 'ok'" large class="green--text text--darken-1">check_circle</v-icon>
          <v-icon v-if="prj.status == 'pending'" large class="blue--text text--darken-1">info</v-icon>
          <v-icon v-if="prj.status == 'error'" large class="red--text text--darken-1">error</v-icon>
        </v-list-tile-avatar>
        <v-list-tile-content>
          <v-list-tile-title>{{prj.name}}</v-list-tile-title>
          <v-list-tile-sub-title>{{prj.planupdated}}</v-list-tile-sub-title>
        </v-list-tile-content>
      </v-list-tile>
    </v-list-item>
  </v-list>
</template>
<script>
  import { mapGetters } from 'vuex'
  let computed = mapGetters({
    'projects': 'projectList'
  })
  export default {
    name: 'project-list',
    computed: computed
  }
</script>
