import Vue from 'vue'
import Router from 'vue-router'
import ProjectDetail from '@/components/ProjectDetail'
import ProjectList from '@/components/ProjectList'
import ProjectUpdate from '@/components/ProjectUpdate'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'ProjectList',
      component: ProjectList
    },
    {
      path: '/project/:name',
      name: 'ProjectDetail',
      component: ProjectDetail
    },
    {
      path: '/project/:name/update',
      name: 'ProjectUpdate',
      component: ProjectUpdate
    }
  ]
})
