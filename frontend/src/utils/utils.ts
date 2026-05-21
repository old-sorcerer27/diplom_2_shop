import {reactive } from 'vue'

export let format = reactive({
    formatPrice: (price: number) => {
        return price.toLocaleString('ru-RU');
    },

    formatDate: (date: Date) => {
        if (!date) return '';
        return new Date(date).toLocaleDateString('ru-RU', {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        });
    }
    
})