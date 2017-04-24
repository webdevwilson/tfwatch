<template>
    <div>
        <h1>Project <span v-if="project.name">'{{ project.name }}'</span> Pending Updates</h1>
        <button v-on:click="apply()">Apply</button>
        <div v-if="plan.resources">
            <ul>
                <li v-for="res in plan.resources" class="resource">
                  <span v-bind:class="res.action">{{ res.name }} - {{ res.action }}</span>
                </li>
            </ul>
        </div>
        <router-link :to="{ name: 'ProjectList' }">Return to list</router-link>
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
      },
      plan () {
        return this.$store.state.plan
      }
    },
    methods: {
      apply () {
        this.$store.dispatch('APPLY', this.$route.params.name)
      }
    },
    created () {
      this.$store.dispatch('LOAD_PROJECT', this.$route.params.name)
      this.$store.dispatch('GET_PLAN', this.$route.params.name)
    }
  }
</script>

<style>
  .resource { font-weight: bold; }
  .Create { color: green; }
  .Update { color: orange; }
  .Destroy { color: red; }
</style>
