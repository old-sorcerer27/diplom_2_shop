import {reactive, computed } from 'vue'
import { API_URL } from '../main'
import { auth } from './auth'
import { currentHomeState, HomePageStates, SetHomeState } from './home'

export interface Product {
    ID: number;
    name: string;
    description: string;
    category: string;
    price: number;
    created_At: Date;
    stock: number;
    images?: string[];
    user_purchased?: boolean; 
    average_rating?: number;
}

export let products = reactive({
    products: <Product[]>([]),
    product: <Product | null>(null),
    currentImage: <string>('https://via.placeholder.com/500x400?text=No+Image'),
    loading: false,



    canLeaveReview: computed<boolean>((): boolean => {
        return auth.isAuthenticated() && Boolean(products.product?.user_purchased);
    }),

    loadProducts : async () => {
        try {
            const response = await fetch(`${API_URL}/products`);
            products.products = await response.json();
            products.loading = false;
            console.log('Товары загружены:', products.products);
        } catch (error) {
            console.error('Ошибка загрузки товаров:', error);
        }
    },

    loadProduct : async (productId: number) => {
        console.log(productId)
        if (!productId) {
            alert('Товар не найден');
            return;
        }

        try {
            const response = await fetch(`${API_URL}/products/${productId}`);

            if (!response.ok) throw new Error('Failed to load product');

            products.product = await response.json();

            // Настройка изображений
            products.currentImage  = products.product?.images?.[0] || 'https://via.placeholder.com/500x400?text=No+Image';

        } catch (error) {
            console.error('Error loading product:', error);
            alert('Ошибка загрузки товара');
        } finally {
            console.log(currentHomeState)
            SetHomeState(HomePageStates.PRODUCT)
            console.log(currentHomeState)
            products.loading = false;
        }
    },

    reloadProduct : () => {
        if (!products.product) return;
        products.loadProduct(products.product.ID);
    }

    
})