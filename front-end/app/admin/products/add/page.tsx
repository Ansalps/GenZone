'use client'

import ProductForm, { FormErrors } from "@/components/product-form"
import axios from "axios"
import { useRouter } from "next/navigation"
import { useState } from "react"
import { useCategories } from "@/hooks/useCategories"
import { ProductFormData } from "@/types/productFormData"

const validateForm = (data: ProductFormData): FormErrors => {
    const errors: FormErrors = {}

    if (!data.categoryName.trim()) {
        errors.categoryName = "Please select a category"
    }

    if (!data.productName.trim()) {
        errors.productName = "Product name is required"
    } else if (data.productName.trim().length < 3) {
        errors.productName = "Product name must be at least 3 characters"
    }

    if (!data.productDescription.trim()) {
        errors.productDescription = "Description is required"
    }

    if (data.price <= 0) {
        errors.price = "Price must be greater than 0"
    }

    const inventoryValues = Object.values(data.inventory ?? {})
    if (inventoryValues.length === 0 || inventoryValues.some((value) => value < 0 || !Number.isInteger(value))) {
        errors.inventory = "Inventory stock must be a non-negative whole number for each size"
    }

    if (data.discountPercentage < 0 || data.discountPercentage > 100) {
        errors.discountPercentage = "Discount percentage must be between 0 and 100"
    }

    return errors
}

export default function AddProduct() {
    const { categories } = useCategories()
    const router = useRouter()

    const [isLoading, setIsLoading] = useState(false)

    const [formData, setFormData] = useState<ProductFormData>({
        categoryName: "",
        productName: "",
        productDescription: "",
        price: 0,
        stock: 0,
        size: "",
        inventory: {
            Small: 0,
            Medium: 0,
            Large: 0,
        },
        popular: false,
        discountPercentage: 0,
    })

    const [imageFile, setImageFile] = useState<File | null>(null)
    const [errors, setErrors] = useState<FormErrors>({})

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

    const onInventoryChange = (size: string, value: number) => {
        setFormData((prev) => ({
            ...prev,
            inventory: {
                ...(prev.inventory ?? {}),
                [size]: value,
            },
        }))
    }

    const onImageChange = (
        e: React.ChangeEvent<HTMLInputElement>
    ) => {
        const file = e.target.files?.[0] ?? null

        setImageFile(file)

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

        const newErrors = validateForm(formData)

        if (!imageFile) {
            newErrors.productImage = "Product image is required"
        } else if (!imageFile.type.startsWith("image/")) {
            newErrors.productImage = "Please select a valid image"
        } else if (imageFile.size > 5 * 1024 * 1024) {
            newErrors.productImage = "Image must be smaller than 5MB"
        }

        if (Object.keys(newErrors).length > 0) {
            setErrors(newErrors)
            return
        }

        setErrors({})
        setIsLoading(true)

        try {
            const data = new FormData()
            const inventoryEntries = Object.entries(formData.inventory ?? {}).map(([size, stock]) => ({
                size,
                stock: Number(stock) || 0,
            }))

            const primarySizeEntry = inventoryEntries.find((item) => (item.stock ?? 0) > 0) ?? inventoryEntries[0]

            data.append("category_name", formData.categoryName)
            data.append("product_name", formData.productName)
            data.append("product_description", formData.productDescription)
            data.append("price", String(formData.price))
            data.append("popular", String(formData.popular))
            data.append("inventory", JSON.stringify(inventoryEntries))
            data.append("size", primarySizeEntry?.size ?? "Small")
            data.append("stock", String(primarySizeEntry?.stock ?? 0))
            data.append("discount_percentage", String(formData.discountPercentage))

            if (imageFile) {
                data.append("product_image", imageFile)
            }

            const response = await axios.post(
                "http://localhost:8080/admin/product",
                data,
                {
                    withCredentials: true,
                }
            )

            if (
                response.status === 200 ||
                response.status === 201
            ) {
                router.push("/admin/products")
            }

        } catch (error) {
            console.error("Failed to add product:", error)

            setErrors({
                productImage: "Failed to add product. Please try again.",
            })
        } finally {
            setIsLoading(false)
        }
    }

    return (
        <div className="min-h-screen bg-slate-950 px-4 py-8 text-white sm:px-6 lg:px-8">
            <div className="mx-auto max-w-7xl">
                <header className="mb-6 flex items-center justify-between gap-4 rounded-3xl border border-white/10 bg-white/5 p-5 shadow-2xl shadow-slate-950/30 backdrop-blur-xl">
                    <div>
                        <h1 className="text-2xl font-bold">Add Product</h1>
                        <p className="text-sm text-slate-300">Create a new product listing</p>
                    </div>
                </header>

                <main className="rounded-3xl border border-white/10 bg-slate-900/60 p-6 shadow-2xl shadow-slate-950/30 max-w-4xl mx-auto">
                    <form onSubmit={handleSubmit} className="w-full">
                        <ProductForm
                            formData={formData}
                            onChange={onChange}
                            onInventoryChange={onInventoryChange}
                            onImageChange={onImageChange}
                            isLoading={isLoading}
                            categories={categories}
                            errors={errors}
                        />
                    </form>
                </main>
            </div>
        </div>
    )
}