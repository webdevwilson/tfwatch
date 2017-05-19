import Vue from 'vue'
import VueResource from 'vue-resource'

//Vue.http.options.crossOrigin = true
//Vue.http.options.credentials = true

Vue.use(VueResource)

export default {
    ProjectResource: require('./project'),
    ConfigurationResource: require('./configuration')
}