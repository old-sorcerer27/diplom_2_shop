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
.container {
    max-width: 1400px;
    margin: 0 auto;
    padding: 0 20px;
}

/* Заголовок раздела */
h2 {
    font-size: 1.8rem;
    font-weight: 600;
    margin-bottom: 1.5rem;
    color: #333;
    position: relative;
    display: inline-block;
}

h2::after {
    content: '';
    position: absolute;
    bottom: -8px;
    left: 0;
    width: 60px;
    height: 3px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 3px;
}

/* Сетка товаров */
.row {
    display: flex;
    flex-wrap: wrap;
    margin: 0 -15px;
}

.col-md-6 {
    flex: 0 0 50%;
    max-width: 50%;
    padding: 0 15px;
}

/* Карточка товара */
.card {
    background: white;
    border: none;
    border-radius: 15px;
    overflow: hidden;
    transition: all 0.3s ease;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
    height: 100%;
    display: flex;
    flex-direction: column;
}

.card:hover {
    transform: translateY(-5px);
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
}

.card-body {
    padding: 1.5rem;
    flex: 1;
    display: flex;
    flex-direction: column;
}

/* Кнопка-ссылка на карточку товара */
.card-body button:first-child {
    background: none;
    border: none;
    text-align: left;
    cursor: pointer;
    padding: 0;
    margin: 0;
    width: 100%;
    flex: 1;
}

/* Заголовок товара */
.card-title {
    font-size: 1.2rem;
    font-weight: 600;
    color: #333;
    margin-bottom: 0.75rem;
    transition: color 0.3s;
    line-height: 1.4;
}

button:hover .card-title {
    color: #667eea;
}

/* Описание товара */
.card-text {
    color: #666;
    font-size: 0.9rem;
    line-height: 1.5;
    margin-bottom: 1rem;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
}

/* Цена */
.price {
    font-size: 1.5rem;
    font-weight: bold;
    color: #28a745;
    margin: 0.5rem 0;
}

/* Информация о наличии */
.stock {
    font-size: 0.85rem;
    color: #6c757d;
    margin-bottom: 1rem;
    display: flex;
    align-items: center;
    gap: 0.5rem;
}

.stock::before {
    content: '●';
    font-size: 0.8rem;
}

.stock:has(+ .btn-primary:not(:disabled))::before {
    color: #28a745;
}

/* Кнопка добавления в корзину */
.btn-primary {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border: none;
    padding: 0.75rem 1.5rem;
    font-size: 0.95rem;
    font-weight: 600;
    border-radius: 10px;
    color: white;
    cursor: pointer;
    transition: all 0.3s;
    width: 100%;
}

.btn-primary:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
}

.btn-primary:disabled {
    background: #ccc;
    cursor: not-allowed;
    opacity: 0.7;
}

/* Адаптивность */
@media (max-width: 768px) {
    .col-md-6 {
        flex: 0 0 100%;
        max-width: 100%;
    }
    
    .card-title {
        font-size: 1.1rem;
    }
    
    .price {
        font-size: 1.3rem;
    }
}

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