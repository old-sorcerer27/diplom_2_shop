<script setup lang="ts">
import {getCurrentInstance, onMounted } from 'vue'

import Navbar from './components/navbar/Navbar.vue' 

import ProductsTable from './components/home/product/ProductsTable.vue' 
import ProductPage from './components/home/product/ProductPage.vue' 
import Cart from './components/home/Cart.vue'

import {auth} from './utils/auth'
import {cart} from './utils/cart'
import { CheckHomeState, HomePageStates } from './utils/home'


onMounted(() => {
    cart.loadCart();
})

console.log('App component loaded. Auth state:', auth.isAuthenticated(), 'Cart items:', cart.cart.length);
</script>

<template>
  <main>
  <div id="app">
        <Navbar/>
        <div class="container mt-4">
            <div>
                <div class="col-md-4">
                    <div v-if="CheckHomeState(HomePageStates.PRODUCTS)" class="card">
                        <div class="card-header">
                            <h5>Информация</h5>
                        </div>
                        <div class="card-body">
                            <p>Добро пожаловать в наш магазин!</p>
                            <p>Выберите товары и оформите заказ.</p>
                        </div>
                    </div>
                </div>
                <div v-if="CheckHomeState(HomePageStates.PRODUCTS)">
                    <ProductsTable/>
                </div>
                <div v-if="CheckHomeState(HomePageStates.PRODUCT)">
                    <ProductPage/>
                </div>
                <div v-if="CheckHomeState(HomePageStates.CART)">
                    <Cart/>
                </div>
            </div>
        </div>
    </div>
  </main>
</template>

<style>


</style>
