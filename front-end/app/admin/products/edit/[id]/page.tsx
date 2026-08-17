'use client'

import axios from "axios";
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

    if (data.stock < 0 || !Number.isInteger(data.stock)) {
        errors.stock = "Stock must be a non-negative whole number";
    }

    if (!data.size) {
        errors.size = "Please select a size";
    } else if (
        !["Small", "Medium", "Large"].includes(data.size)
    ) {
        errors.size = "Size must be Small, Medium, or Large";
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
                "size",
                formData.size
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
            <div className="flex justify-center items-center h-screen text-gray-500">
                Loading product details...
            </div>
        );
    }

    return (
        <div className="p-6 max-w-4xl mx-auto">
            <h1 className="font-bold text-2xl mb-6">
                Edit Product
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
                    existingImageUrl={existingImageUrl}
                />
            </form>
        </div>
    );
}