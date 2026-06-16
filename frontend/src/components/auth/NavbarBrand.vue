<script setup lang="ts">
import { auth } from '../../utils/auth'
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
            </div>
        </div>
    </nav>
</template>

