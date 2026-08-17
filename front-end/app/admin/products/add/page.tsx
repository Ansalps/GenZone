'use client'

import ProductForm, { FormErrors } from "@/components/product-form"
import axios from "axios"
import { useRouter } from "next/navigation"
import { useState } from "react"
import { useCategories } from "@/hooks/useCategories"
import { ProductFormData } from "@/types/productFormData"

const validateForm = (data: ProductFormData): FormErrors => {
    const errors: FormErrors = {}

    // Category
    if (!data.categoryName.trim()) {
        errors.categoryName = "Please select a category"
    }

    // Product name
    if (!data.productName.trim()) {
        errors.productName = "Product name is required"
    } else if (data.productName.trim().length < 3) {
        errors.productName = "Product name must be at least 3 characters"
    }

    // Description
    if (!data.productDescription.trim()) {
        errors.productDescription = "Description is required"
    }

    // Price
    if (data.price <= 0) {
        errors.price = "Price must be greater than 0"
    }

    // Stock
    if (data.stock < 0 || !Number.isInteger(data.stock)) {
        errors.stock = "Stock must be a non-negative whole number"
    }

    // Size
    if (!data.size) {
        errors.size = "Please select a size"
    } else if (
        !["Small", "Medium", "Large"].includes(data.size)
    ) {
        errors.size = "Size must be Small, Medium, or Large"
    }

    // Discount
    if (
        data.discountPercentage < 0 ||
        data.discountPercentage > 100
    ) {
        errors.discountPercentage =
            "Discount percentage must be between 0 and 100"
    }

    return errors
}

export default function AddProduct() {
    const { categories } = useCategories()
    const router = useRouter()

    const [isLoading, setIsLoading] = useState(false)

    // Product form data
    const [formData, setFormData] = useState<ProductFormData>({
        categoryName: "",
        productName: "",
        productDescription: "",
        price: 0,
        stock: 0,
        size: "",
        popular: false,
        discountPercentage: 0,
    })

    // Actual image file
    const [imageFile, setImageFile] = useState<File | null>(null)

    // Validation errors
    const [errors, setErrors] = useState<FormErrors>({})

    // Handle text, number, select and checkbox inputs
    const onChange = (
        e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
    ) => {
        const { name, value, type } = e.target

        setFormData((prev) => ({
            ...prev,
            [name]:
                type === "checkbox"
                    ? (e.target as HTMLInputElement).checked
                    : type === "number"
                    ? Number(value)
                    : value,
        }))
    }

    // Handle image selection
    const onImageChange = (
        e: React.ChangeEvent<HTMLInputElement>
    ) => {
        const file = e.target.files?.[0] ?? null

        setImageFile(file)

        // Clear previous image error when user selects a file
        if (file) {
            setErrors((prev) => ({
                ...prev,
                productImage: undefined,
            }))
        }
    }

    const handleSubmit = async (
        e: React.FormEvent<HTMLFormElement>
    ) => {
        e.preventDefault()

        // -----------------------------
        // Validate normal form fields
        // -----------------------------

        const newErrors = validateForm(formData)

        // -----------------------------
        // Validate image
        // -----------------------------

        if (!imageFile) {
            newErrors.productImage =
                "Product image is required"
        } else if (!imageFile.type.startsWith("image/")) {
            newErrors.productImage =
                "Please select a valid image"
        } else if (imageFile.size > 5 * 1024 * 1024) {
            newErrors.productImage =
                "Image must be smaller than 5MB"
        }

        // -----------------------------
        // Stop if validation fails
        // -----------------------------

        if (Object.keys(newErrors).length > 0) {
            setErrors(newErrors)
            return
        }

        setErrors({})
        setIsLoading(true)

        try {
            // -----------------------------
            // Create multipart FormData
            // -----------------------------

            const data = new FormData()

            data.append(
                "category_name",
                formData.categoryName
            )

            data.append(
                "product_name",
                formData.productName
            )

            data.append(
                "product_description",
                formData.productDescription
            )

            data.append(
                "price",
                String(formData.price)
            )

            data.append(
                "stock",
                String(formData.stock)
            )

            data.append(
                "popular",
                String(formData.popular)
            )

            data.append(
                "size",
                formData.size
            )

            data.append(
                "discount_percentage",
                String(formData.discountPercentage)
            )

            // Add actual image file
            if (imageFile) {
                data.append("product_image", imageFile)
            }

            // -----------------------------
            // Send to backend
            // -----------------------------

            const response = await axios.post(
                "http://localhost:8080/admin/product",
                data,
                {
                    withCredentials: true,
                }
            )

            // -----------------------------
            // Success
            // -----------------------------

            if (
                response.status === 200 ||
                response.status === 201
            ) {
                router.push("/admin/products")
            }

        } catch (error) {
            console.error(
                "Failed to add product:",
                error
            )

            setErrors({
                productImage:
                    "Failed to add product. Please try again.",
            })
        } finally {
            setIsLoading(false)
        }
    }

    return (
        <div className="p-6 max-w-4xl mx-auto">

            <h1 className="font-bold text-2xl mb-6">
                Add Product
            </h1>

            <form
                onSubmit={handleSubmit}
                className="w-full"
            >
                <ProductForm
                    formData={formData}
                    onChange={onChange}
                    onImageChange={onImageChange}
                    isLoading={isLoading}
                    categories={categories}
                    errors={errors}
                />
            </form>

        </div>
    )
}