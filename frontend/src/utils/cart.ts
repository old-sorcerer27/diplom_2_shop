import {reactive, computed } from 'vue'
import { auth, getToken } from './auth';
import { Product } from './products';
import { API_URL } from '@/main';



// type OrderRequest struct {
//     CustomerName  string `json:"customer_name"`
//     CustomerEmail string `json:"customer_email"`
//     Items         []struct {
//         ProductID int `json:"product_id"`
//         Quantity  int `json:"quantity"`
//     } `json:"items"`
// }

export interface OrderData {
  email: string;
  name: string;
  comment: string;
  phone: string;
}

export interface OrderItem {
  product_id: number;
  quantity: number;
//   product_name: string;
//   price: number;
}

export interface OrderPayload {
  items: OrderItem[]
  customer_name: string
  customer_email: string
  customer_phone: string
  comment: string
  total: number
  is_authenticated: boolean
}

export interface OrderResponse {
  success: boolean
  order?: {
    id: number
    total: number
    status: string
    created_at: string
  }
  error?: string
}

export const totalPrice = computed(() => {
  return cart.cart.reduce((sum: number, item: { product: Product; quantity: number }) => sum + (item.product?.price || 0) * item.quantity, 0)
})

export let cart = reactive({
    cart: <{ product: Product; quantity: number }[]>([]),
    loading: false,
    showModalOrder: false,


    checkoutCart() {
        if (this.cart.length === 0) {
            alert('Корзина пуста');
            return false;
        }
        return true;
    },

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

    // submitOrder: async () => {
    //     const userData = auth.getCurrentUser;
    //     if (!userData) {
    //         alert('Пожалуйста, войдите в систему, чтобы оформить заказ.');
    //         return;
    //     }
    //     try {
    //         let orderData = {
    //             customer_name: userData.name,
    //             customer_email: userData.email,
    //             items: this.cart.map(item => ({
    //                 product_id: item.product.id,
    //                 quantity: item.quantity
    //             }))
    //         };
    //     } catch (error) {
    //         console.error('Error creating order:', error);
    //         alert('Ошибка оформления заказа');
    //     }
    // }

    createOrder : async (payload: OrderPayload) : Promise<OrderResponse> => {
        try {
          const response = await fetch(`${API_URL}/orders`, {
            method: 'POST',
            headers: {'Authorization': `Bearer ${getToken()}`},
            body: JSON.stringify(payload)
          })
          console.log('Order creation response status:', JSON.stringify(payload), response.status, response.statusText);
          
          const data = await response.json()
          
          if (!response.ok) {
            throw new Error(data.error || 'Ошибка создания заказа')
          }
          
          return {
            success: true,
            order: {
              id: data.order_id,
              total: data.total,
              status: data.status,
              created_at: data.created_at
            }
          }
        } catch (error: any) {
          console.error('Order creation error:', error)
          return {
            success: false,
            error: error.message
          }
        }
    }

})

 

