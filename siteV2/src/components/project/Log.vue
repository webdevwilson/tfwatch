<template>
    <v-card>
        <v-card-row>
            <v-card-title>
                <span class="grey--text">View Logs</span>
                <v-spacer></v-spacer>
                <v-btn icon="icon" class="green--text">
                    <v-icon>cached</v-icon>
                </v-btn>
            </v-card-title>
        </v-card-row>
        <v-card-text>View Recent Tasks</v-card-text>
    </v-card>
</template>
<script>
export default {
    name: 'log',
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