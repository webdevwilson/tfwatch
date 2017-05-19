import Vue from 'vue'
import Vuex from 'vuex'
import project from './modules/project'
import configuration from './modules/configuration'

Vue.use(Vuex)

const actions = {}
const getters = {}
const mutations = {}

export default new Vuex.Store({
    actions,
    getters,
    mutations,
    modules: {
        project,
        configuration
    }
})