<script setup lang="ts">
import {reactive} from 'vue'
const API_URL = 'http://localhost:8080/api/v1'
import { cart, totalPrice } from '../utils/cart'
import { products } from '@/utils/products'
import { Modal } from 'bootstrap'

const orderForm = reactive({
  name: '',
  email: ''
})

const  submitOrder = async () => {
    if (!orderForm.name || !orderForm.email) {
        alert('Пожалуйста, заполните все поля');
        return;
    }
    
    const orderData = {
        customer_name: orderForm.name,
        customer_email: orderForm.email,
        items: cart.cart.map(item => ({
            product_id: item.product.ID,
            quantity: item.quantity
        }))
    };

    console.log('Отправляем данные заказа:', orderData);
    
    try {
        const response = await fetch(`${API_URL}/orders`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(orderData)
        });
   

        if (response.ok) {
            const result = await response.json();
            alert(`Заказ #${result.order_id} оформлен! Сумма: ${result.total} ₽`);
            cart.cart = [];
            cart.saveCart();
            orderForm.email, orderForm.name = '';
            cart.showModalOrder = false;
            products.loadProducts();
        } else {
            alert('Ошибка оформления заказа');
        }
    } catch (error) {
        console.error('Ошибка:', error);
        alert('Ошибка оформления заказа');
    }
}

</script>

<template>
    <div class="modal fade" id="checkoutModal" tabindex="-1">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title">Оформление заказа</h5>
                    <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                </div>
                <div class="modal-body">
                    <form onsubmit.prevent="submitOrder()">
                        <div class="mb-3">
                            <label>Ваше имя</label>
                            <input type="text" class="form-control" v-model="orderForm.name" required>
                        </div>
                        <div class="mb-3">
                            <label>Email</label>
                            <input type="email" class="form-control" v-model="orderForm.email" required>
                        </div>
                        <div class="mb-3">
                            <label>Сумма заказа</label>
                            <input type="text" class="form-control" :value="totalPrice + ' ₽'" disabled>
                        </div>
                        <button type="submit" class="btn btn-primary">Подтвердить заказ</button>
                    </form>
                </div>
            </div>
        </div>
    </div>
</template>
