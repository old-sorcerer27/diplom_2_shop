<script setup lang="ts">
import { ref, computed } from 'vue'
import { products } from '@/utils/products'
import { auth } from '@/utils/auth';


const props = defineProps<{
  productId: number
  initialMainImage?: string
  initialGallery?: string[]
}>()

const emit = defineEmits<{
  (e: 'update', data: { mainImage: string; gallery: string[] }): void
}>()

const mainImage = ref(props.initialMainImage || '')
const galleryImages = ref(props.initialGallery || [])
const uploading = ref(false)
const uploadProgress = ref(0)
const uploadError = ref('')
const isDragging = ref(false)

// Загрузка основного изображения
const uploadMainImage = async (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  
  // Валидация
  if (!validateFile(file)) return
  
  uploadError.value = ''
  uploading.value = true
  uploadProgress.value = 0
  
  try {
    const result = await products.uploadImage(props.productId, file)
    uploadProgress.value = 100
    
    if (result.success) {
      mainImage.value = result.image_url
      emit('update', {
        mainImage: mainImage.value,
        gallery: galleryImages.value
      })
    } else {
      throw new Error(result.error || 'Ошибка загрузки')
    }
  } catch (error: any) {
    uploadError.value = error.message || 'Ошибка загрузки изображения'
  } finally {
    uploading.value = false
    uploadProgress.value = 0
    input.value = ''
  }
}

// Загрузка галереи
const uploadGalleryImages = async (event: Event) => {
  const input = event.target as HTMLInputElement
  const files = input.files
  if (!files || files.length === 0) return
  
  // Проверяем лимит
  const totalImages = galleryImages.value.length + files.length
  if (totalImages > 5) {
    uploadError.value = `Максимум 5 изображений. Сейчас ${galleryImages.value.length}`
    return
  }
  
  uploadError.value = ''
  uploading.value = true
  uploadProgress.value = 0
  
  const formData = new FormData()
  for (const file of files) {
    if (validateFile(file)) {
      formData.append('images', file)
    }
  }
  
  try {
    // Имитация прогресса
    const progressInterval = setInterval(() => {
      uploadProgress.value = Math.min(uploadProgress.value + 10, 90)
    }, 200)
    
    const result = await products.uploadGallery(props.productId, Array.from(files))
    clearInterval(progressInterval)
    uploadProgress.value = 100
    
    if (result.success) {
      galleryImages.value = [...galleryImages.value, ...result.images]
      emit('update', {
        mainImage: mainImage.value,
        gallery: galleryImages.value
      })
    } else {
      throw new Error(result.error || 'Ошибка загрузки галереи')
    }
  } catch (error: any) {
    uploadError.value = error.message || 'Ошибка загрузки изображений'
  } finally {
    uploading.value = false
    uploadProgress.value = 0
    input.value = ''
    isDragging.value = false
  }
}

// Удаление основного изображения
const deleteMainImage = async () => {
  if (!confirm('Удалить основное изображение?')) return
  
  uploading.value = true
  try {
    const result = await products.deleteMainImage(props.productId)
    if (result.success) {
      mainImage.value = ''
      emit('update', {
        mainImage: '',
        gallery: galleryImages.value
      })
    } else {
      throw new Error(result.error || 'Ошибка удаления')
    }
  } catch (error: any) {
    uploadError.value = error.message || 'Ошибка удаления изображения'
  } finally {
    uploading.value = false
  }
}

// Удаление изображения из галереи
const removeGalleryImage = async (index: number) => {
  if (!confirm('Удалить это изображение?')) return
  
  const imageUrl = galleryImages.value[index]
  if (!imageUrl) return
  uploading.value = true
  
  try {
    const result = await products.deleteGalleryImage(props.productId, imageUrl)
    if (result.success) {
      galleryImages.value.splice(index, 1)
      emit('update', {
        mainImage: mainImage.value,
        gallery: galleryImages.value
      })
    } else {
      throw new Error(result.error || 'Ошибка удаления')
    }
  } catch (error: any) {
    uploadError.value = error.message || 'Ошибка удаления изображения'
  } finally {
    uploading.value = false
  }
}

// // Перемещение изображения в галерее
// const moveGalleryImage = (index: number, direction: number) => {
//   const newIndex = index + direction
//   if (newIndex < 0 || newIndex >= galleryImages.value.length) return
  
//   const temp = galleryImages.value[index]
//   galleryImages.value[index] = galleryImages.value[newIndex]
//   galleryImages.value[newIndex] = temp
  
//   emit('update', {
//     mainImage: mainImage.value,
//     gallery: galleryImages.value
//   })
// }

// Drag and Drop
const handleDrop = (event: DragEvent) => {
  isDragging.value = false
  const files = event.dataTransfer?.files
  if (files && files.length > 0) {
    const input = document.createElement('input')
    input.type = 'file'
    input.files = files
    uploadGalleryImages({ target: input } as any)
  }
}

// Валидация файла
const validateFile = (file: File): boolean => {
  const allowedTypes = ['image/jpeg', 'image/png', 'image/webp', 'image/gif']
  const maxSize = 5 * 1024 * 1024 // 5MB
  
  if (!allowedTypes.includes(file.type)) {
    uploadError.value = 'Поддерживаются только JPEG, PNG, WebP и GIF'
    return false
  }
  
  if (file.size > maxSize) {
    uploadError.value = 'Размер файла не должен превышать 5MB'
    return false
  }
  
  return true
}
</script>

<template>
    <Navbar/>
  <div v-if="auth.getUserRole() === 'admin' || auth.getUserRole() === 'owner'" class="admin-image-manager">
    <!-- Основное изображение -->
    <div class="image-section">
      <h5><i class="fas fa-image"></i> Основное изображение</h5>
      <div class="main-image-container">
        <div class="main-image-wrapper">
          <ProductImage
            v-if="mainImage"
            :src="mainImage"
            alt="Основное изображение"
            size="original"
            :lazy="false"
          />
          <div v-else class="no-image">
            <i class="fas fa-camera"></i>
            <span>Нет основного изображения</span>
          </div>
        </div>
        
        <div class="image-actions">
          <label class="btn-upload">
            <i class="fas fa-upload"></i> Загрузить
            <input
              type="file"
              accept="image/*"
              @change="uploadMainImage"
              :disabled="uploading"
            >
          </label>
          <button
            v-if="mainImage"
            class="btn-delete"
            @click="deleteMainImage"
            :disabled="uploading"
          >
            <i class="fas fa-trash"></i> Удалить
          </button>
        </div>
      </div>
    </div>

    <!-- Галерея изображений -->
    <div class="image-section">
      <div class="gallery-header">
        <h5><i class="fas fa-images"></i> Галерея изображений</h5>
        <span class="image-count">{{ galleryImages.length }} / 5</span>
      </div>
      
      <div class="gallery-container">
        <!-- Существующие изображения -->
        <div
          v-for="(img, index) in galleryImages"
          :key="index"
          class="gallery-item"
        >
          <ProductImage
            :src="img"
            :alt="`Изображение ${index + 1}`"
            size="medium"
            :lazy="false"
          />
          <div class="gallery-item-overlay">
            <!-- <button
              class="btn-move-left"
              v-if="index > 0"
              @click="moveGalleryImage(index, -1)"
              title="Переместить влево"
            >
              <i class="fas fa-chevron-left"></i>
            </button>
            <button
              class="btn-move-right"
              v-if="index < galleryImages.length - 1"
              @click="moveGalleryImage(index, 1)"
              title="Переместить вправо"
            >
              <i class="fas fa-chevron-right"></i>
            </button> -->
            <button
              class="btn-remove-gallery"
              @click="removeGalleryImage(index)"
              title="Удалить"
            >
              <i class="fas fa-times"></i>
            </button>
          </div>
        </div>
        
        <!-- Добавление нового изображения -->
        <div
          v-if="galleryImages.length < 5"
          class="gallery-item add-item"
          @dragover.prevent="isDragging = true"
          @dragleave.prevent="isDragging = false"
          @drop.prevent="handleDrop"
        >
          <label class="add-gallery-label" :class="{ dragging: isDragging }">
            <i class="fas fa-plus-circle"></i>
            <span>Добавить изображение</span>
            <input
              type="file"
              accept="image/*"
              multiple
              @change="uploadGalleryImages"
              :disabled="uploading"
            >
          </label>
        </div>
      </div>
      
      <div v-if="uploadProgress > 0" class="upload-progress">
        <div class="progress-bar" :style="{ width: uploadProgress + '%' }"></div>
        <span>{{ uploadProgress }}%</span>
      </div>
      
      <div v-if="uploadError" class="upload-error">
        <i class="fas fa-exclamation-circle"></i>
        {{ uploadError }}
      </div>
    </div>
  </div>
</template>


<style scoped>
.admin-image-manager {
  background: white;
  border-radius: 15px;
  padding: 1.5rem;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.image-section {
  margin-bottom: 2rem;
}

.image-section:last-child {
  margin-bottom: 0;
}

.image-section h5 {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 1rem;
  color: #333;
  font-weight: 600;
}

.gallery-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.image-count {
  font-size: 0.85rem;
  color: #999;
}

.main-image-container {
  display: flex;
  gap: 1.5rem;
  flex-wrap: wrap;
  align-items: flex-start;
}

.main-image-wrapper {
  width: 300px;
  height: 300px;
  border-radius: 12px;
  overflow: hidden;
  background: #f8f9fa;
  border: 2px dashed #e0e0e0;
}

.no-image {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #adb5bd;
  gap: 0.5rem;
}

.no-image i {
  font-size: 3rem;
}

.image-actions {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.btn-upload {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.6rem 1.2rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s;
}

.btn-upload:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
}

.btn-upload input {
  display: none;
}

.btn-delete {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.6rem 1.2rem;
  background: #dc3545;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s;
}

.btn-delete:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(220, 53, 69, 0.4);
}

button:disabled {
  opacity: 0.6;
  cursor: not-allowed !important;
  transform: none !important;
}

.gallery-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 1rem;
}

.gallery-item {
  position: relative;
  aspect-ratio: 1;
  border-radius: 10px;
  overflow: hidden;
  background: #f8f9fa;
  border: 2px solid #e0e0e0;
  transition: all 0.3s;
}

.gallery-item:hover {
  border-color: #667eea;
}

.gallery-item-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  opacity: 0;
  transition: opacity 0.3s;
}

.gallery-item:hover .gallery-item-overlay {
  opacity: 1;
}

.gallery-item-overlay button {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: none;
  background: white;
  color: #333;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.gallery-item-overlay button:hover {
  background: #667eea;
  color: white;
  transform: scale(1.1);
}

.btn-remove-gallery {
  background: #dc3545 !important;
  color: white !important;
}

.btn-remove-gallery:hover {
  background: #c82333 !important;
}

.add-item {
  border: 2px dashed #e0e0e0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.add-gallery-label {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  cursor: pointer;
  color: #999;
  transition: all 0.3s;
}

.add-gallery-label:hover {
  color: #667eea;
  background: rgba(102, 126, 234, 0.05);
}

.add-gallery-label.dragging {
  color: #667eea;
  background: rgba(102, 126, 234, 0.1);
}

.add-gallery-label i {
  font-size: 2rem;
}

.add-gallery-label input {
  display: none;
}

.upload-progress {
  margin-top: 1rem;
  background: #f0f0f0;
  border-radius: 20px;
  overflow: hidden;
  height: 24px;
  position: relative;
}

.progress-bar {
  height: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  transition: width 0.3s;
}

.upload-progress span {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 0.75rem;
  font-weight: 600;
  color: white;
}

.upload-error {
  margin-top: 0.5rem;
  padding: 0.5rem 1rem;
  background: #f8d7da;
  color: #721c24;
  border-radius: 8px;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.9rem;
}

@media (max-width: 768px) {
  .main-image-wrapper {
    width: 100%;
    height: 250px;
  }
  
  .image-actions {
    flex-direction: row;
    flex-wrap: wrap;
  }
  
  .gallery-container {
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  }
}
</style>