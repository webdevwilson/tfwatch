import Vue from 'vue'
import Vuex from 'vuex'
import {PlanResource, ProjectResource} from '../resources'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    project: {},
    projects: [],
    plan: {}
  },
  mutations: {
    project (state, data) {
      state.project = data
    },
    projects (state, data) {
      state.projects = data
    },
    plan (state, data) {
      state.plan = data
    }
  },
  actions: {
    APPLY (ctx, projectId) {
    },
    LOAD_PROJECT (ctx, projectId) {
      ProjectResource.get({id: projectId})
        .then(response => {
          ctx.commit('project', response.body)
        })
    },
    LOAD_PROJECT_LIST (ctx) {
      ProjectResource.get()
        .then(response => {
          ctx.commit('projects', response.body)
        })
    },
    GET_PLAN (ctx, projectId) {
      PlanResource.get({
        id: projectId
      }).then(response => {
        ctx.commit('plan', response.body)
      })
    }
  }
})
