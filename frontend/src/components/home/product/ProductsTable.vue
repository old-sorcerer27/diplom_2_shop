<script setup lang="ts">
import {onMounted} from 'vue'
const API_URL = 'http://localhost:8080/api/v1'
import { products } from '@/utils/products'
import { cart } from '@/utils/cart'
import ProductPage from './ProductPage.vue'
import { CheckHomeState, HomePageStates } from '@/utils/home'

onMounted(() => {
    products.loadProducts();
})

</script>

<template>
    <div  class = "container mt-4">
            <div v-if="products.loading" class="loading">
                <div class="spinner"></div> 
                <p class="mt-3">Загрузка...</p>
            </div>
            <div v-else class="row">
            <div class="col-md-8">
                <h2>Товары</h2>
                <div class="row">
                    <div class="col-md-6 mb-4" v-for="product in products.products" :key="product.ID">
                        <div class="card h-100">
                            <div class="card-body">
                                <button @click="products.loadProduct(product.ID)">

                                    <h5 class="card-title">{{ product.name }}</h5>
                                    <p class="card-text">{{ product.description }}</p>
                                    <p class="price">{{ product.price }} ₽</p>
                                    <p class="stock">В наличии: {{ product.stock }}</p>
                                </button>
                                <button class="btn btn-primary" 
                                      @click="cart.addToCart(product, 1)"
                                       :disabled="product.stock === 0">
                                    В корзину
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>


<style>


</style>


<!-- <script lang="ts">
import { defineComponent, PropType, ref} from 'vue'
const API_URL = 'http://localhost:8080/api/v1'
export default defineComponent({
    name: "ProdutsTable",
    data() {
      return {
        products: [],
        loading: false,
        showcart: false
      }
    },
    methods: {
        const  loadProducts = async () => {
        try {
            const response = await fetch(`${API_URL}/products`);
            products.value = await response.json();
            loading.value = false;
            console.log('Товары загружены:', products.value);
        } catch (error) {
            console.error('Ошибка загрузки товаров:', error);
        }
        },
    }
})
</script> -->