import {reactive } from 'vue'
import {cart} from './cart'
import { API_URL } from '@/main'

const TOKEN_KEY = 'auth_token';
const USER_KEY = 'user_data';
const REMEMBER_ME_KEY = 'remember_me';

export interface User {
  id: number;
  email: string;
  name: string;
  role: 'client' | 'admin' | 'owner';
  isActive: boolean;
  createdAt: string;
  lastLoginAt?: string;
}

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterData {
  name: string;
  email: string;
  password: string;
}

export interface AuthResponse {
  success: boolean;
  token?: string;
  user?: User;
  message?: string;
  error?: string;
}

export interface ChangePasswordData {
  currentPassword: string;
  newPassword: string;
}

export let auth = reactive({
    loading: false,
    orders: [],
    userRole: <string | null>(null),


    isAuthenticated(): boolean {
      const token = getToken();
      if (!token) return false;
      
      const payload = decodeToken(token);
      if (!payload) return false;
      
       const now = Date.now() / 1000;
       return payload.exp > now;
    },

    isAdmin(): boolean {
      const role = this.getUserRole();
      return role === 'admin' || role === 'owner';
    },

    isOwner(): boolean {
      return this.getUserRole() === 'owner';
    },

    isValidEmail(email: string): boolean {
      const regex = /^[^\s@]+@([^\s@.,]+\.)+[^\s@.,]{2,}$/;
      return regex.test(email);
    },

    async register(data: RegisterData): Promise<AuthResponse> {
      console.log('Registering user with data:', data);
      try {
        // Валидация
        if (!data.name || !data.email || !data.password) {
          console.log('Validation failed: Missing fields');
          return {
            success: false,
            error: 'Пожалуйста, заполните все поля'
          };
        }

        if (data.name.length < 2) {
          console.log('Validation failed: Name too short');
          return {
            success: false,
            error: 'Имя должно содержать минимум 2 символа'
          };
        }

        if (!this.isValidEmail(data.email)) {
          console.log('Validation failed: Invalid email');
          return {
            success: false,
            error: 'Введите корректный email адрес'
          };
        }

        if (data.password.length < 6) {
          console.log('Validation failed: Password too short');
          return {
            success: false,
            error: 'Пароль должен содержать минимум 6 символов'
          };
        }

        const response = await fetch(`${API_URL}/auth/register`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            name: data.name.trim(),
            email: data.email.toLowerCase().trim(),
            password: data.password
          })
        });

        const result = await response.json();

        console.log('Register response:', result);

        if (!response.ok) {
          return {
            success: false,
            error: result.error || 'Ошибка регистрации'
          };
        }

        // Автоматически входим после регистрации
        if (result.token && result.user) {
            setToken(result.token);
            setUser(result.user);
        }

        return {
          success: true,
          token: result.token,
          user: result.user,
          message: 'Регистрация прошла успешно'
        };

      } catch (error) {
        console.error('Register error:', error);
        return {
          success: false,
          error: 'Ошибка сети. Попробуйте позже.'
        };
      }
    },

    async login(credentials: LoginCredentials): Promise<AuthResponse> {
        try {
          // Валидация на клиенте
          if (!credentials.email || !credentials.password) {
            return {
              success: false,
              error: 'Пожалуйста, заполните все поля'
            };
          }

          if (!this.isValidEmail(credentials.email)) {
            return {
              success: false,
              error: 'Введите корректный email адрес'
            };
          }

          if (credentials.password.length < 6) {
            return {
              success: false,
              error: 'Пароль должен содержать минимум 6 символов'
            };
          }

          const response = await fetch(`${API_URL}/auth/login`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              // email: credentials.email.toLowerCase().trim(),
              email: credentials.email,
              password: credentials.password
            })
          });

          const data = await response.json();

          if (!response.ok) {
            return {
              success: false,
              error: data.error || 'Ошибка входа. Проверьте email и пароль.'
            };
          }
          // Сохраняем данные
          if (data.token && data.user) {
            setToken(data.token);
            setUser(data.user);

            // Сохраняем email для "Запомнить меня"
            if (localStorage.getItem(REMEMBER_ME_KEY) === 'true') {
              localStorage.setItem('last_email', credentials.email);
            }
          }
                    console.log('isAuthenticated:', this.isAuthenticated);
          return {
            success: true,
            token: data.token,
            user: data.user,
            message: data.message || 'Вход выполнен успешно'
          };

        } catch (error) {
          console.error('Login error:', error);
          return {
            success: false,
            error: 'Ошибка сети. Проверьте подключение к интернету.'
          };
        }
    },
    

    logout(): void {
      clearToken();
      clearUser();
      localStorage.removeItem(REMEMBER_ME_KEY);
      localStorage.removeItem('last_email');
      cart.cart = [];
      window.location.reload();
      
      fetch(`${API_URL}/auth/logout`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${getToken()}`
        }
      }).catch(() => {});
    },

    async changePassword(data: ChangePasswordData): Promise<{ success: boolean; error?: string }> {
      const token = getToken();
      if (!token) {
        return { success: false, error: 'Не авторизован' };
      }
      
      if (data.newPassword.length < 6) {
        return { success: false, error: 'Новый пароль должен содержать минимум 6 символов' };
      }
      
      if (data.currentPassword === data.newPassword) {
        return { success: false, error: 'Новый пароль должен отличаться от текущего' };
      }
      
      try {
        const response = await fetch(`${API_URL}/auth/change-password`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          },
          body: JSON.stringify(data)
        });
        
        const result = await response.json();
        
        if (!response.ok) {
          return { success: false, error: result.error || 'Ошибка смены пароля' };
        }
        
        return { success: true };
        
      } catch (error) {
        return { success: false, error: 'Ошибка сети' };
      }
    },

    async loadOrders() {
        const token = localStorage.getItem('token');
        try {
            const response = await fetch(`${API_URL}/orders/user`, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
            
            if (response.ok) {
                this.orders = await response.json();
            }
        } catch (error) {
            console.error('Error loading orders:', error);
        }
    },

    getCurrentUser(): User | null {
        const userStr = localStorage.getItem(USER_KEY);
        if (!userStr) return null;
        
        try {
          return JSON.parse(userStr) as User;
        } catch {
          return null;
        }
    },

    getUserRole(): string {
        const user = this.getCurrentUser();
        return user?.role || 'client';
    },

    getRoleIcon(): string {
      const role = this.getUserRole();  
        switch (role) {
            case 'admin':
                return 'bi-shield-lock';
            case 'owner':
                return 'bi-cup-straw';
            default:
                return 'bi-person';
      }
    },

    togglePassword(fieldId: string) {
        const field = document.getElementById(fieldId);
        if (!field) return;
        const type = field.getAttribute('type') === 'password' ? 'text' : 'password';
        field.setAttribute('type', type);
    },

})


function setToken(token: string): void {
    localStorage.setItem(TOKEN_KEY, token);
}

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY);
}

function clearToken(): void {
  localStorage.removeItem(TOKEN_KEY);
}

function decodeToken(token: string): any {
  try {
    const base64Url = token.split('.')[1];
    if (!base64Url) return null;
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(c => {
      return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));
    return JSON.parse(jsonPayload);
  } catch {
    return null;
  }
}

function setUser(user: User): void {
  localStorage.setItem(USER_KEY, JSON.stringify(user));
}

function clearUser(): void {
  localStorage.removeItem(USER_KEY);
}

