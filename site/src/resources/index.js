import Vue from 'vue'
import VueResource from 'vue-resource'

Vue.use(VueResource)

Vue.http.options.crossOrigin = true
Vue.http.options.credentials = true

export const ProjectResource = Vue.resource('http://localhost:3000/api/projects{/id}')
export const PlanResource = Vue.resource('http://localhost:3000/api/projects/{id}/tfplan')
