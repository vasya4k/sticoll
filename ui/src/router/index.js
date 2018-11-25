import Vue from 'vue'
import Router from 'vue-router'

// Containers
const DefaultContainer = () => import('@/containers/DefaultContainer')

// Views - Devices
const Devices = () => import('@/views/Devices')
const Newdevice = () => import('@/views/Newdevice')

Vue.use(Router)

export default new Router({
  mode: 'hash', // https://router.vuejs.org/api/#mode
  linkActiveClass: 'open active',
  scrollBehavior: () => ({ y: 0 }),
  routes: [
    {
      path: '/',
      redirect: '/devices',
      name: 'Home',
      component: DefaultContainer,
      props: true,
      children: [
        {
          path: 'devices',
          name: 'Devices',
          component: Devices,
          props: true
        },
        {
          path: 'newdevice',
          name: 'Newdevice',
          component: Newdevice,
          props: true,                  
        }      
      ]
    }
  ]
})
