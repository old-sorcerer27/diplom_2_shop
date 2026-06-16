<script setup lang="ts">
    import { ref } from 'vue';
    import { auth, LoginCredentials } from '../../utils/auth'
    import { navigateTo } from '../../utils/router';
    import NavbarBrand from './NavbarBrand.vue';

    const form = ref<LoginCredentials>({
      email: '',
      password: ''
    });

    const errorMessage = ref('');
    const successMessage = ref('');
    const loading = ref(false);
    const showPassword = ref(false);
    const rememberMe = ref(false);

    async function handleLogin() {
     errorMessage.value = '';
     loading.value = true;
     
     try {
       const result = await auth.login(form.value);
       
       if (result.success && result.user) {
        console.log('Login successful:', result);
            successMessage.value = result.message || 'Вход выполнен успешно!';
            if (rememberMe.value) {
              localStorage.setItem('remember_me', 'true');
              console.log('Remember me enabled: session will persist');
            } else {
              localStorage.removeItem('remember_me');
                console.log('Remember me disabled: session will not persist');
            } 

            setTimeout(() => {
              navigateTo('/');
            }, 1000);

        } else {
         errorMessage.value = result.error || 'Ошибка входа';
        }
     } catch (error) {
       errorMessage.value = 'Произошла ошибка. Попробуйте позже.';
     } finally {
       loading.value = false;
       console.log('isAuthenticated:', auth.isAuthenticated);
     }
    }
    
</script>

<template>
      <NavbarBrand/>
    <div class="auth-container">
        <div class="auth-card">
            <div class="auth-header">
                <h2>Добро пожаловать!</h2>
                <p>Войдите в свой аккаунт</p>
            </div>
            <div class="auth-body">
                <div id="alertMessage"></div>
                
                <form id="loginForm"  @submit.prevent="handleLogin">
                    <div class="form-group">
                        <i class="fas fa-envelope"></i>
                        <input type="email" 
                               class="form-control" 
                               id="email" 
                               placeholder="Email" 
                               v-model="form.email"
                               required>
                    </div>
                    
                    <div class="form-group password-field">
                        <i class="fas fa-lock"></i>
                        <input type="password" 
                               class="form-control" 
                               id="password" 
                               placeholder="Пароль" 
                               v-model="form.password"
                               required>
                        <i class="fas fa-eye toggle-password" onclick="togglePassword('password')"></i>
                    </div>
                    
                    <div class="form-group">
                        <div class="form-check">
                            <input type="checkbox" class="form-check-input" id="remember" v-model="rememberMe" />
                            <label class="form-check-label" for="remember">Запомнить меня</label>
                        </div>
                    </div>
                    
                    <button type="submit" class="btn-auth" id="loginBtn" >
                        <i v-if="loading" class="fas fa-spinner fa-spin"></i>
                        <i v-else class="fas fa-sign-in-alt"></i>
                        {{ loading ? 'Вход...' : 'Войти' }}
                    </button>
                </form>
                
                <div class="auth-links">
                    <a href="#" onclick="showForgotPassword()">Забыли пароль?</a>
                    <span class="mx-2 text-muted">|</span>
                    <a href="register">Нет аккаунта? Зарегистрироваться</a>
                </div>
                
                <div class="social-login">
                    <p>Или войдите через</p>
                    <div class="social-buttons">
                        <div class="social-btn" onclick="socialLogin('google')">
                            <i class="fab fa-google" style="color: #DB4437;"></i>
                        </div>
                        <div class="social-btn" onclick="socialLogin('vk')">
                            <i class="fab fa-vk" style="color: #0077FF;"></i>
                        </div>
                        <div class="social-btn" onclick="socialLogin('github')">
                            <i class="fab fa-github"></i>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
    
    <div class="modal fade" id="forgotPasswordModal" tabindex="-1">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title">Восстановление пароля</h5>
                    <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
                </div>
                <div class="modal-body">
                    <p>Введите ваш email, и мы отправим инструкции по восстановлению пароля</p>
                    <input type="email" id="resetEmail" class="form-control" placeholder="Email">
                    <div id="resetMessage" class="mt-2"></div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Отмена</button>
                    <button type="button" class="btn btn-primary" onclick="resetPassword()">Отправить</button>
                </div>
            </div>
        </div>
    </div>
</template>

<style>
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

.auth-container {
    min-height: 100vh;
    min-width: 100vh;

    display: flex;
    align-items: center;
    justify-content: center;

    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.auth-card {
    background: white;
    border-radius: 20px;
    box-shadow: 0 20px 40px rgba(0,0,0,0.1);
    overflow: hidden;
    animation: fadeInUp 0.6s ease;
}

@keyframes fadeInUp {
    from {
        opacity: 0;
        transform: translateY(30px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

.auth-header {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    padding: 40px 30px;
    text-align: center;
}

.auth-header h2 {
    margin: 0;
    font-size: 28px;
    font-weight: 600;
}

.auth-header p {
    margin: 10px 0 0;
    opacity: 0.9;
}

.auth-body {
    padding: 40px 30px;
}

.form-group {
    margin-bottom: 25px;
    position: relative;
}

.form-group i {
    position: absolute;
    left: 15px;
    top: 50%;
    transform: translateY(-50%);
    color: #999;
}

.form-control {
    padding: 12px 15px 12px 45px;
    border: 2px solid #e0e0e0;
    border-radius: 10px;
    font-size: 16px;
    transition: all 0.3s;
}

.form-control:focus {
    border-color: #667eea;
    box-shadow: 0 0 0 0.2rem rgba(102, 126, 234, 0.25);
}

.btn-auth {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    padding: 12px;
    border: none;
    border-radius: 10px;
    font-size: 16px;
    font-weight: 600;
    width: 100%;
    cursor: pointer;
    transition: transform 0.2s;
}

.btn-auth:hover {
    transform: translateY(-2px);
}

.btn-auth:active {
    transform: translateY(0);
}

.auth-links {
    text-align: center;
    margin-top: 20px;
}

.auth-links a {
    color: #667eea;
    text-decoration: none;
    font-size: 14px;
}

.auth-links a:hover {
    text-decoration: underline;
}

.alert {
    border-radius: 10px;
    margin-bottom: 20px;
}

.password-strength {
    margin-top: 5px;
    font-size: 12px;
}

.strength-bar {
    height: 3px;
    margin-top: 5px;
    border-radius: 3px;
    transition: all 0.3s;
}

.terms {
    margin: 20px 0;
}

.loading {
    display: inline-block;
    width: 20px;
    height: 20px;
    border: 2px solid white;
    border-radius: 50%;
    border-top-color: transparent;
    animation: spin 0.6s linear infinite;
    margin-right: 10px;
}

@keyframes spin {
    to { transform: rotate(360deg); }
} 

</style>