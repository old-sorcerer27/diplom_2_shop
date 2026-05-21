<script setup lang="ts">
    import { ref } from 'vue';
    import { auth, RegisterData} from '../../utils/auth'
    import { navigateTo } from '../../utils/router';

    const form = ref<RegisterData>({
      name: '',
      email: '',
      password: ''
    });

    const errorMessage = ref('');
    const successMessage = ref('');
    const loading = ref(false);

    async function handleRegister() {
     errorMessage.value = '';
     loading.value = true;
     
     try {
        const result = await auth.register(form.value);
       
        if (result.success && result.user) {
            successMessage.value = result.message || 'Регистрация прошла успешно!';
            setTimeout(() => {
              navigateTo('/');
            }, 1000);
        } else {
          errorMessage.value = result.error || 'Ошибка регистрации';
        }
        } catch (error) {
          errorMessage.value = 'Произошла ошибка. Попробуйте позже.';
        } finally {
           console.log('Регистрация завершена');
          loading.value = false;
        }
    }

</script>

<template>
    <div class="auth-container">
        <div class="auth-card">
            <div class="auth-header">
                <h2>Создать аккаунт</h2>
                <p>Присоединяйтесь к нашему магазину</p>
            </div>
            <div class="auth-body">
                <div id="alertMessage"></div>
                
                <form id="registerForm" @submit.prevent="handleRegister">
                    <div class="form-group">
                        <i class="fas fa-user"></i>
                        <input type="text" 
                               class="form-control" 
                               id="name" 
                               v-model="form.name"
                               placeholder="Ваше имя" 
                               required>
                    </div>
                    
                    <div class="form-group">
                        <i class="fas fa-envelope"></i>
                        <input type="email" 
                               class="form-control" 
                               id="email" 
                               v-model="form.email"
                               placeholder="Email" 
                               required>
                    </div>
                    
                    <div class="form-group password-field">
                        <i class="fas fa-lock"></i>
                        <input type="password" 
                               class="form-control" 
                               id="password" 
                               v-model="form.password"
                               placeholder="Пароль" 
                               required>
                        <i class="fas fa-eye toggle-password" onclick="auth.togglePassword('password')"></i>
                    </div>
                    <div class="password-strength" id="passwordStrength"></div>
                    
                    <div class="form-group password-field">
                        <i class="fas fa-lock"></i>
                        <input type="password" 
                               class="form-control" 
                               id="confirmPassword" 
                               placeholder="Подтвердите пароль" 
                               required>
                        <i class="fas fa-eye toggle-password" onclick="auth.togglePassword('confirmPassword')"></i>
                    </div>
                    
                    <div class="terms form-check">
                        <input type="checkbox" class="form-check-input" id="terms" required>
                        <label class="form-check-label" for="terms">
                            Я согласен с <a href="#" data-bs-toggle="modal" data-bs-target="#termsModal">условиями использования</a>
                        </label>
                    </div>
                    
                    <button type="submit" class="btn-auth" id="registerBtn">
                        <span id="registerText">Зарегистрироваться</span>
                        <span id="registerLoading" style="display: none;">
                            <span class="loading"></span> Регистрация...
                        </span>
                    </button>
                </form>
                
                <div class="auth-links">
                    <a href="/login">Уже есть аккаунт? Войти</a>
                </div>
            </div>
        </div>
    </div>
    
    <div class="modal fade" id="termsModal" tabindex="-1">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title">Условия использования</h5>
                    <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                </div>
                <div class="modal-body">
                    <h6>1. Общие положения</h6>
                    <p>Регистрируясь на сайте, вы соглашаетесь с правилами магазина.</p>
                    <h6>2. Конфиденциальность</h6>
                    <p>Мы не передаем ваши данные третьим лицам.</p>
                    <h6>3. Ответственность</h6>
                    <p>Вы несете ответственность за сохранность вашего пароля.</p>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-primary" data-bs-dismiss="modal">Закрыть</button>
                </div>
            </div>
        </div>
    </div>
</template>