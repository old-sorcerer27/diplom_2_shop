<!-- <script setup lang="ts">
import { ref, watch, debounce } from 'vue'
// import products from '../../utils/products'

const query = defineModel<string>('query')
const emit = defineEmits<{ (e: 'search'): void }>()

const suggestions = ref<string[]>([])
const showSuggestions = ref(false)
const popularSearches = ref(['iPhone', 'Ноутбук', 'Наушники'])

const onInput = debounce(async () => {
  if (query.value && query.value.length >= 2) {
    suggestions.value = await productsService.getSuggestions(query.value)
    showSuggestions.value = true
  } else {
    suggestions.value = []
    showSuggestions.value = false
  }
  emit('search')
}, 300)

function clearSearch() {
  query.value = ''
  emit('search')
}

function applySuggestion(suggestion: string) {
  query.value = suggestion
  showSuggestions.value = false
  emit('search')
}
</script>

<template>
  <div class="search-section">
    <div class="search-box">
      <i class="fas fa-search search-icon"></i>
      <input type="text" 
             class="search-input" 
             v-model="query"
             @input="onInput"
             @keyup.enter="search"
             placeholder="Поиск товаров...">
      <button class="search-clear" v-if="query" @click="clearSearch">
        <i class="fas fa-times-circle"></i>
      </button>
      
      <div v-if="showSuggestions && suggestions.length" class="suggestions">
        <div v-for="suggestion in suggestions" 
             :key="suggestion"
             @click="applySuggestion(suggestion)"
             class="suggestion-item">
          <i class="fas fa-search me-2"></i> {{ suggestion }}
        </div>
      </div>
    </div>
    
    <div v-if="!query && popularSearches.length" class="popular-searches">
      <span><i class="fas fa-fire"></i> Популярное:</span>
      <span v-for="term in popularSearches" 
            :key="term"
            @click="applySuggestion(term)"
            class="popular-term">
        {{ term }}
      </span>
    </div>
  </div>
</template> -->
