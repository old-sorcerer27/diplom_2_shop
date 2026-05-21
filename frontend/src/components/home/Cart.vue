<script setup lang="ts">

import { CheckHomeState, HomePageStates, SetHomeState } from '@/utils/home';
import { auth } from '../../utils/auth'
import { cart, totalPrice } from '../../utils/cart'

</script>

<template>
    <div>
        <h2>Корзина</h2>
        <div v-if="cart.cart.length === 0" class="alert alert-info">
            Корзина пуста
            <button class="btn btn-secondary ms-2" @click="SetHomeState(HomePageStates.PRODUCTS)">
                Продолжить покупки
            </button>
        </div>
        <div v-else>
            <div class="card mb-3" v-for="(item, index) in cart.cart" :key="index">
                <div class="card-body">
                    <h5>{{ item.product.name }}</h5>
                    <p>Цена: {{ item.product.price }} ₽</p>
                    <p>Количество: 
                        <button class="btn btn-sm btn-secondary" @click="cart.updateQuantity(index, -1)">-</button>
                        {{ item.quantity }}
                        <button class="btn btn-sm btn-secondary" @click="cart.updateQuantity(index, 1)">+</button>
                    </p>
                    <p>Итого: {{ item.product.price * item.quantity }} ₽</p>
                    <button class="btn btn-danger btn-sm" @click="cart.removeFromCart(index)">Удалить</button>
                </div>
            </div>
            <div class="alert alert-success">
                <strong>Общая сумма: {{ totalPrice }} ₽</strong>
            </div>
            <button class="btn btn-success" @click="auth.checkouth()" data-bs-toggle="modal" data-bs-target="#checkoutModal">
                Оформить заказ
            </button>
            <button class="btn btn-secondary ms-2" @click="SetHomeState(HomePageStates.PRODUCTS)">
                Продолжить покупки
            </button>
            <button class="btn btn-danger ms-2" @click="cart.clearCart">
                Очистить корзину
            </button>
        </div>
    </div>
</template>