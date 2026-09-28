'use client'

import axios from "axios";
import Link from "next/link";
import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import ProductForm, { FormErrors } from "@/components/product-form";
import { useCategories } from "@/hooks/useCategories";
import { ProductFormData } from "@/types/productFormData";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

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

    if (data.price <= 0) {
        errors.price = "Price must be greater than 0";
    }

    const inventoryValues = Object.values(data.inventory ?? {});
    if (inventoryValues.length === 0 || inventoryValues.some((value) => value < 0 || !Number.isInteger(value))) {
        errors.inventory = "Inventory stock must be a non-negative whole number for each size";
    }

    if (data.discountPercentage < 0 || data.discountPercentage > 100) {
        errors.discountPercentage = "Discount percentage must be between 0 and 100";
    }

        // Date validation: if either date provided, both required; end >= start
        if ((data.startDate && !data.endDate) || (!data.startDate && data.endDate)) {
            if (!data.startDate) errors.startDate = "Start date is required when end date is provided"
            if (!data.endDate) errors.endDate = "End date is required when start date is provided"
        }

        if (data.startDate && data.endDate) {
            const s = new Date(data.startDate)
            const e = new Date(data.endDate)
            if (isNaN(s.getTime())) errors.startDate = "Invalid start date"
            if (isNaN(e.getTime())) errors.endDate = "Invalid end date"
            if (!errors.startDate && !errors.endDate && e < s) {
                errors.endDate = "End date must be the same or after start date"
            }
        }

    if (data.discountPercentage > 0) {
        if (!data.startDate) {
            errors.startDate = "Start date is required when discount is provided";
        }
        if (!data.endDate) {
            errors.endDate = "End date is required when discount is provided";
        }
        if (data.startDate && data.endDate && new Date(data.startDate) >= new Date(data.endDate)) {
            errors.endDate = "End date must be after start date";
        }
    }

    return errors;
};

export default function EditProduct() {
    const { categories } = useCategories();
    const router = useRouter();
    const params = useParams();
    const id = params?.id as string;

    const [isLoading, setIsLoading] = useState(false);
    const [isFetching, setIsFetching] = useState(true);
    const [submitError, setSubmitError] = useState<string>("");

    const [formData, setFormData] = useState<ProductFormData>({
        categoryName: "",
        productName: "",
        productDescription: "",
        price: 0,
        stock: 0,
        size: "Small",
        inventory: {
            Small: 0,
            Medium: 0,
            Large: 0,
        },
        popular: false,
        discountPercentage: 0,
        startDate: "",
        endDate: "",
    });

    const [existingImageUrl, setExistingImageUrl] = useState<string>("");
    const [imageFile, setImageFile] = useState<File | null>(null);
    const [errors, setErrors] = useState<FormErrors>({});

    useEffect(() => {
        if (!id) return;

        const fetchProduct = async () => {
            try {
                const response = await axios.get(`${API_BASE_URL}/admin/product/${id}`, {
                    withCredentials: true,
                });

                if (response.data?.data) {
                    const product = response.data.data;
                    const inventoryMap: Record<string, number> = {
                        Small: 0,
                        Medium: 0,
                        Large: 0,
                    };

                    (product.variants || []).forEach((variant: { size: string; stock: number }) => {
                        if (variant.size && inventoryMap[variant.size] !== undefined) {
                            inventoryMap[variant.size] = variant.stock ?? 0;
                        }
                    });

                    // Format dates for HTML date/time inputs if existing
                    const formattedStartDate = product.offer?.start_at 
                        ? new Date(product.offer.start_at).toISOString().slice(0, 16) 
                        : product.start_date || "";
                    const formattedEndDate = product.offer?.end_at 
                        ? new Date(product.offer.end_at).toISOString().slice(0, 16) 
                        : product.end_date || "";

                    setFormData({
                        categoryName: product.category?.category_name || product.category_name || "",
                        productName: product.product_name || "",
                        productDescription: product.description || product.product_description || "",
                        price: product.price || 0,
                        stock: product.stock || 0,
                        size: product.size || "Small",
                        inventory: inventoryMap,
                        popular: product.popular || false,
                        discountPercentage: product.offer?.discount_percentage || product.discount_percentage || 0,
                        startDate: formattedStartDate,
                        endDate: formattedEndDate,
                    });

                    setExistingImageUrl(product.product_image_url || product.image_url || "");
                }
            } catch (error) {
                console.error("Failed to fetch product details:", error);
                setSubmitError("Failed to load product details.");
            } finally {
                setIsFetching(false);
            }
        };

        fetchProduct();
    }, [id]);

    const onChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        const { name, value, type } = e.target;
        setFormData((prev) => ({
            ...prev,
            [name]: type === "checkbox" 
                ? (e.target as HTMLInputElement).checked 
                : type === "number" 
                ? Number(value) 
                : value,
        }));
        if (submitError) setSubmitError("");
    };

    const onInventoryChange = (size: string, value: number) => {
        setFormData((prev) => ({
            ...prev,
            inventory: {
                ...(prev.inventory ?? {}),
                [size]: value,
            },
        }));
        if (submitError) setSubmitError("");
    };

    const onImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0] ?? null;
        setImageFile(file);
    };

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setSubmitError("");

        const newErrors = validateForm(formData);

        if (imageFile && !imageFile.type.startsWith("image/")) {
            newErrors.productImage = "Please select a valid image file.";
        }

        if (Object.keys(newErrors).length > 0) {
            setErrors(newErrors);
            return;
        }

        setErrors({});
        setIsLoading(true);

        try {
            const data = new FormData();
            data.append("category_name", formData.categoryName);
            data.append("product_name", formData.productName);
            data.append("product_description", formData.productDescription);
            data.append("price", String(formData.price));
            data.append("popular", String(formData.popular));
            data.append("size", formData.size || "Small");
            data.append("stock", String(formData.inventory?.[formData.size] ?? 0));
            
            data.append(
                "inventory",
                JSON.stringify(
                    Object.entries(formData.inventory ?? {}).map(([size, stock]) => ({
                        size,
                        stock: Number(stock) || 0,
                    }))
                )
            );

            data.append("discount_percentage", String(formData.discountPercentage));
            if (formData.startDate) data.append("start_date", formData.startDate);
            if (formData.endDate) data.append("end_date", formData.endDate);

            if (imageFile) {
                data.append("product_image", imageFile);
            }

            const response = await axios.put(`${API_BASE_URL}/admin/product/${id}`, data, {
                withCredentials: true,
            });

            if (response.status === 200 || response.status === 201) {
                router.push("/admin/products");
            }
        } catch (error: unknown) {
            if (axios.isAxiosError(error)) {
                setSubmitError(
                    error.response?.data?.message || "Failed to update product. Please try again."
                );
            } else {
                setSubmitError("An unexpected error occurred.");
            }
            console.error("Failed to update product:", error);
        } finally {
            setIsLoading(false);
        }
    };

    if (isFetching) {
        return (
            <div className="flex min-h-screen items-center justify-center bg-slate-950 px-4 py-10 text-slate-300">
                <div className="rounded-2xl border border-white/10 bg-white/5 px-6 py-5 shadow-2xl shadow-slate-950/40 backdrop-blur-xl">
                    Loading product details...
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-slate-950 px-4 py-8 text-white sm:px-6 lg:px-8">
            <div className="mx-auto max-w-6xl">
                <header className="mb-8 rounded-[2rem] border border-white/10 bg-white/5 p-5 shadow-2xl shadow-slate-950/30 backdrop-blur-xl">
                    <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
                        <div>
                            <p className="text-xs font-semibold uppercase tracking-[0.3em] text-cyan-300">
                                Product management
                            </p>
                            <h1 className="mt-2 text-3xl font-bold text-white">Edit product</h1>
                        </div>

                        <div className="flex items-center gap-3">
                            <Link
                                href="/admin/products"
                                className="rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm font-medium text-slate-200 transition hover:border-cyan-400/40 hover:bg-cyan-500/10"
                            >
                                Back to products
                            </Link>
                        </div>
                    </div>
                </header>

                <div className="rounded-[2rem] border border-white/10 bg-slate-900/80 p-5 shadow-2xl shadow-slate-950/30 backdrop-blur-xl md:p-8">
                    <div className="mb-6 flex items-center gap-3">
                        <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-gradient-to-br from-cyan-500 to-violet-500 text-xl font-bold text-white shadow-lg shadow-cyan-500/20">
                            ✎
                        </div>
                        <div>
                            <h2 className="text-xl font-semibold text-white">Product details</h2>
                            <p className="text-sm text-slate-400">Update pricing, inventory, image, and offer data.</p>
                        </div>
                    </div>

                    {submitError && (
                        <div className="mb-6 rounded-xl border border-rose-500/40 bg-rose-500/10 px-4 py-3 text-sm text-rose-300">
                            {submitError}
                        </div>
                    )}

                    <form onSubmit={handleSubmit} className="w-full">
                        <ProductForm
                            formData={formData}
                            onChange={onChange}
                            onInventoryChange={onInventoryChange}
                            onImageChange={onImageChange}
                            isLoading={isLoading}
                            categories={categories}
                            errors={errors}
                            existingImageUrl={existingImageUrl}
                        />
                    </form>
                </div>
            </div>
        </div>
    );
}