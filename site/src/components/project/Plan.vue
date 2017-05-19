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
        <v-card-text v-if="project.status == 'ok'">
            <p>The project is up-to-date.</p>
        </v-card-text>
        <v-card-text v-if="project.status == 'error'">
            <p class="red--text">This project contains errors. Please view the logs to get more information.</p>
        </v-card-text>
        <v-card-text v-if="project.status == 'pending'">
            <v-card>
                <v-toolbar>
                    <v-toolbar-title>Pending Updates</v-toolbar-title>
                    <v-spacer></v-spacer>
                    <v-icon>view_module</v-icon>
                </v-toolbar>
                <v-list two-line subheader>
                    <v-subheader inset>Resources</v-subheader>
                    <v-list-item v-for="item in project.pending_changes" v-bind:key="item.resource_id">
                        <v-list-tile avatar>
                            <v-list-tile-avatar>
                                <v-icon>check_circle</v-icon>
                            </v-list-tile-avatar>
                            <v-list-tile-content>
                            <v-list-tile-title>{{ item.resource_id }}</v-list-tile-title>
                            <v-list-tile-sub-title>{{ item.action }}</v-list-tile-sub-title>
                            </v-list-tile-content>
                            <v-list-tile-action>
                            <v-btn icon ripple>
                                <v-icon class="grey--text text--lighten-1">info</v-icon>
                            </v-btn>
                            </v-list-tile-action>
                        </v-list-tile>
                    </v-list-item>
                </v-list>
            </v-card>
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