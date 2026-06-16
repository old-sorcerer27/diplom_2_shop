import {reactive, computed, ref } from 'vue'
import { API_URL } from '../main'
import { auth } from './auth'
import { currentHomeState, HomePageStates, SetHomeState } from './home'
import { images } from './image';
import { api } from './api';


export interface Product {
    ID: number;
    name: string;
    description: string;
    category: string;
    image_url: string
    thumbnail_url: string
    medium_url: string
    gallery: string[]
    price: number;
    created_at: Date;
    stock: number;
    images?: string[];
    user_purchased?: boolean; 
    average_rating?: number;
}

export let products = reactive({
    products: <Product[]>([]),
    product: <Product | null>(null),
    currentImage: ref(''),  
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
            if (products.product?.image_url) {
                products.currentImage = images.getImageUrl(products.product.image_url, 'original')
            } else if (products.product?.gallery && products.product.gallery.length > 0 && products.product.gallery[0]) {
                products.currentImage = images.getImageUrl(products.product.gallery[0], 'original')
            }

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
    },

    createProduct : async (productData: Partial<Product>) => {
        if (auth.getUserRole() !== 'admin') {
            throw new Error('Access denied');
        }
        const token = localStorage.getItem('token');
        try {
            const response = await api.post('/admin/products', productData, token || undefined);
            return response;
        }
        catch (error) {
            console.error('Error creating product:', error);
            throw error;
        }
    },

    updateProduct : async (productId: number, productData: Partial<Product>) => {
        if (auth.getUserRole() !== 'admin') {
            throw new Error('Access denied');
        }
        const token = localStorage.getItem('token');
        try {
            const response = await api.put(`/admin/products/${productId}`, productData, token || undefined);
            return response;
        }
        catch (error) {
            console.error('Error updating product:', error);
            throw error;
        }
    },

    deleteProduct : async (productId: number) => {
        if (auth.getUserRole() !== 'admin') {
            throw new Error('Access denied');
        }
        const token = localStorage.getItem('token');
        try {
            const response = await api.delete(`/admin/products/${productId}`, token || undefined);
            return response;
        }
        catch (error) {
            console.error('Error deleting product:', error);
            throw error;
        }
    },

    // searchProducts : async (query: string, p0: { category: string; }) => {
    //     try {
    //         const response = await fetch(`${API_URL}/products/search?query=${encodeURIComponent(query)}&category=${encodeURIComponent(p0.category)}`);
    //         if (!response.ok) throw new Error('Failed to search products');
    //         products.products = await response.json();
    //     }   
    //     catch (error) {
    //         console.error('Error searching products:', error);
    //         throw error;
    //     }
    // },

    setCurrentImage: (url: string) => {
        products.currentImage = url
    },

    uploadGallery : async (productId: number, files: File[]) => {
      const token = localStorage.getItem('token')
      const formData = new FormData()
      files.forEach(file => formData.append('images', file))
        
      return api.upload(`/admin/products/${productId}/gallery`, formData, token || undefined)
    },

    uploadImage: async  (productId: number, file: File) => {
      const token = localStorage.getItem('token')
      const formData = new FormData()
      formData.append('image', file)
      
      return api.upload(`/admin/products/${productId}/image`, formData, token || undefined)
    },


    deleteMainImage: async (productId: number): Promise<any> => {
      const token = localStorage.getItem('token')
      return api.delete(`/admin/products/${productId}/image`, token || undefined)
    },

    deleteGalleryImage: async (productId: number, imageUrl: string): Promise<any>=>  {
      const token = localStorage.getItem('token')
      return api.delete(`/admin/products/${productId}/gallery`, token || undefined)
    },

    getProductImages: async (productId: number): Promise<any> => {
      return api.get(`/products/${productId}/images`)
    },
    
})