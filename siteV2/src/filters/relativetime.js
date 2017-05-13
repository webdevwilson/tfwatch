import Vue from 'vue'
import moment from 'moment'

// Register a global custom directive called v-focus
Vue.filter('relativeTime', function(val) {
    return moment(val).fromNow()
})