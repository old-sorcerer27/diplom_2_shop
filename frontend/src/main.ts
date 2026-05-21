import './assets/main.css'
// import './assets/auth.css'
// import './assets/profile.css'

import { createApp} from 'vue'
import {currentRoute, render} from './utils/router'

import 'bootstrap/dist/css/bootstrap.min.css'
import 'bootstrap/dist/js/bootstrap.bundle.min.js'
import { HomePageStates, SetHomeState } from './utils/home'


export const API_URL = 'http://localhost:8080/api/v1'


export enum AsyncState {
  IDLE = 'IDLE',
  WAITING = 'WAITING',
  ERROR = 'ERROR',
  DONE = 'DONE',
}

SetHomeState(HomePageStates.PRODUCTS)
const app = createApp(render)
app.mount('#app')


 



