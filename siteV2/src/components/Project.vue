<template>
    <v-tabs grow scroll-bars class="elevation-3">
        <v-card class="secondary">
            <v-card-text>
                <v-card-row>
                    <v-card-title>{{project.name}}</v-card-title>
                    <v-spacer></v-spacer>
                    <v-icon v-if="project.status == 'ok'" large class="green--text text--darken-1">check_circle</v-icon>
                    <v-icon v-if="project.status == 'pending'" large class="blue--text text--darken-1">info</v-icon>
                    <v-icon v-if="project.status == 'error'" large class="red--text text--darken-1">error</v-icon>
                </v-card-row>
            </v-card-text>
        </v-card>
        <v-tab-item href="#tab-variables" slot="activators">Variables</v-tab-item>
        <v-tab-item href="#tab-plan" slot="activators">Plan</v-tab-item>
        <v-tab-item href="#tab-resources" slot="activators">Resources</v-tab-item>
        <v-tab-item href="#tab-permissions" slot="activators">Permissions</v-tab-item>
        <v-tab-item href="#tab-log" slot="activators">Log</v-tab-item>
        <v-tab-content id="tab-variables" slot="content">
            <div class="ma-4">
                <variables></variables>
            </div>
        </v-tab-content>
        <v-tab-content id="tab-plan" slot="content">
            <div class="ma-4">
                <plan :guid="guid"></plan>
            </div>
        </v-tab-content>
        <v-tab-content id="tab-resources" slot="content">
            <div class="ma-4">
                <resources :guid="guid"></resources>
            </div>
        </v-tab-content>
        <v-tab-content id="tab-log" slot="content">
            <div class="ma-4">
                <log :guid="guid"></log>
            </div>
        </v-tab-content>
        <v-tab-content id="tab-permissions" slot="content">
            <div class="ma-4">
                <permissions :guid="guid"></permissions>
            </div>
        </v-tab-content>
    </v-tabs>
</template>
<script>
import Log from './project/Log.vue'
import Permissions from './project/Permissions.vue'
import Plan from './project/Plan.vue'
import Resources from './project/Resources.vue'
import Variables from './project/Variables.vue'

export default {
    props: {
        guid: String
    },
    computed: {
        project () {
            return this.$store.getters.project(this.guid)
        }
    },
    components: {
        Log,
        Permissions,
        Plan,
        Resources,
        Variables
    }
}
</script>