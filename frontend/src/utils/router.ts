import {h} from 'vue'
import App from '../App.vue'
import Login from '@/components/auth/Login.vue'
import Register from '@/components/auth/Register.vue'
import NotFoundComponent from '@/components/NotFoundComponent.vue'
import Profile from '@/components/Profile/Profile.vue'
import AdminPanel from '@/components/Profile/admin/AdminPanel.vue'


export const currentRoute = {
  data: window.location.pathname
}

const routes: {[key: string]: typeof App} = {
  '/': App,
  '/login': Login,
  '/register': Register,
  '/profile': Profile,
  '/admin/products': AdminPanel,
  '/admin/products/new': AdminPanel,
  '/owner': Profile,
}

const CurrentComponent =  routes[currentRoute.data] || NotFoundComponent
 
export function navigateTo(path: string) {
  window.history.pushState(null, '', path)
  currentRoute.data = path
  window.location.reload();
}

export function render() {return h(CurrentComponent)}
