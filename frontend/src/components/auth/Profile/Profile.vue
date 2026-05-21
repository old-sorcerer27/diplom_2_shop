<script setup lang="ts">
import Navbar from '@/components/navbar/Navbar.vue';
import { auth, ChangePasswordData } from '@/utils/auth';
import { reviews } from '@/utils/reviews';
import { products } from '@/utils/products';
import { format } from '@/utils/utils';
import { onMounted, ref } from 'vue';


let activeTab = ref('orders');

const form = ref<ChangePasswordData>({
  currentPassword: '',
  newPassword: ''
});


onMounted(() => {
    reviews.loadReviews();    
    auth.loadOrders();
})


// mounted() {
//     auth.checkAuth();
//     auth.loadUserData();
//     auth.loadOrders();
//     auth.loadComments();
// }

        
    
//     formatPrice(price) {
//         return price.toLocaleString('ru-RU');
//     },
    
//     formatDate(date) {
//         return new Date(date).toLocaleDateString('ru-RU', {
//             year: 'numeric',
//             month: 'long',
//             day: 'numeric'
//         });
//     },
    
//     getRoleName(role) {
//         const roles = {
//             client: 'Клиент',
//             admin: 'Администратор',
//             owner: 'Владелец'
//         };
//         return roles[role] || role;
//     },
    
//     getRoleClass(role) {
//         return `role-${role}`;
//     },
    
//     getRoleIcon(role) {
//         const icons = {
//             client: 'fas fa-user',
//             admin: 'fas fa-user-shield',
//             owner: 'fas fa-crown'
//         };
//         return icons[role] || 'fas fa-user';
//     },
    
//     getStatusName(status) {
//         const statuses = {
//             pending: 'Ожидает обработки',
//             paid: 'Оплачен',
//             shipped: 'Отправлен',
//             delivered: 'Доставлен',
//             cancelled: 'Отменен'
//         };
//         return statuses[status] || status;
//     }
// }
</script>

<template >
    <Navbar/>    
        <div v-if="auth.isAuthenticated()" class="container mt-4">
            <div v-if="auth.loading" class="loading">
                <div class="spinner"></div>
                <p class="mt-3">Загрузка...</p>
            </div>
            
            <div v-else>
                <!-- Профиль пользователя -->
                <div class="profile-header">
                    <div class="row align-items-center">
                        <div class="col-md-2 text-center">
                            <div class="avatar">
                                <i class="fas fa-user"></i>
                            </div>
                        </div>
                        <div class="col-md-6">
                            <h2>{{ auth.getCurrentUser()?.name }}</h2>
                            <p class="text-muted">
                                <i class="fas fa-envelope"></i> {{ auth.getCurrentUser()?.email }}
                            </p>
                            <span :class="['role-badge', auth.getCurrentUser()?.role]">
                                <i :class="auth.getRoleIcon()"></i> 
                                {{ auth.getUserRole() }}
                            </span>
                            <p class="text-muted mt-2">
                                <i class="fas fa-calendar-alt"></i> 
                                Зарегистрирован: {{ auth.getCurrentUser()?.createdAt}}
                            </p>
                        </div>
                        <div class="col-md-4">
                            <div class="row text-center">
                                <div class="col-6">
                                    <div class="stat-card" @click="activeTab = 'orders'">
                                        <div class="stat-icon">
                                            <i class="fas fa-shopping-bag"></i>
                                        </div>
                                        <div class="stat-value">{{ auth.orders.length }}</div>
                                        <div class="text-muted">Заказов</div>
                                    </div>
                                </div>
                                <div class="col-6">
                                    <div class="stat-card" @click="activeTab = 'comments'">
                                        <div class="stat-icon">
                                            <i class="fas fa-comments"></i>
                                        </div>
                                        <div class="stat-value">{{ reviews.Reviews.length }}</div>
                                        <div class="text-muted">Отзывов</div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                
 
                <ul class="nav nav-tabs">
                    <li class="nav-item">
                        <a class="nav-link" :class="{ active: activeTab === 'orders' }" 
                           @click="activeTab = 'orders'" href="#">
                            <i class="fas fa-shopping-bag"></i> Мои заказы
                        </a>
                    </li>
                    <li class="nav-item">
                        <a class="nav-link" :class="{ active: activeTab === 'comments' }" 
                           @click="activeTab = 'comments'" href="#">
                            <i class="fas fa-star"></i> Мои отзывы
                        </a>
                    </li>
                    <li class="nav-item">
                        <a class="nav-link" :class="{ active: activeTab === 'settings' }" 
                           @click="activeTab = 'settings'" href="#">
                            <i class="fas fa-cog"></i> Настройки
                        </a>
                    </li>
                </ul>
                
           
                <!-- <div v-if="activeTab === 'orders'">
                    <div v-if="auth.orders.length === 0" class="alert alert-info">
                        <i class="fas fa-info-circle"></i> У вас пока нет заказов
                    </div>
                    <div v-else>
                        <div v-for="order in auth.orders" :key="order.id" class="order-card">
                            <div class="d-flex justify-content-between align-items-start">
                                <div>
                                    <h5>Заказ #{{ order.id }}</h5>
                                    <p class="text-muted">{{ formatDate(order.created_at) }}</p>
                                </div>
                                <div>
                                    <span :class="['order-status', 'status-' + order.status]">
                                        {{ getStatusName(order.status) }}
                                    </span>
                                </div>
                            </div>
                            <hr>
                            <div v-for="item in order.items" :key="item.id" class="mb-2">
                                <div class="d-flex justify-content-between">
                                    <div>
                                        <strong>{{ item.product.name }}</strong> x {{ item.quantity }}
                                    </div>
                                    <div>{{ formatPrice(item.price * item.quantity) }} ₽</div>
                                </div>
                            </div>
                            <hr>
                            <div class="d-flex justify-content-between align-items-center">
                                <div>
                                    <strong>Итого: {{ formatPrice(order.total) }} ₽</strong>
                                </div>
                                <div v-if="order.status === 'delivered'">
                                    <button class="btn btn-sm btn-outline-success" 
                                            @click="goToProduct(order.items[0].product_id)">
                                        <i class="fas fa-star"></i> Оставить отзыв
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                </div> -->
                
              
                <!-- <div v-if="activeTab === 'comments'">
                    <div v-if="comments.comments.length === 0" class="alert alert-info">
                        <i class="fas fa-info-circle"></i> Вы еще не оставляли отзывы
                    </div>
                    <div v-else>
                        <div v-for="Comment in comments.comments" :key="Comment.id" class="comment-card">
                            <div class="d-flex justify-content-between">
                                <div>
                                    <a :href="'/product-detail.html?id=' + Comment.productId" 
                                       style="text-decoration: none; color: #333;">
                                        <h6><i class="fas fa-box"></i> {{ Comment.product.name }}</h6>
                                    </a>
                                    <div class="rating mb-2">
                                        <i v-for="i in 5" :key="i" 
                                           :class="i <= Comment.rating ? 'fas fa-star' : 'far fa-star'"
                                           style="color: #ffc107;"></i>
                                    </div>
                                    <p>{{ Comment.text }}</p>
                                    <small class="text-muted">{{ formatDate(Comment.created_at) }}</small>
                                </div>
                                <div>
                                    <button class="btn btn-sm btn-outline-primary" 
                                            @click="comments.editComment(Comment)">
                                        <i class="fas fa-edit"></i>
                                    </button>
                                    <button class="btn btn-sm btn-outline-danger ms-2" 
                                            @click="comments.deleteComment(Comment.id)">
                                        <i class="fas fa-trash"></i>
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                </div> -->
                
                <div v-if="activeTab === 'settings'">
                    <div class="card">
                        <div class="card-body">
                            <h5><i class="fas fa-key"></i> Смена пароля</h5>
                            <form @submit.prevent="auth.changePassword(form)">
                                <div class="mb-3">
                                    <label>Текущий пароль</label>
                                    <input 
                                        type="password" 
                                        class="form-control" 
                                        id="currentPassword"
                                        v-model="form.currentPassword" 
                                        required
                                    >
                                </div>
                                <div class="mb-3">
                                    <label>Новый пароль</label>
                                    <input 
                                        type="password" 
                                        class="form-control" 
                                        id="newPassword"
                                        v-model="form.newPassword" 
                                        required 
                                        minlength="6"
                                    >
                                    <i class="fas fa-eye toggle-password" onclick="togglePassword('newPassword')"></i>
                                </div>
                                <div class="mb-3">
                                    <label>Подтверждение пароля</label>
                                    <input 
                                        type="password" 
                                        class="form-control" 
                                        id="confirmPassword"
                                        required
                                    >
                                    <i class="fas fa-eye toggle-password" onclick="togglePassword('confirmPassword')"></i>
                                </div>
                                <button type="submit" class="btn btn-primary">Изменить пароль</button>
                            </form>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div v-else class="container mt-4">
            <div class="alert alert-warning text-center">
                <i class="fas fa-exclamation-triangle"></i> Пожалуйста, войдите в свой аккаунт, чтобы просмотреть профиль
            </div>
        </div>
</template>

<style>

</style>

