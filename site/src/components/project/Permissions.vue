<template>
    <v-card>
        <v-card-row>
            <v-card-title>
                <span class="grey--text">Permissions</span>
            </v-card-title>
        </v-card-row>
        <v-card-text>
            <v-subheader>Who has permissions to view and apply on this project?</v-subheader>
            <v-list three-line subheader>

            </v-list>
        </v-card-text>
    </v-card>
</template>
<script>
export default {
    name: 'permissions',
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