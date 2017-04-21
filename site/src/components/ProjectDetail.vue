<template>
    <div>
        <h1>Project <span v-if="project.Name">'{{ project.Name }}'</span> Pending Updates</h1>
        <div v-if="plan.Diff">
            <ul v-for="mod in plan.Diff.Modules">
                <li v-for="(res, id) in mod.Resources">
                    <h3>{{ id }}</h3>
                    <p>{{ res }}</p>
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
    created () {
      this.$store.dispatch('LOAD_PROJECT', this.$route.params.name)
      this.$store.dispatch('GET_PLAN', this.$route.params.name)
    }
  }
</script>
