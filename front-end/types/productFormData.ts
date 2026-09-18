export interface ProductFormData {
    categoryName: string
    productName: string
    productDescription: string
    price: number
    stock: number
    size: string
    inventory?: Record<string, number>
    popular: boolean
    discountPercentage: number
}