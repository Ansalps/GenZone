'use client'

import axios from 'axios';
import Link from 'next/link';
import { useEffect, useState } from 'react';
import { Product } from '@/types/product';
import ProductCard from '@/components/admin/product-card';
import { useCategories } from '@/hooks/useCategories';

export default function Products() {
    const [products, setProducts] = useState<Product[]>([]);
    const [sortOrder, setSortOrder] = useState("DSC");
    const [category, setCategory] = useState("");

    const {
        categories,
        isCategoriesLoading,
        categoriesError,
    } = useCategories();

    useEffect(() => {
        async function fetchData() {
            try {
                const response = await axios.get(
                    "http://localhost:8080/admin/product",
                    {
                        params: {
                            list_order: sortOrder,
                            category: category || undefined,
                        },
                        withCredentials: true,
                    }
                );

                if (
                    response.data?.status &&
                    response.data?.data?.products
                ) {
                    setProducts(response.data.data.products);
                }
            } catch (error) {
                console.error(
                    'Failed to load products:',
                    error
                );
            }
        }

        fetchData();
    }, [sortOrder, category]);

    const handleDelete = async (id: number) => {
        const confirmed = window.confirm(
            'Are you sure you want to delete this product?'
        );

        if (!confirmed) return;

        try {
            await axios.delete(
                `http://localhost:8080/admin/product/${id}`,
                {
                    withCredentials: true,
                }
            );

            setProducts((prev) =>
                prev.filter(
                    (product) => product.id !== id
                )
            );
        } catch (error) {
            console.error(
                'Failed to delete product:',
                error
            );
        }
    };

    return (
        <div className="min-h-screen bg-gray-100">

            {/* Header */}
            <div className="flex justify-between items-center bg-green-700 text-white px-8 py-6">

                <div className="flex items-center gap-4">

                    <Link
                        href="/admin"
                        className="bg-green-800 hover:bg-green-900 text-white px-4 py-2 rounded text-sm font-medium transition border border-green-600"
                    >
                        ← Back
                    </Link>

                    <h1 className="text-3xl font-bold">
                        Products Management
                    </h1>

                </div>

                <div className="flex items-center gap-4">

                    {/* Category Filter */}
                    <select
                        value={category}
                        onChange={(e) => setCategory(e.target.value)}
                        disabled={isCategoriesLoading}
                        className="bg-white text-black border rounded px-3 py-2"
                    >
                        <option value="">
                            {isCategoriesLoading
                                ? "Loading categories..."
                                : "All Categories"}
                        </option>

                        {categories.map((cat) => (
                            <option
                                key={cat.id}
                                value={cat.category_name}
                            >
                                {cat.category_name}
                            </option>
                        ))}
                    </select>

                    {/* Sort */}
                    <select
                        value={sortOrder}
                        onChange={(e) =>
                            setSortOrder(e.target.value)
                        }
                        className="bg-white text-black border rounded px-3 py-2"
                    >
                        <option value="DSC">
                            Newest First
                        </option>

                        <option value="ASC">
                            Oldest First
                        </option>
                    </select>

                    {/* Add Product */}
                    <Link
                        href="/admin/products/add"
                        className="bg-blue-600 hover:bg-blue-700 px-5 py-2 rounded font-medium transition"
                    >
                        + Add Product
                    </Link>

                </div>
            </div>

            {/* Products */}
            <div className="p-6">

                {products.length === 0 ? (
                    <div className="bg-white rounded-lg border p-12 text-center text-gray-500">
                        No products found.
                    </div>
                ) : (
                    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">

                        {products.map((product) => (
                            <ProductCard
                                key={product.id}
                                product={product}
                                onDelete={handleDelete}
                            />
                        ))}

                    </div>
                )}

            </div>

        </div>
    );
}