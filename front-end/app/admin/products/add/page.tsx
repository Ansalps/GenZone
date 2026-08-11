'use client'

import ProductForm, { FormErrors } from "@/components/product-form"
import axios from "axios"
import { useRouter } from "next/navigation"
import { useState } from "react"
import { useCategories } from "@/hooks/useCategories"

// 1. Define explicit Form State interface to fix "formData is undefined"
export interface ProductFormData {
    categoryName: string;
    productName: string;
    productDescription: string;
    productImageUrl: string;
    price: number;
    stock: number;
    size: string;
    popular: boolean;
    discountPercentage: number;
}

// 2. Pass ProductFormData as type parameter
const validateForm = (data: ProductFormData): FormErrors => {
    const errors: FormErrors = {};

    if (!data.categoryName.trim()) {
        errors.categoryName = "Please select a category";
    }

    if (!data.productName.trim()) {
        errors.productName = "Product name is required";
    } else if (data.productName.trim().length < 3) {
        errors.productName = "Product name must be at least 3 characters";
    }

    if (!data.productDescription.trim()) {
        errors.productDescription = "Description is required";
    }

    if (!data.productImageUrl.trim()) {
        errors.productImageUrl = "Image URL is required";
    } else if (!/^https?:\/\/.+/i.test(data.productImageUrl.trim())) {
        errors.productImageUrl = "Must be a valid HTTP or HTTPS URL";
    }

    if (data.price <= 0) {
        errors.price = "Price must be greater than 0";
    }

    if (data.stock < 0 || !Number.isInteger(data.stock)) {
        errors.stock = "Stock must be a non-negative whole number";
    }

    if (!data.size) {
        errors.size = "Please select a size";
    } else if (!["Small", "Medium", "Large"].includes(data.size)) {
        errors.size = "Size must be Small, Medium, or Large";
    }

    if (data.discountPercentage < 0 || data.discountPercentage > 100) {
        errors.discountPercentage = "Discount percentage must be between 0 and 100";
    }

    return errors;
};

export default function AddProduct() {
    const { categories } = useCategories();
    const router = useRouter();
    const [isLoading, setIsLoading] = useState(false);

    // Cleaned state matching ProductFormData interface
    const [formData, setFormData] = useState<ProductFormData>({
        categoryName: '',
        productName: '',
        productDescription: '',
        productImageUrl: '',
        price: 0,
        stock: 0,
        size: '',
        popular: false,
        discountPercentage: 0,
    });

    const [errors, setErrors] = useState<FormErrors>({});

    const onChange = (
        e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
    ) => {
        const { name, value, type } = e.target;

        setFormData((prev) => ({
            ...prev,
            [name]:
                type === "checkbox"
                    ? (e.target as HTMLInputElement).checked
                    : type === "number"
                    ? Number(value)
                    : value,
        }));
    };

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        // Run validation
        const newErrors = validateForm(formData);

        // If errors exist, set state and STOP submission
        if (Object.keys(newErrors).length > 0) {
            setErrors(newErrors);
            return;
        }

        // Clear previous errors if valid
        setErrors({});
        setIsLoading(true);

        try {
            const response = await axios.post(
                'http://localhost:8080/admin/product',
                {
                    category_name: formData.categoryName,
                    product_name: formData.productName,
                    product_description: formData.productDescription,
                    product_image_url: formData.productImageUrl,
                    price: formData.price,
                    stock: formData.stock,
                    popular: formData.popular,
                    size: formData.size,
                    discount_percentage: formData.discountPercentage,
                },
                { withCredentials: true }
            );

            if (response.status === 200 || response.status === 201) {
                router.push("/admin/products");
            }
        } catch (error) {
            console.error("Failed to add product:", error);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="p-6 max-w-4xl mx-auto">
            <h1 className="font-bold text-2xl mb-6">Add Product</h1>
            <form onSubmit={handleSubmit} className="w-full">
                <ProductForm
                    formData={formData}
                    onChange={onChange}
                    isLoading={isLoading}
                    categories={categories}
                    errors={errors}
                />
            </form>
        </div>
    );
}