import {reactive } from 'vue'
import {cart} from './cart'
import { API_URL } from '@/main'
import { auth, getToken } from './auth';

export interface Review {
    ID: number;
    productId: number;
    user: {
        name: string;
        userId: number;
    };
    rating: number;
    text: string;
    created_At: string;
    updated_At: string;
    likes?: number;
    user_liked?: boolean;
    can_edit?: boolean;
    can_delete?: boolean;
}

export interface editingData {
    ID: number;
    rating: number;
    text: string;
}

export interface NewReviewData {
    rating: number;
    text: string;
    productId: number;
}

export interface EdditReviewResponse {
    success: boolean;
    message?: string;
    error?: string;
}


export const reviews = reactive({
    loading: false,
    Reviews: <Review[]>[],
    editRating: 0,
    editText: '',

    async loadReviews() {
        const token = localStorage.getItem('token');
        try {
            const response = await fetch(`${API_URL}/user/Reviews`, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
            
            if (response.ok) {
                this.Reviews = await response.json();
            }
        } catch (error) {
            console.error('Error loading Reviews:', error);
        }
    },

    async loadProductReviews(productId: number) {
        const token = localStorage.getItem('token');
        try {
            const response = await fetch(`${API_URL}/products/${productId}/Reviews`, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
            
            if (response.ok) {
                this.Reviews = await response.json();
            }
        } catch (error) {
            console.error('Error loading Reviews:', error);
        }
    },
    
    async deleteReview(ReviewId: number) {
        if (!confirm('Вы уверены, что хотите удалить этот отзыв?')) return;
        
        const token = localStorage.getItem('token');
        try {
            const response = await fetch(`${API_URL}/Reviews/${ReviewId}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
            
            if (response.ok) {
                alert('Отзыв удален');
                this.loadReviews();
            } else {
                alert('Ошибка удаления отзыва');
            }
        } catch (error) {
            console.error('Error deleting Review:', error);
            alert('Ошибка удаления отзыва');
        }
    },

    async toggleLike (reviewID: number) {
        if (!auth.isAuthenticated) {
            alert('Войдите, чтобы оценивать отзывы');
            window.location.href = 'login.html';
            return;
        }
        
        try {
            const response = await fetch(`${API_URL}/comments/${reviewID}/like`, {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${getToken()}`
                }
            });
            
            const result = await response.json();
            const comment = this.Reviews.find(c => c.ID === reviewID);
            if (comment) {
                comment.likes = result.likes;
            }
            
        } catch (error) {
            console.error('Error toggling like:', error);
        }
    }
    
})