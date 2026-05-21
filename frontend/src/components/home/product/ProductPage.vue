<script setup lang="ts">
import {computed} from 'vue'
import { cart } from '../../../utils/cart'

import { products } from '@/utils/products'
import { format } from '@/utils/utils'
import { reviews } from '@/utils/reviews'
import ProductReviews from './ProductReviews.vue'


const totalPrice = computed(() => {
    return cart.cart.reduce((total, item) => total + (item.product.price * item.quantity), 0);
})

let quantity = 1;

const decreaseQuantity = () => {
    if (quantity > 1) {
        quantity--;
    }
}
    
const increaseQuantity = () => {
    if (quantity > 1) {
        quantity++;
    }
}

const declensionReviews = (count: number) => {
    const lastDigit = count % 10;
    const lastTwo = count % 100;
    
    if (lastTwo >= 11 && lastTwo <= 19) return 'отзывов';
    if (lastDigit === 1) return 'отзыв';
    if (lastDigit >= 2 && lastDigit <= 4) return 'отзыва';
    return 'отзывов';
}

const buyNow = async () => {

}

</script>

<template>
    <div class="container mt-4" v-if="products.product">
            <div class="row">
                <!-- Левая колонка - Галерея -->
                <div class="col-lg-6 mb-4">
                    <div class="product-gallery">
                        <div class="main-image">
                            <img :src="products.currentImage" :alt="products.product.name">
                        </div>
                        <div class="thumbnail-list">
                            <div v-for="(img, idx) in products.product.images" :key="idx"
                                 :class="['thumbnail', { active: products.currentImage === img }]"
                                 @click="products.currentImage = img">
                                <img :src="img" :alt="products.product.name + ' ' + idx">
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Правая колонка - Информация -->
                <div class="col-lg-6 mb-4">
                    <div class="product-info-card">
                        <h1 class="product-title">{{ products.product.name }}</h1>
                        
                        <div class="product-rating">
                            <div class="stars">
                                <i v-for="i in 5" :key="i" 
                                   :class="i <= Math.round(products.product.average_rating || 0) ? 'fas fa-star' : 'far fa-star'"></i>
                            </div>
                            <a href="#reviews" class="reviews-count">
                                {{ reviews.Reviews.length }} {{ declensionReviews(reviews.Reviews.length) }}
                            </a>
                        </div>
                        
                        <div class="product-price">
                            {{ format.formatPrice(products.product.price) }} ₽
                        </div>
                        
                        <div class="product-meta">
                            <div class="meta-item">
                                <span class="meta-label"><i class="fas fa-tag"></i> Категория</span>
                                <span>{{ products.product.category || 'Без категории' }}</span>
                            </div>
                            <div class="meta-item">
                                <span class="meta-label"><i class="fas fa-box"></i> Наличие</span>
                                <span>
                                    <span :class="['stock-badge', products.product.stock > 0 ? 'stock-in' : 'stock-out']">
                                        <i :class="products.product.stock > 0 ? 'fas fa-check-circle' : 'fas fa-times-circle'"></i>
                                        {{ products.product.stock > 0 ? 'В наличии: ' + products.product.stock + ' шт.' : 'Нет в наличии' }}
                                    </span>
                                </span>
                            </div>
                            <div class="meta-item">
                                <span class="meta-label"><i class="fas fa-calendar"></i> Добавлен</span>
                                <span>{{ format.formatDate(products.product.created_At) }}</span>
                            </div>
                        </div>
                        
                        <div class="product-description">
                            <h5><i class="fas fa-align-left"></i> Описание</h5>
                            <p>{{ products.product.description }}</p>
                        </div>
                        
                        <div class="quantity-selector" v-if="products.product.stock > 0">
                            <button class="quantity-btn" @click="decreaseQuantity">-</button>
                            <input type="number" class="quantity-input" v-model.number="quantity" min="1" :max="products.product.stock">
                            <button class="quantity-btn" @click="increaseQuantity">+</button>
                            <span class="text-muted">доступно {{ products.product.stock }} шт.</span>
                        </div>
                        
                        <div class="buy-buttons">
                            <button class="btn btn-primary btn-buy" 
                                    @click="cart.addToCart(products.product, quantity)"
                                    :disabled="products.product.stock === 0">
                                <i class="fas fa-shopping-cart"></i> В корзину
                            </button>
                            <button class="btn btn-success btn-buy" 
                                    @click="buyNow"
                                    :disabled="products.product.stock === 0">
                                <i class="fas fa-bolt"></i> Купить сейчас
                            </button>
                        </div>
                    </div>
                </div>
            </div>
            <ProductReviews :productId="products.product.ID" id="reviews"/>
        </div> 
</template>


<style>


</style>
