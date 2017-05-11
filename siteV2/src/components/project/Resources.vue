<template>
    <v-card>
        <v-card-row>
            <v-card-title>
                <span class="grey--text">Project Resources</span>
            </v-card-title>
        </v-card-row>
        <v-card-text>Display list of resources or graph</v-card-text>
    </v-card>
</template>
<script>
export default {
    name: 'resources',
    props: {
        guid: String
    },
    computed: {
        project () {
            return this.$store.getters.project(this.guid)
        }
    }
}
</script>