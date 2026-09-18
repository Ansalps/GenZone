export interface ProductInventoryItem {
    size: string;
    stock: number;
}

export interface Product {
    id: number;
    created_at: string;

    category_id: number;
    category_name: string;

    product_name: string;
    product_description: string;

    product_image_url: string;

    price: number;
    stock: number;

    popular: boolean;
    size: string;
    inventory?: ProductInventoryItem[];

    discount_percentage?: number;
}