'use client'

import axios from "axios";
import Link from "next/link";
import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import ProductForm, { FormErrors } from "@/components/product-form";
import { useCategories } from "@/hooks/useCategories";
import { ProductFormData } from "@/types/productFormData";

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

    if (
        data.discountPercentage < 0 ||
        data.discountPercentage > 100
    ) {
        errors.discountPercentage =
            "Discount percentage must be between 0 and 100";
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

    /*
     * Product information
     *
     * Notice that image is NOT inside formData.
     * The existing image is stored separately and
     * the newly selected image is stored as a File.
     */
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
    });

    // Image currently stored in S3
    const [existingImageUrl, setExistingImageUrl] = useState<string>("");

    // New image selected by the user
    const [imageFile, setImageFile] = useState<File | null>(null);

    const [errors, setErrors] = useState<FormErrors>({});

    /*
     * Fetch existing product
     */
    useEffect(() => {
        if (!id) return;

        const fetchProduct = async () => {
            try {
                const response = await axios.get(
                    `http://localhost:8080/admin/product/${id}`,
                    {
                        withCredentials: true,
                    }
                );

                if (response.data?.data) {
                    const product = response.data.data;
                    const inventoryMap: Record<string, number> = {
                        Small: 0,
                        Medium: 0,
                        Large: 0,
                    };

                    (product.inventory || []).forEach((entry: { size: string; stock: number }) => {
                        if (entry.size && inventoryMap[entry.size] !== undefined) {
                            inventoryMap[entry.size] = entry.stock ?? 0;
                        }
                    });

                    setFormData({
                        categoryName:
                            product.category_name || "",

                        productName:
                            product.product_name || "",

                        productDescription:
                            product.product_description || "",

                        price:
                            product.price || 0,

                        stock:
                            product.stock || 0,

                        size:
                            product.size || "",

                        inventory: inventoryMap,

                        popular:
                            product.popular || false,

                        discountPercentage:
                            product.discount_percentage || 0,
                    });

                    /*
                     * Store the existing S3 image URL separately.
                     */
                    setExistingImageUrl(
                        product.product_image_url || ""
                    );
                }
            } catch (error) {
                console.error(
                    "Failed to fetch product details:",
                    error
                );
            } finally {
                setIsFetching(false);
            }
        };

        fetchProduct();
    }, [id]);

    /*
     * Handle text/select/number/checkbox fields
     */
    const onChange = (
        e: React.ChangeEvent<
            HTMLInputElement | HTMLSelectElement
        >
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

    const onInventoryChange = (size: string, value: number) => {
        setFormData((prev) => ({
            ...prev,
            inventory: {
                ...(prev.inventory ?? {}),
                [size]: value,
            },
        }));
    };

    /*
     * Handle new image selection
     */
    const onImageChange = (
        e: React.ChangeEvent<HTMLInputElement>
    ) => {
        const file = e.target.files?.[0] ?? null;

        setImageFile(file);
    };

    /*
     * Submit updated product
     */
    const handleSubmit = async (
        e: React.FormEvent<HTMLFormElement>
    ) => {
        e.preventDefault();

        /*
         * Validate normal fields
         */
        const newErrors = validateForm(formData);

        /*
         * IMPORTANT:
         *
         * Image is NOT required when editing.
         *
         * If imageFile is null:
         *     keep existing image.
         *
         * If imageFile exists:
         *     upload new image.
         */

        /*
         * Validate new image only if user selected one.
         */
        if (imageFile) {
            if (!imageFile.type.startsWith("image/")) {
                newErrors.productImage =
                    "Please select a valid image";
            }
        }

        if (Object.keys(newErrors).length > 0) {
            setErrors(newErrors);
            return;
        }

        setErrors({});
        setIsLoading(true);

        try {
            /*
             * Because we are uploading a file,
             * we MUST use FormData.
             */
            const data = new FormData();

            data.append(
                "category_name",
                formData.categoryName
            );

            data.append(
                "product_name",
                formData.productName
            );

            data.append(
                "product_description",
                formData.productDescription
            );

            data.append(
                "price",
                String(formData.price)
            );

            data.append(
                "stock",
                String(formData.stock)
            );

            data.append(
                "popular",
                String(formData.popular)
            );

            data.append(
                "inventory",
                JSON.stringify(
                    Object.entries(formData.inventory ?? {}).map(([size, stock]) => ({
                        size,
                        stock: Number(stock) || 0,
                    }))
                )
            );

            data.append(
                "discount_percentage",
                String(formData.discountPercentage)
            );

            /*
             * Only append product_image when
             * the user selected a NEW image.
             */
            if (imageFile) {
                data.append(
                    "product_image",
                    imageFile
                );
            }

            const response = await axios.put(
                `http://localhost:8080/admin/product/${id}`,
                data,
                {
                    withCredentials: true,
                }
            );

            if (
                response.status === 200 ||
                response.status === 201
            ) {
                router.push("/admin/products");
            }
        } catch (error) {
            console.error(
                "Failed to update product:",
                error
            );
        } finally {
            setIsLoading(false);
        }
    };

    /*
     * Loading existing product
     */
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
                            <h1 className="mt-2 text-3xl font-bold text-white">
                                Edit product
                            </h1>
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