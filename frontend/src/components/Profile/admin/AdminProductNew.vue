<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { products } from '@/utils/products'



const formRef = ref<InstanceType<typeof AdminProductForm>>()
const saving = ref(false)
const deleting = ref(false)
const showDeleteModal = ref(false)
const productName = ref('')

const productId = computed(() => {
  const id = route.params.id
  return id ? Number(id) : undefined
})

const isEditMode = computed(() => !!productId.value && productId.value > 0)

const initialData = ref<any>(null)
const notification = ref<{
  type: 'success' | 'error' | 'warning'
  icon: string
  message: string
} | null>(null)

// Загрузка товара для редактирования
const loadProduct = async () => {
  if (!isEditMode.value) return
  
  try {
    await products.loadProduct(productId.value!)
    if (products.product.value) {
      const p = products.product.value
      initialData.value = {
        name: p.name,
        description: p.description,
        price: p.price,
        stock: p.stock,
        category: p.category || '',
        sku: p.sku || '',
        image_url: p.image_url || '',
        gallery: p.gallery || [],
        attributes: p.attributes || [],
        meta_title: p.meta_title || '',
        meta_description: p.meta_description || ''
      }
      productName.value = p.name
    }
  } catch (error: any) {
    showNotification('error', 'Ошибка загрузки товара', error.message)
  }
}

// Сохранение товара
const saveProduct = async () => {
  if (!formRef.value) return
  
  const { valid, errors } = formRef.value.validate()
  
  if (!valid) {
    showNotification('error', 'Ошибка валидации', errors.join('. '))
    return
  }
  
  const formData = formRef.value.getFormData()
  
  saving.value = true
  
  try {
    let result
    if (isEditMode.value) {
      result = await products.updateProduct(productId.value!, formData)
    } else {
      result = await products.createProduct(formData)
    }
    
    if (result) {
      showNotification('success', 
        isEditMode.value ? 'Товар обновлен' : 'Товар создан',
        `Товар "${formData.name}" успешно ${isEditMode.value ? 'обновлен' : 'создан'}`
      )
      
      // Перенаправление через секунду
      setTimeout(() => {
        router.push('/admin/products')
      }, 1500)
    }
  } catch (error: any) {
    showNotification('error', 'Ошибка сохранения', error.message)
  } finally {
    saving.value = false
  }
}

// Удаление товара
const deleteProduct = () => {
  if (!isEditMode.value) return
  showDeleteModal.value = true
}

const confirmDelete = async () => {
  if (!productId.value) return
  
  deleting.value = true
  
  try {
    await products.deleteProduct(productId.value)
    showNotification('success', 'Товар удален', 'Товар успешно удален')
    
    setTimeout(() => {
      router.push('/admin/products')
    }, 1500)
  } catch (error: any) {
    showNotification('error', 'Ошибка удаления', error.message)
    showDeleteModal.value = false
  } finally {
    deleting.value = false
  }
}

// Сброс формы
const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetForm()
    if (!isEditMode.value) {
      showNotification('success', 'Форма сброшена', 'Все поля очищены')
    } else {
      loadProduct()
      showNotification('success', 'Изменения отменены', 'Данные восстановлены')
    }
  }
}

// Обновление формы
const handleFormUpdate = (data: any) => {
  // Можно добавить автосохранение черновика
}

// Уведомления
const showNotification = (type: 'success' | 'error' | 'warning', title: string, message: string) => {
  const icons = {
    success: 'fas fa-check-circle',
    error: 'fas fa-exclamation-circle',
    warning: 'fas fa-exclamation-triangle'
  }
  
  notification.value = {
    type,
    icon: icons[type],
    message: `${title}: ${message}`
  }
  
  setTimeout(() => {
    notification.value = null
  }, 5000)
}

const goBack = () => {
  router.push('/admin/products')
}

onMounted(() => {
  loadProduct()
})
</script>

<template>
  <div class="admin-product-create">
    <!-- Хлебные крошки -->
    <nav class="breadcrumb-nav">
      <ol class="breadcrumb">
        <li class="breadcrumb-item">
          <a href="/admin">Админ-панель</a>
        </li>
        <li class="breadcrumb-item">
          <a href="/admin/products">Товары</a>
        </li>
        <li class="breadcrumb-item active">
          {{ isEditMode ? 'Редактирование' : 'Создание' }} товара
        </li>
      </ol>
    </nav>

    <!-- Заголовок -->
    <div class="page-header">
      <div class="header-left">
        <button class="btn-back" @click="goBack">
          <i class="fas fa-arrow-left"></i> Назад
        </button>
        <h1>{{ isEditMode ? 'Редактирование товара' : 'Создание нового товара' }}</h1>
      </div>
      <div class="header-actions">
        <button class="btn btn-secondary" @click="resetForm">
          <i class="fas fa-undo"></i> Сбросить
        </button>
        <button class="btn btn-primary" @click="saveProduct" :disabled="saving">
          <i v-if="saving" class="fas fa-spinner fa-spin"></i>
          <i v-else class="fas fa-save"></i>
          {{ saving ? 'Сохранение...' : 'Сохранить товар' }}
        </button>
      </div>
    </div>

    <!-- Уведомления -->
    <div v-if="notification" class="notification" :class="notification.type">
      <i :class="notification.icon"></i>
      {{ notification.message }}
      <button class="btn-close" @click="notification = null">×</button>
    </div>

    <!-- Форма -->
    <AdminProductForm
      ref="formRef"
      :product-id="productId"
      :initial-data="initialData"
      @update="handleFormUpdate"
    />

    <!-- Дополнительные действия -->
    <div class="form-actions">
      <div class="actions-left">
        <button v-if="isEditMode" class="btn btn-danger" @click="deleteProduct">
          <i class="fas fa-trash"></i> Удалить товар
        </button>
      </div>
      <div class="actions-right">
        <button class="btn btn-secondary" @click="goBack">
          <i class="fas fa-times"></i> Отмена
        </button>
        <button class="btn btn-primary" @click="saveProduct" :disabled="saving">
          <i v-if="saving" class="fas fa-spinner fa-spin"></i>
          <i v-else class="fas fa-save"></i>
          {{ saving ? 'Сохранение...' : 'Сохранить товар' }}
        </button>
      </div>
    </div>

    <!-- Модальное окно подтверждения удаления -->
    <div v-if="showDeleteModal" class="modal-overlay" @click.self="showDeleteModal = false">
      <div class="modal-content">
        <div class="modal-header">
          <h5><i class="fas fa-exclamation-triangle"></i> Подтверждение удаления</h5>
          <button class="btn-close" @click="showDeleteModal = false">×</button>
        </div>
        <div class="modal-body">
          <p>Вы уверены, что хотите удалить товар <strong>"{{ productName }}"</strong>?</p>
          <p class="text-danger">Это действие нельзя отменить. Все данные о товаре будут удалены.</p>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showDeleteModal = false">
            Отмена
          </button>
          <button class="btn btn-danger" @click="confirmDelete" :disabled="deleting">
            <i v-if="deleting" class="fas fa-spinner fa-spin"></i>
            <i v-else class="fas fa-trash"></i>
            {{ deleting ? 'Удаление...' : 'Удалить' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>



<style scoped>
.admin-product-create {
  max-width: 1200px;
  margin: 0 auto;
  padding: 1rem;
}

.breadcrumb-nav {
  margin-bottom: 1.5rem;
}

.breadcrumb {
  display: flex;
  flex-wrap: wrap;
  padding: 0.5rem 1rem;
  margin: 0;
  background: transparent;
  list-style: none;
  gap: 0.5rem;
}

.breadcrumb-item {
  display: flex;
  align-items: center;
}

.breadcrumb-item + .breadcrumb-item::before {
  content: '/';
  margin-right: 0.5rem;
  color: #999;
}

.breadcrumb-item a {
  color: #667eea;
  text-decoration: none;
}

.breadcrumb-item.active {
  color: #999;
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
  font-size: 0.95rem;
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

.btn-danger {
  background: #dc3545;
  color: white;
}

.btn-danger:hover:not(:disabled) {
  background: #c82333;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none !important;
}

.form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 1rem;
  margin-top: 2rem;
  padding-top: 1.5rem;
  border-top: 2px solid #f0f0f0;
}

.actions-left,
.actions-right {
  display: flex;
  gap: 0.5rem;
}

.notification {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem 1.5rem;
  border-radius: 10px;
  margin-bottom: 1.5rem;
  position: relative;
  animation: slideDown 0.3s ease;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.notification.success {
  background: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.notification.error {
  background: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.notification.warning {
  background: #fff3cd;
  color: #856404;
  border: 1px solid #ffeeba;
}

.notification .btn-close {
  position: absolute;
  right: 1rem;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  opacity: 0.5;
  transition: opacity 0.3s;
}

.notification .btn-close:hover {
  opacity: 1;
}

/* Модальное окно */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  animation: fadeIn 0.3s;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modal-content {
  background: white;
  border-radius: 15px;
  max-width: 500px;
  width: 90%;
  overflow: hidden;
  animation: slideUp 0.3s;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(50px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.25rem 1.5rem;
  border-bottom: 1px solid #e0e0e0;
}

.modal-header h5 {
  margin: 0;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.modal-header .btn-close {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #999;
  transition: color 0.3s;
}

.modal-header .btn-close:hover {
  color: #333;
}

.modal-body {
  padding: 1.5rem;
}

.modal-body .text-danger {
  color: #dc3545;
  font-size: 0.9rem;
  margin-top: 0.5rem;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
  padding: 1rem 1.5rem;
  border-top: 1px solid #e0e0e0;
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
  
  .form-actions {
    flex-direction: column;
    align-items: stretch;
  }
  
  .actions-left,
  .actions-right {
    justify-content: center;
  }
}
</style>