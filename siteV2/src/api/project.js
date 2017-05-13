import Vue from 'vue'

export default Vue.resource('api/projects{/id}')/*, {}, {
    list: {
        method: 'GET',
        url: 'api/projects'
    },
    tfplan: {
        method: 'GET',
        url: 'api/projects{/id}/tfplan'
    },
    executions: {
        method: 'GET',
        url: 'api/projects{/id}/tfplan'
    }
})
*/
