<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { products } from '@/utils/products'
import { format } from '@/utils/utils'
import { navigateTo } from '@/utils/router'

const productId = computed(() => Number(route.params.id))
const isNew = computed(() => productId.value === 0)

const saving = ref(false)
const error = ref('')

interface ProductForm {
  name: string
  description: string
  price: number
  stock: number
  category: string
  image_url: string
  gallery: string[]
  created_at: string
}

const form = ref<ProductForm>({
  name: '',
  description: '',
  price: 0,
  stock: 0,
  category: '',
  image_url: '',
  gallery: [],
  created_at: '',
})

// Загрузка товара
const loadProduct = async () => {
  if (isNew.value) return
  
  try {
    await products.loadProduct(productId.value)
    if (products.product) {
      const p = products.product
      form.value = {
        name: p.name,
        description: p.description,
        price: p.price,
        stock: p.stock,
        category: p.category || '',
        image_url: p.image_url || '',
        gallery: p.gallery || [],
        created_at: p.created_at.toString(),
      }
    }
  } catch (err: any) {
    error.value = err.message || 'Ошибка загрузки товара'
  }
}

// Сохранение товара
const saveProduct = async () => {
  // Валидация
  if (!form.value.name.trim()) {
    error.value = 'Введите название товара'
    return
  }
  
  if (!form.value.description.trim()) {
    error.value = 'Введите описание товара'
    return
  }
  
  if (form.value.price <= 0) {
    error.value = 'Цена должна быть больше 0'
    return
  }
  
  if (form.value.stock < 0) {
    error.value = 'Количество не может быть отрицательным'
    return
  }
  
  saving.value = true
  error.value = ''
  
  try {
    const data = {
      name: form.value.name,
      description: form.value.description,
      price: form.value.price,
      stock: form.value.stock,
      category: form.value.category,
      image_url: form.value.image_url,
      gallery: form.value.gallery
    }
    
    let result
    if (isNew.value) {
      result = await products.createProduct(data)
    } else {
      result = await products.updateProduct(productId.value, data)
    }
    
    if (result) {
      navigateTo(`/admin/products/${result.id}`)  
    }
  } catch (err: any) {
    error.value = err.message || 'Ошибка сохранения'
  } finally {
    saving.value = false
  }
}

// Обновление изображений
const updateImages = (data: { mainImage: string; gallery: string[] }) => {
  form.value.image_url = data.mainImage
  form.value.gallery = data.gallery
}

// Сброс формы
const resetForm = () => {
  if (isNew.value) {
    form.value = {
      name: '',
      description: '',
      price: 0,
      stock: 0,
      category: '',
      image_url: '',
      gallery: [],
      created_at: '',
    }
  } else {
    loadProduct()
  }
}

const goBack = () => {
  navigateTo('/admin/products')
}

onMounted(() => {
  loadProduct()
})
</script>


<template>
  <div class="admin-product-edit">
    <div class="page-header">
      <div class="header-left">
        <button class="btn-back" @click="goBack">
          <i class="fas fa-arrow-left"></i> Назад
        </button>
        <h1>{{ isNew ? 'Создание товара' : 'Редактирование товара' }}</h1>
      </div>
      <div class="header-actions">
        <button class="btn btn-secondary" @click="resetForm">
          <i class="fas fa-undo"></i> Сбросить
        </button>
        <button class="btn btn-primary" @click="saveProduct" :disabled="saving">
          <i v-if="saving" class="fas fa-spinner fa-spin"></i>
          <i v-else class="fas fa-save"></i>
          {{ saving ? 'Сохранение...' : 'Сохранить' }}
        </button>
      </div>
    </div>
    
    <div v-if="error" class="alert alert-danger">
      <i class="fas fa-exclamation-circle"></i>
      {{ error }}
    </div>
    
    <div class="row">
      <!-- Основная информация -->
      <div class="col-lg-8">
        <div class="form-card">
          <h5><i class="fas fa-info-circle"></i> Основная информация</h5>
          
          <div class="form-group">
            <label class="required">Название товара</label>
            <input
              type="text"
              class="form-control"
              v-model="form.name"
              placeholder="Введите название товара"
              required
            >
          </div>
          
          <div class="form-group">
            <label>Категория</label>
            <select class="form-control" v-model="form.category">
              <option value="">Выберите категорию</option>
              <option value="Ноутбуки">Ноутбуки</option>
              <option value="Смартфоны">Смартфоны</option>
              <option value="Планшеты">Планшеты</option>
              <option value="Комплектующие">Комплектующие</option>
              <option value="Периферия">Периферия</option>
              <option value="Аудио">Аудио</option>
              <option value="Мониторы">Мониторы</option>
              <option value="Аксессуары">Аксессуары</option>
              <option value="Мебель">Мебель</option>
              <option value="Сетевое">Сетевое</option>
              <option value="Хранение">Хранение</option>
            </select>
          </div>
          
          <div class="form-group">
            <label class="required">Описание</label>
            <textarea
              class="form-control"
              v-model="form.description"
              rows="5"
              placeholder="Подробное описание товара"
            ></textarea>
          </div>
        </div>
      </div>
      
      <!-- Цена и наличие -->
      <div class="col-lg-4">
        <div class="form-card">
          <h5><i class="fas fa-tag"></i> Цена и наличие</h5>
          
          <div class="form-group">
            <label class="required">Цена (₽)</label>
            <input
              type="number"
              class="form-control"
              v-model.number="form.price"
              min="0"
              step="0.01"
              placeholder="0.00"
              required
            >
          </div>
          
          <div class="form-group">
            <label class="required">Количество на складе</label>
            <input
              type="number"
              class="form-control"
              v-model.number="form.stock"
              min="0"
              placeholder="0"
              required
            >
          </div>
        </div>
        
        <div class="form-card">
          <h5><i class="fas fa-calendar-alt"></i> Даты</h5>
          <div class="form-group">
            <label>Создан</label>
            <p class="form-control-static">{{ form.created_at || 'Не создан' }}</p>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Управление изображениями -->
    <div class="row mt-4">
      <div class="col-12">
        <AdminImageManager
          :product-id="productId"
          :initial-main-image="form.image_url"
          :initial-gallery="form.gallery || []"
          @update="updateImages"
        />
      </div>
    </div>
  </div>
</template>


<style scoped>
.admin-product-edit {
  max-width: 1400px;
  margin: 0 auto;
  padding: 1rem;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 1rem;
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 2px solid #f0f0f0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.header-left h1 {
  margin: 0;
  font-size: 1.8rem;
  font-weight: 700;
  color: #333;
}

.btn-back {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  background: none;
  border: none;
  color: #667eea;
  font-size: 1rem;
  cursor: pointer;
  padding: 0.5rem 1rem;
  border-radius: 8px;
  transition: all 0.3s;
}

.btn-back:hover {
  background: rgba(102, 126, 234, 0.1);
}

.header-actions {
  display: flex;
  gap: 0.5rem;
}

.btn {
  padding: 0.6rem 1.2rem;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
}

.btn-secondary {
  background: #f0f0f0;
  color: #333;
}

.btn-secondary:hover:not(:disabled) {
  background: #e0e0e0;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.form-card {
  background: white;
  border-radius: 15px;
  padding: 1.5rem;
  margin-bottom: 1.5rem;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.form-card h5 {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
  color: #333;
  font-weight: 600;
}

.form-group {
  margin-bottom: 1.25rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #555;
}

.form-group .required::after {
  content: ' *';
  color: #dc3545;
}

.form-control {
  width: 100%;
  padding: 0.75rem;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 1rem;
  transition: all 0.3s;
}

.form-control:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-control-static {
  padding: 0.75rem;
  background: #f8f9fa;
  border-radius: 8px;
  color: #666;
  margin: 0;
}

.alert {
  padding: 1rem;
  border-radius: 10px;
  margin-bottom: 1rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.alert-danger {
  background: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: stretch;
  }
  
  .header-left {
    flex-wrap: wrap;
  }
  
  .header-actions {
    flex-wrap: wrap;
  }
  
  .header-actions .btn {
    flex: 1;
    justify-content: center;
  }
}
</style>