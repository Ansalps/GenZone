'use client'

import axios from "axios";
import { useEffect, useState } from "react";
import { useParams, useRouter } from 'next/navigation';
import ProductForm from "@/components/product-form";
import { useCategories } from "@/hooks/useCategories";

export default function EditProduct() {
    const { categories } = useCategories();
    const router = useRouter();
    const params = useParams();
    const id = params?.id as string;

    const [isLoading, setIsLoading] = useState(false);
    const [isFetching, setIsFetching] = useState(true);

    // Refactored state: removed hasOffer, discountAmount, totalDiscountedAmount
    const [formData, setFormData] = useState({
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

    useEffect(() => {
        if (!id) return;

        async function fetchProduct() {
            try {
                const response = await axios.get(
                    `http://localhost:8080/admin/product/${id}`,
                    { withCredentials: true }
                );

                if (response.data?.data) {
                    const product = response.data.data;

                    setFormData({
                        categoryName: product.category_name || '',
                        productName: product.product_name || '',
                        productDescription: product.product_description || '',
                        productImageUrl: product.product_image_url || '',
                        price: product.price || 0,
                        stock: product.stock || 0,
                        size: product.size || '',
                        popular: product.popular || false,
                        // Map discount_percentage directly from backend response
                        discountPercentage: product.discount_percentage || 0,
                    });
                }
            } catch (error) {
                console.error("Failed to fetch product details:", error);
            } finally {
                setIsFetching(false);
            }
        }

        fetchProduct();
    }, [id]);

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
        setIsLoading(true);

        try {
            const response = await axios.put(
                `http://localhost:8080/admin/product/${id}`,
                {
                    category_name: formData.categoryName,
                    product_name: formData.productName,
                    product_description: formData.productDescription,
                    product_image_url: formData.productImageUrl,
                    price: formData.price,
                    stock: formData.stock,
                    popular: formData.popular,
                    size: formData.size,
                    // Send discount_percentage directly
                    discount_percentage: formData.discountPercentage,
                },
                { withCredentials: true }
            );

            if (response.status === 200) {
                router.push('/admin/products');
            }
        } catch (error) {
            console.error("Failed to update product:", error);
        } finally {
            setIsLoading(false);
        }
    };

    if (isFetching) {
        return (
            <div className="flex justify-center items-center h-screen text-gray-500">
                Loading product details...
            </div>
        );
    }

    return (
        <div className="p-6 max-w-4xl mx-auto">
            <h1 className="font-bold text-2xl mb-6">Edit Product</h1>
            
            <form onSubmit={handleSubmit} className="w-full">
                <ProductForm
                    formData={formData}
                    onChange={onChange}
                    isLoading={isLoading}
                    categories={categories}
                />
            </form>
        </div>
    );
}