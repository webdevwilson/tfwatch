import api from '../../api'

const state = {
    configuration: []
}

const getters = {
  configuration: state => {
    return state.configuration
  }
}

const actions = {
  GET_CONFIGURATION (ctx) {
    api.ConfigurationResource.default.get()
      .then(response => {
        ctx.commit('configuration', response.body)
      })
  }
}

const mutations = {
  configuration (state, data) {
    state.configuration = data
  }
}

export default {
  state,
  getters,
  actions,
  mutations
}