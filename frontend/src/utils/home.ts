import { ref } from 'vue'

export const enum HomePageStates {
  LOADING = 'LOADING',
  PRODUCTS = 'PRODUCTS',
  PRODUCT = 'PRODUCT',
  CART = 'CART',
  ERROR = 'ERROR',
}

export let currentHomeState = ref(HomePageStates.PRODUCTS);
let previousHomeState: HomePageStates = HomePageStates.LOADING;


export function CheckHomeState (state: HomePageStates) {
    return currentHomeState.value === state
}

export function SetHomeState (state: HomePageStates) {
    previousHomeState = currentHomeState.value
    currentHomeState.value  = state
}