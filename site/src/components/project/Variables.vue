<template>
    <div>
        <v-card>
            <v-card-row>
                <v-card-title>
                    <span class="grey--text">Update Variables</span>
                </v-card-title>
            </v-card-row>
            <v-card-text>
                <v-card-row>
                    <v-row>
                        <v-col xs4>
                            <v-text-field name="input-1-name" label="Name"></v-text-field>
                        </v-col>
                        <v-col xs8>
                            <v-text-field name="input-1-value" label="Value"></v-text-field>
                        </v-col>
                    </v-row>
                </v-card-row>
            </v-card-text>
            <v-card-row actions>
                <v-btn flat>Save</v-btn>
            </v-card-row>
        </v-card>
    </div>
</template>
<script>
export default {
    name: 'variables',
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