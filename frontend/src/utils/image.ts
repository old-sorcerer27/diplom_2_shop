import { API_URL } from "@/main"

export const images = {
    getImageUrl(url: string, size: 'original' | 'medium' | 'thumbnail' = 'medium'): string {
    if (!url) return '/placeholder.jpg'

    if (url.startsWith('http') || url.startsWith('//')) {
      return url
    }
    
    const baseUrl = API_URL || 'http://localhost:8080'
     if (url.startsWith('/uploads')) {
      // Для разных размеров
      if (size === 'thumbnail') {
        return `${baseUrl}${url.replace('/original/', '/thumbnail/')}`
      }
      if (size === 'medium') {
        return `${baseUrl}${url.replace('/original/', '/medium/')}`
      }
      return `${baseUrl}${url}`
    }

    const cleanUrl = url.startsWith('/') ? url : `/${url}`
    const uploadPath = `/uploads/products/${size}${cleanUrl}`
    
    if (cleanUrl.includes('/original/')) {
      const newUrl = cleanUrl.replace('/original/', `/${size}/`)
      return `${baseUrl}${newUrl}`
    }
    
    return `${baseUrl}${uploadPath}`
  },


    isValidImageUrl(url: string): boolean {
      if (!url) return false
      return url.match(/\.(jpeg|jpg|gif|png|webp|svg)$/i) !== null
    },

    getMainImage(images: string[] | undefined): string {
      if (!images || images.length === 0) return ''
      return images[0] || ''
    },
}