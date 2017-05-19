<template>
    <div id="getting-started">
        <h2>Welcome!</h2>
        <v-row>
            <v-col xs6>
                <v-card class="mr-2">
                    <v-toolbar class="grey darken-2">
                        <v-toolbar-title>Projects</v-toolbar-title>
                        <v-spacer></v-spacer>
                    </v-toolbar>
                    <project-list></project-list>
                </v-card>
            </v-col>
            <v-col xs6>
                <v-card class="ml-2">
                    <v-toolbar class="grey darken-2">
                        <v-toolbar-title>Executions</v-toolbar-title>
                        <v-spacer></v-spacer>
                    </v-toolbar>
                    <execution-list></execution-list>
                </v-card>
            </v-col>
        </v-row>
    </div>
</template>
<script>
    import ExecutionList from './dashboard/ExecutionList.vue'
    import ProjectList from './dashboard/ProjectList.vue'
    export default {
        components: {
            ExecutionList,
            ProjectList
        }
    }
</script>