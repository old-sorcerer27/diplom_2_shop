<script setup lang="ts">
import {computed} from 'vue'
import { cart } from '../../../utils/cart'

import { products } from '@/utils/products'
import { format } from '@/utils/utils'
import { reviews } from '@/utils/reviews'
import ProductReviews from './ProductReviews.vue'
import ProductGallery from './gallery/ProductGallery.vue'


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
                  <ProductGallery
                    :images="products.product.images"
                    :product-name="products.product.name"
                  />
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
                                <span>{{ format.formatDate(products.product.created_at) }}</span>
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

.product-gallery {
    background: white;
    border-radius: 15px;
    padding: 1rem;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.main-image {
    width: 100%;
    height: 400px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #f8f9fa;
    border-radius: 12px;
    overflow: hidden;
    margin-bottom: 1rem;
}

.main-image img {
    max-width: 100%;
    max-height: 100%;
    object-fit: contain;
    transition: transform 0.3s;
}

.main-image img:hover {
    transform: scale(1.05);
}

.thumbnail-list {
    display: flex;
    gap: 0.5rem;
    overflow-x: auto;
    padding: 0.5rem 0;
}

.thumbnail {
    width: 80px;
    height: 80px;
    flex-shrink: 0;
    background: #f8f9fa;
    border-radius: 8px;
    overflow: hidden;
    cursor: pointer;
    border: 2px solid transparent;
    transition: all 0.3s;
}

.thumbnail:hover {
    transform: translateY(-2px);
}

.thumbnail.active {
    border-color: #667eea;
    box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

.thumbnail img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

/* Информационная карточка */
.product-info-card {
    background: white;
    border-radius: 15px;
    padding: 1.5rem;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
    height: 100%;
}

.product-title {
    font-size: 2rem;
    font-weight: 700;
    color: #333;
    margin-bottom: 1rem;
    line-height: 1.3;
}

/* Рейтинг */
.product-rating {
    display: flex;
    align-items: center;
    gap: 1rem;
    margin-bottom: 1rem;
    flex-wrap: wrap;
}

.stars {
    display: flex;
    gap: 4px;
}

.stars i {
    color: #ffc107;
    font-size: 1rem;
}

.reviews-count {
    color: #667eea;
    text-decoration: none;
    font-size: 0.9rem;
    transition: color 0.3s;
}

.reviews-count:hover {
    color: #5a67d8;
    text-decoration: underline;
}

/* Цена */
.product-price {
    font-size: 2.5rem;
    font-weight: bold;
    color: #28a745;
    margin: 1rem 0;
    display: inline-block;
    padding: 0.5rem 1rem;
    background: #d4edda;
    border-radius: 12px;
}

/* Мета-информация */
.product-meta {
    background: #f8f9fa;
    border-radius: 12px;
    padding: 1rem;
    margin: 1rem 0;
}

.meta-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.5rem 0;
    border-bottom: 1px solid #e0e0e0;
}

.meta-item:last-child {
    border-bottom: none;
}

.meta-label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-weight: 600;
    color: #555;
}

.meta-label i {
    color: #667eea;
    width: 20px;
}

.stock-badge {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    padding: 4px 12px;
    border-radius: 20px;
    font-size: 0.85rem;
    font-weight: 500;
}

.stock-in {
    background: #d4edda;
    color: #155724;
}

.stock-out {
    background: #f8d7da;
    color: #721c24;
}

/* Описание */
.product-description {
    margin: 1.5rem 0;
}

.product-description h5 {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 1.1rem;
    font-weight: 600;
    margin-bottom: 0.75rem;
    color: #333;
}

.product-description p {
    color: #666;
    line-height: 1.6;
    font-size: 0.95rem;
}

/* Выбор количества */
.quantity-selector {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin: 1rem 0;
    flex-wrap: wrap;
}

.quantity-btn {
    width: 36px;
    height: 36px;
    background: #f0f0f0;
    border: 1px solid #ddd;
    border-radius: 8px;
    font-size: 1.2rem;
    font-weight: bold;
    cursor: pointer;
    transition: all 0.3s;
}

.quantity-btn:hover {
    background: #667eea;
    color: white;
    border-color: #667eea;
}

.quantity-input {
    width: 60px;
    height: 36px;
    text-align: center;
    border: 1px solid #ddd;
    border-radius: 8px;
    font-size: 1rem;
}

.quantity-input::-webkit-inner-spin-button,
.quantity-input::-webkit-outer-spin-button {
    opacity: 1;
}

/* Кнопки покупки */
.buy-buttons {
    display: flex;
    gap: 1rem;
    margin-top: 1.5rem;
}

.btn-buy {
    flex: 1;
    padding: 0.875rem;
    font-size: 1rem;
    font-weight: 600;
    border: none;
    border-radius: 10px;
    cursor: pointer;
    transition: all 0.3s;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
}

.btn-buy:first-child {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
}

.btn-buy:first-child:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
}

.btn-buy:last-child {
    background: #28a745;
    color: white;
}

.btn-buy:last-child:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 5px 15px rgba(40, 167, 69, 0.4);
}

.btn-buy:disabled {
    background: #ccc;
    cursor: not-allowed;
    opacity: 0.7;
}

/* Секция отзывов */
#reviews {
    margin-top: 2rem;
    scroll-margin-top: 80px;
}

/* Анимация загрузки */
.loading {
    text-align: center;
    padding: 3rem;
}

.spinner {
    width: 50px;
    height: 50px;
    border: 3px solid #f3f3f3;
    border-top: 3px solid #667eea;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin: 0 auto;
}

@keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
}

/* Адаптивность */
@media (max-width: 992px) {
    .main-image {
        height: 350px;
    }
    
    .product-title {
        font-size: 1.6rem;
    }
    
    .product-price {
        font-size: 2rem;
    }
}

@media (max-width: 768px) {
    .container {
        padding: 0 15px;
    }
    
    .main-image {
        height: 300px;
    }
    
    .thumbnail {
        width: 60px;
        height: 60px;
    }
    
    .product-title {
        font-size: 1.4rem;
    }
    
    .product-price {
        font-size: 1.8rem;
    }
    
    .buy-buttons {
        flex-direction: column;
    }
    
    .meta-item {
        flex-direction: column;
        align-items: flex-start;
        gap: 0.5rem;
    }
}

@media (max-width: 576px) {
    .product-info-card {
        padding: 1rem;
    }
    
    .main-image {
        height: 250px;
    }
    
    .stars i {
        font-size: 0.85rem;
    }
}
</style>
