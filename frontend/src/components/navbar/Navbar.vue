<script setup lang="ts">
import { HomePageStates, SetHomeState } from '@/utils/home';
import { auth } from '../../utils/auth'
import { cart } from '../../utils/cart'
</script>

<template>
    <nav class="navbar navbar-dark">
        <div>
        </div>
        <div class="container">
            <a href="/" class="navbar-brand">
                <i class="fas fa-store"></i> Мой Магазин
            </a>
            <div class="d-flex gap-2">
                <div v-if="!auth.isAuthenticated()">
                    <a href="login" class="btn btn-outline-light me-2">
                        <i class="fas fa-sign-in-alt"></i> Войти
                    </a>
                    <a href="register" class="btn btn-light">
                        <i class="fas fa-user-plus"></i> Регистрация
                    </a>
                </div>
                <div v-else>
                    <div class="dropdown d-inline-block">
                        <button class="btn btn-light dropdown-toggle" data-bs-toggle="dropdown">
                            <i class="fas fa-user-circle"></i> {{ auth.getCurrentUser()?.username || 'Профиль' }}
                        </button>
                        <ul class="dropdown-menu dropdown-menu-end">
                            <li><a class="dropdown-item" href="profile">
                                <i class="fas fa-id-card"></i> Личный кабинет
                            </a></li>
                            <!-- <li v-if="auth.user?.role === 'admin'"><a class="dropdown-item" href="AdminPanel.vue">
                                <i class="fas fa-user-shield"></i> Админ-панель
                            </a></li>
                            <li v-if="auth.user?.role === 'owner'"><a class="dropdown-item" href="OwnerPanel.vue">
                                <i class="fas fa-crown"></i> Панель владельца
                            </a></li> -->
                            <li><hr class="dropdown-divider"></li>
                            <li><a class="dropdown-item text-danger" href="#" @click.prevent="auth.logout()">
                                <i class="fas fa-sign-out-alt"></i> Выйти
                            </a></li>
                        </ul>
                    </div>
                </div>
                 <button class="btn btn-outline-light me-2" @click="SetHomeState(HomePageStates.CART) ">
                        Корзина <span class="badge bg-light text-dark ms-1">{{ cart.cart.length }}</span>
                </button>
            </div>
        </div>
    </nav>
</template>
