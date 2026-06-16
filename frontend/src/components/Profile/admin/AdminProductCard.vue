<script setup lang="ts">
import { computed } from 'vue'

interface Product {
  id: number
  name: string
  description: string
  price: number
  stock: number
  category: string
  image_url: string
  gallery: string[]
}

const props = defineProps<{
  product: Product
}>()

defineEmits<{
  (e: 'edit', id: number): void
  (e: 'delete', product: Product): void
}>()

const imageCount = computed(() => {
  let count = 0
  if (props.product.image_url) count++
  if (props.product.gallery) count += props.product.gallery.length
  return count
})

const formatPrice = (price: number) => {
  return price.toLocaleString('ru-RU')
}
</script>

<template>
  <div class="admin-product-card">
    <div class="product-image">
      <ProductImage
        :src="product.image_url || product.gallery?.[0]"
        :alt="product.name"
        size="medium"
        :lazy="true"
      />
    </div>
    
    <div class="product-info">
      <div class="product-header">
        <h3 class="product-name">{{ product.name }}</h3>
        <span class="product-id">ID: {{ product.id }}</span>
      </div>
      
      <div class="product-meta">
        <span class="product-category">{{ product.category || 'Без категории' }}</span>
        <span class="product-price">{{ formatPrice(product.price) }} ₽</span>
      </div>
      
      <div class="product-stats">
        <div class="stat">
          <i class="fas fa-box"></i>
          <span :class="['stock', product.stock > 0 ? 'in-stock' : 'out-of-stock']">
            {{ product.stock > 0 ? product.stock + ' шт.' : 'Нет в наличии' }}
          </span>
        </div>
        <div class="stat">
          <i class="fas fa-image"></i>
          <span>{{ imageCount }} изображений</span>
        </div>
      </div>
      
      <div class="product-actions">
        <button class="btn btn-edit" @click="$emit('edit', product.id)">
          <i class="fas fa-edit"></i> Редактировать
        </button>
        <button class="btn btn-delete" @click="$emit('delete', product)">
          <i class="fas fa-trash"></i>
        </button>
      </div>
    </div>
  </div>
</template>


<style scoped>
.admin-product-card {
  background: white;
  border-radius: 15px;
  overflow: hidden;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  transition: all 0.3s;
}

.admin-product-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
}

.product-image {
  height: 200px;
  background: #f8f9fa;
}

.product-info {
  padding: 1rem;
}

.product-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 0.5rem;
}

.product-name {
  margin: 0;
  font-size: 1rem;
  font-weight: 600;
  color: #333;
  flex: 1;
}

.product-id {
  font-size: 0.7rem;
  color: #999;
  background: #f0f0f0;
  padding: 2px 8px;
  border-radius: 20px;
  white-space: nowrap;
}

.product-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
}

.product-category {
  font-size: 0.8rem;
  color: #999;
  text-transform: uppercase;
}

.product-price {
  font-size: 1.1rem;
  font-weight: bold;
  color: #28a745;
}

.product-stats {
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
  padding-top: 0.75rem;
  border-top: 1px solid #f0f0f0;
}

.stat {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.85rem;
  color: #666;
}

.stat i {
  color: #667eea;
  width: 16px;
}

.stock.in-stock {
  color: #28a745;
}

.stock.out-of-stock {
  color: #dc3545;
}

.product-actions {
  display: flex;
  gap: 0.5rem;
}

.btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.85rem;
}

.btn-edit {
  flex: 1;
  background: #f0f0f0;
  color: #333;
}

.btn-edit:hover {
  background: #667eea;
  color: white;
}

.btn-delete {
  background: #dc3545;
  color: white;
  padding: 0.5rem;
}

.btn-delete:hover {
  background: #c82333;
}
</style>