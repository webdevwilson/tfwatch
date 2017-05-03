<template>
    <div>
        <h1>Project List</h1>
        <input type="button" v-on:click="fetchProjects" value="Refresh" />
        <ul>
          <li v-for="project in projects">
            <router-link :to="{ name: 'ProjectDetail', params: { name: project.guid }}">{{ project.name }}</router-link>
            <router-link :to="{ name: 'ProjectUpdate', params: { name: project.guid }}">edit</router-link>
          </li>
        </ul>
    </div>
</template>

<script>
  export default {
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
    }
  }
</script>
