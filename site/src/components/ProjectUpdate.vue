<template>
    <div v-if="project">
        <h1>Update Project '{{ project.name }}'</h1>
        <h3>Repo Url: <input type="text" v-bind:value="project.repoUrl" readonly/></h3>
        <h3>Repo Path: <input type="text" v-bind:value="project.repoPath" readonly/></h3>
        <h3>Settings: <a v-on:click="newSetting()">New</a></h3>
        <div v-for="(key,val) in project.settings">
            <li>
                <input type="text" v-bind:value="key" readonly/>
                <input type="text" v-bind:value="val" readonly/>
            </li>
        </div>
    </div>
</template>

<script>
  export default {
    data () {
      return {}
    },
    computed: {
      project () {
        return this.$store.state.project
      }
    },
    methods: {
      newSetting () {
        this.$store.dispatch('UPDATE_SETTING', this.$route.params.name)
      }
    },
    created () {
      this.$store.dispatch('LOAD_PROJECT', this.$route.params.name)
    }
  }
</script>
