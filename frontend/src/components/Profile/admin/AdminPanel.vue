<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { products as productsManager } from '@/utils/products'
import { api } from '@/utils/api'
import { navigateTo } from '@/utils/router'
import { auth } from '@/utils/auth'
import Navbar from '@/components/auth/NavbarBrand.vue'

const products = ref<any[]>([])
const loading = ref(false)
const searchQuery = ref('')
const categoryFilter = ref('')
const categories = ref<string[]>([])

const loadProducts = async () => {
  loading.value = true
  try {
    await productsManager.loadProducts()
    products.value = productsManager.products
  } catch (error) {
    console.error('Failed to load products:', error)
  } finally {
    loading.value = false
  }
}

const loadCategories = async () => {
  try {
    const data = await api.get<string[]>('/products/categories')
    categories.value = data
  } catch (error) {
    console.error('Failed to load categories:', error)
  }
}

// const searchProducts = async () => {
//   if (searchQuery.value || categoryFilter.value) {
//     loading.value = true
//     try {
//       const result = await productsManager.searchProducts(searchQuery.value, {
//         category: categoryFilter.value
//       })
//       if (result) {
//         products.value = result.products
//       }
//     } catch (error) {
//       console.error('Failed to search products:', error)
//     } finally {
//       loading.value = false
//     }
//   } else {
//     await loadProducts()
//   }
// }


const createProduct = () => {
  navigateTo('/admin/products/new')
}

const editProduct = (id: number) => {
  navigateTo(`/admin/products/${id}`)
}

const deleteProduct = async (product: any) => {
  if (!confirm(`Удалить товар "${product.name}"?`)) return
  
  try {
    await productsManager.deleteProduct(product.id)
    await loadProducts()
  } catch (error) {
    console.error('Failed to delete product:', error)
    alert('Ошибка удаления товара')
  }
}

onMounted(() => {
  loadProducts()
  loadCategories()
})
</script>

<template>
     <Navbar/>
  <div v-if="auth.getUserRole() === 'admin' || auth.getUserRole() === 'owner'" class="admin-panel">
    <div class="page-header">
      <h1><i class="fas fa-store-alt"></i> Управление товарами</h1>
      <button class="btn btn-primary" @click="createProduct">
        <i class="fas fa-plus"></i> Добавить товар
      </button>
    </div>
    
    <!-- <div class="filters">
      <input
        type="text"
        v-model="searchQuery"
        placeholder="Поиск товаров..."
        class="search-input"
        @input="searchProducts"
      >
      <select v-model="categoryFilter" @change="searchProducts" class="filter-select">
        <option value="">Все категории</option>
        <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
      </select>
    </div> -->
    
    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <p>Загрузка...</p>
    </div>
    
    <div v-else-if="products.length === 0" class="empty-state">
      <i class="fas fa-box-open"></i>
      <h3>Нет товаров</h3>
      <p>Добавьте первый товар в магазин</p>
      <button class="btn btn-primary" @click="createProduct">
        <i class="fas fa-plus"></i> Добавить товар
      </button>
    </div>
    
    <div v-else class="products-grid">
      <AdminProductCard
        v-for="product in products"
        :key="product.id"
        :product="product"
        @edit="editProduct"
        @delete="deleteProduct"
      />
    </div>
  </div>
</template>


<style scoped>
.admin-panel {
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

.page-header h1 {
  margin: 0;
  font-size: 2rem;
  font-weight: 700;
  color: #333;
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

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
}

.filters {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
  flex-wrap: wrap;
}

.search-input {
  flex: 1;
  min-width: 200px;
  padding: 0.6rem 1rem;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 1rem;
  transition: all 0.3s;
}

.search-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.filter-select {
  padding: 0.6rem 1rem;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 1rem;
  background: white;
  min-width: 150px;
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1.5rem;
}

.loading {
  text-align: center;
  padding: 3rem;
}

.spinner {
  width: 50px;
  height: 50px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.empty-state {
  text-align: center;
  padding: 4rem;
  background: white;
  border-radius: 15px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.empty-state i {
  font-size: 4rem;
  color: #ddd;
  margin-bottom: 1rem;
}

.empty-state h3 {
  color: #333;
  margin-bottom: 0.5rem;
}

.empty-state p {
  color: #999;
  margin-bottom: 1.5rem;
}

@media (max-width: 768px) {
  .filters {
    flex-direction: column;
  }
  
  .search-input {
    min-width: 100%;
  }
  
  .products-grid {
    grid-template-columns: 1fr;
  }
}
</style>