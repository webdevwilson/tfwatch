const state = {
    projects: [
          {
            "guid":"17aadb28-1a90-428c-553a-60604297421d",
            "name":"lambda_monitor_queues",
            "status": "ok",
            "planupdated": "12 mins ago"
          },{
            "guid":"2460ae5a-e7cb-4083-406c-f447aad10a0b",
            "name":"thrx",
            "status": "ok",
            "planupdated": "2 mins ago"
          },{
            "guid":"41ae4c4e-a51d-4b49-6898-f4304d7e7977",
            "name":"monitor_directconnect",
            "status": "ok",
            "planupdated": "4 mins ago"
          },{
            "guid":"43db6ca7-4060-43e4-648b-7ecd0676ca4c",
            "name":"cloudwatch_nonprod",
            "status": "error",
            "planupdated": "< 1 min ago"
          },{
            "guid":"45db5cd0-ffc5-4ef0-42ad-aa0fc0ed34b4",
            "name":"infra_vpc",
            "status": "pending",
            "planupdated": "2 hours ago"
          },{
            "guid":"473071d1-3d4e-42f1-4cb2-7ad53dd3a78f",
            "name":"cloudwatch_prod",
            "status": "pending",
            "planupdated": "27 mins ago"
          },{
            "guid":"86ae9284-946e-4d36-5b84-d20ac2ae1bf3",
            "name":"lambda_manage_maintenance_page",
            "status": "ok",
            "planupdated": "2 days ago"
         }
      ]
}

const getters = {
  projectList: state => {
    return state.projects
  },
  project: (state, getters) => (guid) => {
    return state.projects.find(prj => prj.guid == guid)
  }
}

const actions = {}

const mutations = {}

export default {
  state,
  getters,
  actions,
  mutations
}