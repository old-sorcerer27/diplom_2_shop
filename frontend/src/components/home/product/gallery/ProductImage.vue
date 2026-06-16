<script setup lang="ts">
import { ref, computed } from 'vue'
import { images} from '@/utils/image'

const props = defineProps({
  src: {
    type: String,
    default: ''
  },
  alt: {
    type: String,
    default: 'Product image'
  },
  size: {
    type: String as () => 'original' | 'medium' | 'thumbnail',
    default: 'medium'
  },
  lazy: {
    type: Boolean,
    default: true
  },
  placeholderText: {
    type: String,
    default: 'Нет изображения'
  }
})

const isLoaded = ref(false)
const hasError = ref(false)

const sizeClass = computed(() => {
  const classes = {
    thumbnail: 'product-image-thumbnail',
    medium: 'product-image-medium',
    original: 'product-image-original'
  }
  return classes[props.size] || 'product-image-medium'
})

const imageSrc = computed(() => {
  if (!props.src || props.src === '') {
    return '/placeholder.jpg'
  }
  
  return images.getImageUrl(props.src, props.size)
})

const handleError = () => {
  hasError.value = true
  isLoaded.value = false
}

const handleLoad = () => {
  isLoaded.value = true
  hasError.value = false
}
</script>

<template>
  <div class="product-image-container" :class="sizeClass">
    <!-- Основное изображение -->
    <img
      v-if="src && !hasError"
      :src="imageSrc"
      :alt="alt"
      :loading="lazy ? 'lazy' : 'eager'"
      @error="handleError"
      @load="handleLoad"
      class="product-image"
      :class="{ loaded: isLoaded }"
    />
    

    <div v-if="!isLoaded && src && !hasError" class="image-skeleton">
      <div class="skeleton-spinner"></div>
    </div>
    
    <div v-if="!src || hasError" class="image-placeholder">
      <i class="fas fa-image"></i>
      <span>{{ placeholderText }}</span>
    </div>
  </div>
</template>


<style scoped>
.product-image-container {
  position: relative;
  overflow: hidden;
  background: #f8f9fa;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
}

.product-image {
  width: 100%;
  height: 100%;
  object-fit: contain;
  transition: opacity 0.3s, transform 0.3s;
  opacity: 0;
}

.product-image.loaded {
  opacity: 1;
}

.product-image-container:hover .product-image {
  transform: scale(1.02);
}

.product-image-thumbnail {
  width: 80px;
  height: 80px;
  border-radius: 8px;
}

.product-image-medium {
  width: 100%;
  height: 300px;
  border-radius: 12px;
}

.product-image-original {
  width: 100%;
  height: 400px;
  border-radius: 15px;
}

.image-skeleton {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f0f0;
}

.skeleton-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #e0e0e0;
  border-top: 3px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.image-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  color: #adb5bd;
  background: #f8f9fa;
  width: 100%;
  height: 100%;
}

.image-placeholder i {
  font-size: 2.5rem;
}

.image-placeholder span {
  font-size: 0.85rem;
}

/* Адаптивность */
@media (max-width: 768px) {
  .product-image-original {
    height: 300px;
  }
  
  .product-image-medium {
    height: 250px;
  }
}

@media (max-width: 576px) {
  .product-image-original {
    height: 250px;
  }
  
  .product-image-medium {
    height: 200px;
  }
}
</style>