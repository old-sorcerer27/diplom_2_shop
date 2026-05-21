<script setup lang="ts">
import { API_URL } from '@/main';
import { auth, getToken } from '@/utils/auth';
import { editingData, NewReviewData, Review, reviews } from '@/utils/reviews';
import { products } from '@/utils/products';
import { computed, onMounted, ref } from 'vue';
import bootstrap, { Modal } from "bootstrap";



// import {onMounted, computed} from 'vue'
// const API_URL = 'http://localhost:8080/api/v1'
// import { products } from '../../utils/products'
// import {auth} from '../../utils/auth'
// import { cart } from '../../utils/cart'

let modalEle = ref(null);
let editModal: Modal | null = null;

onMounted(() => {
  editModal = new Modal(modalEle.value!);
});

const newReviewData = ref<NewReviewData>({
    rating: 5,
    text: '',
    productId: products.product?.ID || 0
});



let editReviewData = ref<editingData>({
    ID: 0,
    rating: 5,
    text: ''
});

    
const canSubmitReview = computed(() => {
    return newReviewData.value.text.trim().length >= 3 && newReviewData.value.rating > 0;
})



const formatDate = (date: Date) => {
    if (!date) return '';
    return new Date(date).toLocaleDateString('ru-RU', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });
}

const createReview=  async () => {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_URL}/comments/${newReviewData.value.productId}/like`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
}

const declensionReviews = (count: number) => {
    const lastDigit = count % 10;
    const lastTwo = count % 100;
    
    if (lastTwo >= 11 && lastTwo <= 19) return 'отзывов';
    if (lastDigit === 1) return 'отзыв';
    if (lastDigit >= 2 && lastDigit <= 4) return 'отзыва';
    return 'отзывов';
}

const submitReview = async () => {
    if (!canSubmitReview || !products.product?.ID) return;
    
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_URL}/products/${products.product?.ID}/comments`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                rating: newReviewData.value.rating,
                text: newReviewData .value.text
            })
        });
        
        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Ошибка добавления отзыва');
        }
        
        alert('Спасибо за отзыв!');
        newReviewData.value = { rating: 5, text: '', productId: products.product?.ID || 0 }; 
        await products.loadProduct(products.product.ID); // Обновляем страницу
        
    } catch (error) {
        console.error('Error submitting review:', error);
    }
}

const editReview = (Review: Review, editModal: Modal | null) => {
    if (!editModal) return;
    editReviewData.value.ID = Review.ID;
    editReviewData.value.rating = Review.rating;
    editReviewData.value.text = Review.text;
    editModal.show();
}

const updateReview = async () => {
    try {
        const response = await fetch(`${API_URL}/comments/${editReviewData.value.ID}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${getToken()}`
            },
            body: JSON.stringify({
                rating: editReviewData.value.rating,
                text: editReviewData.value.text
            })
        });
        
        if (!response.ok) throw new Error('Failed to update comment');
        
        alert('Отзыв обновлён');
        editModal?.hide();
        products.reloadProduct();
        
    } catch (error) {
        console.error('Error updating comment:', error);
        alert('Ошибка обновления отзыва');
    }
}

</script>  

<template>
    <div>
     <div class="container mt-4">
            <div id="comments" class="comments-section" v-if="products.product">
                <h3>
                    <i class="fas fa-comments"></i> Отзывы покупателей
                    <span v-if="reviews.Reviews">({{ reviews.Reviews.length }})</span>
                </h3>
                
                <div v-if="products.canLeaveReview" class="add-comment-box">
                    <h5><i class="fas fa-star"></i> Оставить отзыв</h5>
                    <div class="rating-input">
                        <i v-for="star in 5" :key="star"
                           class="fas fa-star rating-star"
                           :class="{ active: star <= newReviewData.rating }"
                           @click="newReviewData.rating = star">
                        </i>
                    </div>
                    <textarea class="form-control mb-3" 
                              v-model="newReviewData.text" 
                              rows="4"
                              placeholder="Поделитесь впечатлениями о товаре..."
                              maxlength="1000"></textarea>
                    <div class="d-flex justify-content-between align-items-center">
                        <small class="text-muted">{{ newReviewData.text.length }}/1000</small>
                        <button class="btn btn-primary" @click="submitReview()" :disabled="!canSubmitReview">
                            <i class="fas fa-paper-plane"></i> Отправить отзыв
                        </button>
                    </div>
                </div>
                
                <div v-else-if="auth.isAuthenticated() && !products.product.user_purchased" class="warning-message">
                    <i class="fas fa-info-circle"></i>
                    Вы можете оставить отзыв только после покупки этого товара.
                    <!-- <a href="#" @click.prevent="scrollToBuy">Перейти к покупке</a> -->
                </div>
                
                <div v-else-if="!auth.isAuthenticated()" class="warning-message">
                    <i class="fas fa-info-circle"></i>
                    <a href="login.html">Войдите</a> или <a href="register.html">зарегистрируйтесь</a>, чтобы оставить отзыв.
                </div>
                
                <div v-if="reviews.Reviews && reviews.Reviews.length > 0">
                    <div v-for="review in reviews.Reviews" :key="review.ID" class="comment-card">
                        <div class="comment-header">
                            <div class="comment-user">
                                <div class="user-avatar">
                                    <i class="fas fa-user"></i>
                                </div>
                                <div>
                                    <span class="user-name">{{ review.user.name }}</span>
                                </div>
                            </div>
                            <div class="comment-rating">
                                <i v-for="i in 5" :key="i"
                                   :class="i <= review.rating ? 'fas fa-star' : 'far fa-star'"
                                   style="color: #ffc107; font-size: 14px;"></i>
                            </div>
                        </div>
                        
                        <div class="comment-text">{{ review.text }}</div>
                        
                        <div class="comment-footer">
                            <div class="comment-date">
                                <i class="fas fa-calendar-alt"></i> {{review.created_At}}
                            </div>
                            <div class="comment-actions">
                                <button class="like-btn" :class="{ liked: review.user_liked }" 
                                        @click="reviews.toggleLike(review.ID)">
                                    <i class="fas fa-heart"></i> {{ review.likes || 0 }}
                                </button>
                                <button v-if="review.can_edit" class="edit-btn" @click="editReview(review, editModal)">
                                    <i class="fas fa-edit"></i>
                                </button>
                                <button v-if="review.can_delete" class="delete-btn" @click="reviews.deleteReview(review.ID)">
                                    <i class="fas fa-trash-alt"></i>
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
                
                <div v-else class="text-center p-5">
                    <i class="fas fa-comment-dots" style="font-size: 48px; color: #ccc;"></i>
                    <p class="mt-3">Пока нет отзывов. Будьте первым!</p>
                </div>
            </div>
        </div>

        <div class="modal fade" id="editCommentModal" tabindex="-1">
            <div class="modal-dialog">
                <div class="modal-content">
                    <div class="modal-header">
                        <h5 class="modal-title"><i class="fas fa-edit"></i> Редактировать отзыв</h5>
                        <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                    </div>
                    <div class="modal-body">
                        <div class="rating-input">
                            <i v-for="star in 5" :key="star"
                               class="fas fa-star rating-star"
                               :class="{ active: star <= editReviewData.rating }"
                               @click="editReviewData.rating = star">
                            </i>
                        </div>
                        <textarea class="form-control mt-3" 
                                  v-model="editReviewData.text" 
                                  rows="4"
                                  maxlength="1000"></textarea>
                    </div>
                    <div class="modal-footer">
                        <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Отмена</button>
                        <button type="button" class="btn btn-primary" @click="updateReview()">Сохранить</button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template> 


<style>


</style>
