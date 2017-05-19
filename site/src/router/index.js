import Vue from 'vue'
import Router from 'vue-router'

// top-level view imports
import About from '../components/About.vue'
import Dashboard from '../components/Dashboard.vue'
import Project from '../components/Project.vue'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Dashboard',
      component: Dashboard
    },
    {
      path: '/project/:guid',
      name: 'Project',
      component: Project,
      props: true
    },
    {
      path: '/about',
      name: 'About',
      component: About
    }
  ]
})
