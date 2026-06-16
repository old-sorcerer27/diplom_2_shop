<script setup lang="ts">

import { HomePageStates, SetHomeState } from '@/utils/home';
import { cart, OrderData, totalPrice } from '../../utils/cart'
import { auth } from '@/utils/auth';
import { ref, watch } from 'vue';

const form = ref<OrderData>({
  email: '',
  name: '',
  comment: '',
  phone: ''
});

const loadUserData = () => {
  if (auth.isAuthenticated()) {    
        const user = auth.getCurrentUser()
        if (user) {
            form.value.name = user.username || ''
            form.value.email = user.email || ''
            form.value.phone = user.phone || ''
            form.value.comment = ''
        }
        console.log('Loaded user data:', form.value)
    }
}


async function handleSubmit() {   
    loadUserData() 
    try {
      const orderPayload = {
        items: cart.cart.map(item => ({
          product_id: item.product.ID,
          quantity: item.quantity
        })),
        customer_name: form.value.name,
        customer_email: form.value.email,
        comment: form.value.comment,
        customer_phone: form.value.phone,
        total: totalPrice.value,
        is_authenticated: auth.isAuthenticated()
      }
      const response = await cart.createOrder(orderPayload)
      if (response.success) {
        alert(`Заказ #${response.order?.id} успешно оформлен!\nНа почту отправлено подтверждение.`)
      } else {
        throw new Error(response.error || 'Ошибка оформления заказа')
      }
    } catch (error: any) {
        console.error('Order error:', error)
        alert(error.message || 'Ошибка оформления заказа. Попробуйте позже.')
    } 
}

// watch(() => auth.getCurrentUser(), () => {
//     loadUserData()
// })

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
            <button class="btn btn-success" @click="cart.checkoutCart()" data-bs-toggle="modal" data-bs-target="#checkoutModal">
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

    <div class="modal fade" id="checkoutModal" tabindex="-1">
        <div class="modal-dialog">
            <div class="modal-content">
                <template v-if="auth.isAuthenticated()">
                    <div class="modal-header">
                        <h5 class="modal-title">Спасибо за покупку!</h5>
                        <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                    </div>
                    <div class="modal-body">
                        <p>Ваш заказ успешно оформлен. Мы свяжемся с вами для подтверждения деталей.</p>
                    </div>
                </template>
                <template v-else>
                    <div class="modal-header">
                        <h5 class="modal-title">Войдите в аккаунт</h5>
                        <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                    </div>
                     <form id="registerForm" @submit.prevent="handleSubmit">
                    <div class="form-group">
                        <i class="fas fa-user"></i>
                        <input type="text" 
                               class="form-control" 
                               id="name" 
                               v-model="form.name"
                               placeholder="Ваше имя" 
                               required>
                    </div>
                    <div class="form-group">
                        <i class="fas fa-envelope"></i>
                        <input type="email" 
                               class="form-control" 
                               id="email" 
                               v-model="form.email"
                               placeholder="Email" 
                               required>
                    </div>
                    <div class="form-group password-field">
                        <i class="fas fa-lock"></i>
                        <input type="phone" 
                               class="form-control" 
                               id="phone" 
                               v-model="form.phone"
                               placeholder="Номер телефона" 
                               required>
                    </div>                  
                     <div class="form-group password-field">
                        <i class="fas fa-lock"></i>
                        <input type="comment" 
                               class="form-control" 
                               id="comment" 
                               v-model="form.comment"
                               placeholder="Комментарий к заказу" 
                               required>
                    </div>           
                    <button type="submit" class="btn btn-primary" id="submitBtn">
                        <span id="submitText">Оформить заказ</span>
                    </button>
                </form>

                </template>
                <div class="modal-footer">
                    <button type="button" class="btn btn-primary" data-bs-dismiss="modal" @click="handleSubmit()">
                        Закрыть
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>