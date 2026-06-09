import './assets/main.css'

import { createApp} from 'vue'
import {currentRoute, render} from './utils/router'


// import 'bootstrap/dist/js/bootstrap.bundle.min.js'
// import 'bootstrap/dist/css/bootstrap.min.css'
// import 'bootstrap-vue/dist/bootstrap-vue.css'
import 'bootstrap/dist/css/bootstrap.css';

import { HomePageStates, SetHomeState } from './utils/home'


// import * as dotenv from 'dotenv';

// dotenv.config();
// export const API_URL = process.env.VITE_API_URL || 'http://localhost:8080/api/v1';

export const API_URL = 'http://localhost:8080/api/v1';

export enum AsyncState {
  IDLE = 'IDLE',
  WAITING = 'WAITING',
  ERROR = 'ERROR',
  DONE = 'DONE',
}

SetHomeState(HomePageStates.PRODUCTS)
const app = createApp(render)
app.mount('#app')