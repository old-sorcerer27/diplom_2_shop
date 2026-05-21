import {reactive, computed } from 'vue'
import { auth } from './auth';
import { Product } from './products';

export const totalPrice = computed(() => {
  return cart.cart.reduce((sum: number, item: { product: Product; quantity: number }) => sum + (item.product?.price || 0) * item.quantity, 0)
})

export let cart = reactive({
    cart: <{ product: Product; quantity: number }[]>([]),
    loading: false,
    showModalOrder: false,

    loadCart() {
        const savedCart = localStorage.getItem('cart');
        if (savedCart) {
            cart.cart = JSON.parse(savedCart);
        }
    },

    saveCart() {
        localStorage.setItem('cart', JSON.stringify(this.cart));
    },

    updateQuantity(index: number, delta: number) {
        const cartItem = this.cart[index];
        if (!cartItem) return;

        const newQuantity = cartItem.quantity + delta;
        if (newQuantity > 0 && newQuantity <= cartItem.product.stock) {
           cartItem.quantity = newQuantity;
           this.saveCart();
        }
    },

    removeFromCart(index : number) {
        this.cart.splice(index, 1);
        this.saveCart();
    },

    addToCart(product: Product, quantity: number) {
        for (let item of this.cart) {
            if (item.product.ID === product.ID) {
                this.updateQuantity(this.cart.indexOf(item), quantity);
                return;
            }
        }
        this.cart.push({
           product: product,
           quantity: quantity
        });
       this.saveCart();
    },

    clearCart() {
        cart.cart = [];
        cart.saveCart();
    },

    submitOrder: async () => {
    try {
        const userData = localStorage.getItem('user');
        if (userData) {
            auth.getCurrentUser = JSON.parse(userData);
        }
    } catch (error) {
        console.error('Error loading user data:', error);
    } finally {
        auth.loading = false;
    }
    }

})