import Vue from 'vue'
import Vuex from 'vuex'

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
    fetchProjects (ctx) {
      ctx.commit('projects', ['x', 'y', 'z'])
    }
  }
})
