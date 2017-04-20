import Vue from 'vue'
import Vuex from 'vuex'
import {ProjectResource} from '../resources'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    projects: []
  },
  mutations: {
    projects (state, data) {
      state.projects = data
    }
  },
  actions: {
    LOAD_PROJECT_LIST (ctx) {
      ProjectResource.get()
        .then(response => {
          ctx.commit('projects', response.body)
        })
    }
  }
})
