<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import ProductImage from './ProductImage.vue'


const props = defineProps({
  images: {
    type: Array as () => string[],
    default: () => []
  },
  productName: {
    type: String,
    default: 'Product'
  }
})

const currentIndex = ref(0)

const allImages = computed(() => {
  if (!props.images || props.images.length === 0) {
    console.log("изображение не найдено ", props.images.length)
    return ['/placeholder.jpg']
  }
  return props.images
})

const mainImage = computed(() => {
  return allImages.value[currentIndex.value] || ''
})

const selectImage = (index: number) => {
  currentIndex.value = index
}

// Сброс при изменении списка изображений
watch(() => props.images, () => {
  currentIndex.value = 0
})
</script>

<template>
  <div class="product-gallery">
    <div class="main-image-wrapper">
      <ProductImage
        :src="mainImage"
        :alt="productName"
        size="original"
        :lazy="false"
        :placeholder-text="'Нет изображения'"
      />
    </div>
    
    <div v-if="allImages.length > 1" class="thumbnail-list">
      <div
        v-for="(img, idx) in allImages"
        :key="idx"
        class="thumbnail-item"
        :class="{ active: currentIndex === idx }"
        @click="selectImage(idx)"
      >
        <ProductImage
          :src="img"
          :alt="`${productName} ${idx + 1}`"
          size="thumbnail"
          :lazy="true"
        />
      </div>
    </div>
    
    <div v-if="allImages.length > 1" class="image-counter">
      {{ currentIndex + 1 }} / {{ allImages.length }}
    </div>
  </div>
</template>


<style scoped>
.product-gallery {
  background: white;
  border-radius: 15px;
  padding: 1rem;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  position: relative;
}

.main-image-wrapper {
  width: 100%;
  height: 400px;
  background: #f8f9fa;
  border-radius: 12px;
  overflow: hidden;
  margin-bottom: 1rem;
}

.thumbnail-list {
  display: flex;
  gap: 0.5rem;
  overflow-x: auto;
  padding: 0.5rem 0;
  scrollbar-width: thin;
}

.thumbnail-list::-webkit-scrollbar {
  height: 4px;
}

.thumbnail-list::-webkit-scrollbar-track {
  background: #f0f0f0;
  border-radius: 4px;
}

.thumbnail-list::-webkit-scrollbar-thumb {
  background: #667eea;
  border-radius: 4px;
}

.thumbnail-item {
  flex-shrink: 0;
  cursor: pointer;
  border-radius: 8px;
  overflow: hidden;
  border: 2px solid transparent;
  transition: all 0.3s;
  width: 80px;
  height: 80px;
}

.thumbnail-item:hover {
  transform: translateY(-2px);
}

.thumbnail-item.active {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

.image-counter {
  position: absolute;
  bottom: 20px;
  right: 20px;
  background: rgba(0, 0, 0, 0.7);
  color: white;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: 500;
  backdrop-filter: blur(4px);
}

/* Адаптивность */
@media (max-width: 992px) {
  .main-image-wrapper {
    height: 350px;
  }
}

@media (max-width: 768px) {
  .main-image-wrapper {
    height: 300px;
  }
  
  .thumbnail-item {
    width: 60px;
    height: 60px;
  }
}

@media (max-width: 576px) {
  .main-image-wrapper {
    height: 250px;
  }
}
</style>