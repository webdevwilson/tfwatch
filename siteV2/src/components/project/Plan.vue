<template>
    <v-card>
        <v-card-row>
            <v-card-title>
                <span class="grey--text">Plan</span>
                <v-spacer></v-spacer>
                <v-btn icon="icon" class="green--text">
                    <v-icon>cached</v-icon>
                </v-btn>
            </v-card-title>
        </v-card-row>
        <v-card-text>
            <p>The project is up-to-date.</p>
        </v-card-text>
        <v-card-row actions>
            <v-btn flat>Apply Changes</v-btn>
        </v-card-row>
    </v-card>
</template>
<script>
export default {
    name: 'plan',
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