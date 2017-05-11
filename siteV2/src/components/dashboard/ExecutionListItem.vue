<template>
    <v-list-tile avatar>
        <v-list-tile-avatar>
          <v-icon large :class="icon.color">{{ icon.name }}</v-icon>
        </v-list-tile-avatar>
        <v-list-tile-content>
          <v-list-tile-title>terraform plan</v-list-tile-title>
          <v-list-tile-sub-title>1 min ago in prod</v-list-tile-sub-title>
        </v-list-tile-content>
    </v-list-tile>
</template>
<script>
    export default {
        name: 'execution-list-item',
        props: {
            command: String,
            status: String,
            time: String
        },
        computed: {
            icon () {
                let icon = {}
                switch(this.status) {
                    case 'in-progress':
                        icon.name = 'hourglass_empty'
                        icon.color = 'orange'
                    case 'successful':
                        icon.name = 'check_circle'
                        icon.color = 'green'
                    default:
                        icon.name = 'error'
                        icon.color = 'red'
                }
                icon.color = `${icon.color}--text text--darken-1`
                return icon
            }
        }
    }
</script>