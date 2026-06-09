<!-- <script setup lang="ts">
import { computed } from 'vue'

const props = defineProps({
    size: {
        type: Number,
        default: 80
    },
    gradientStart: {
        type: String,
        default: '#667eea'
    },
    gradientEnd: {
        type: String,
        default: '#764ba2'
    },
    useGradient: {
        type: Boolean,
        default: true
    },
    name: {
        type: String,
        default: ''
    },
    imageUrl: {
        type: String,
        default: ''
    }
})

const gradientId = `avatar-gradient-${Math.random().toString(36).substr(2, 9)}`

const initials = computed(() => {
    if (!props.name) return ''
    const names = props.name.split(' ')
    if (names.length >= 2) {
        return (names[0][0] + names[1][0]).toUpperCase()
    }
    return names[0][0].toUpperCase()
})

const sizeClass = computed(() => {
    if (props.size <= 40) return 'avatar-sm'
    if (props.size <= 60) return 'avatar-md'
    return 'avatar-lg'
})
</script>

<template>
    <div class="avatar" :class="sizeClass">
        <svg :width="size" :height="size" viewBox="0 0 80 80" fill="none" xmlns="http://www.w3.org/2000/svg">

            <circle cx="40" cy="40" r="38" :fill="gradientStart" stroke="#fff" stroke-width="2"/>
            

            <defs v-if="useGradient">
                <linearGradient :id="gradientId" x1="0%" y1="0%" x2="100%" y2="100%">
                    <stop offset="0%" :stop-color="gradientStart"/>
                    <stop offset="100%" :stop-color="gradientEnd"/>
                </linearGradient>
            </defs>
            
    
            <circle cx="40" cy="30" r="10" fill="#fff" fill-opacity="0.8"/>
            
         
            <path d="M20 60C20 52 24 44 40 44C56 44 60 52 60 60" stroke="#fff" stroke-width="3" fill="none" stroke-linecap="round"/>
            
        
            <text v-if="initials" x="40" y="50" text-anchor="middle" fill="#fff" font-size="24" font-weight="bold" font-family="Arial">
                {{ initials }}
            </text>
        </svg>
    </div>
</template>


<style scoped>
.avatar {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    overflow: hidden;
}

.avatar svg {
    width: 100%;
    height: 100%;
    display: block;
}

.avatar-sm {
    width: 40px;
    height: 40px;
}

.avatar-md {
    width: 60px;
    height: 60px;
}

.avatar-lg {
    width: 80px;
    height: 80px;
}

.avatar-xl {
    width: 120px;
    height: 120px;
}
</style> -->