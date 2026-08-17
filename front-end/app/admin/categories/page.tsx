'use client'

import axios from 'axios';
import Link from 'next/link';
import { useEffect, useState } from 'react';
import CategoryCard from '@/components/admin/category-card';
import { Category } from '@/types/category';



export default function Categories() {
    const [categories, setCategories] = useState<Category[]>([]);
    const [sortOrder, setSortOrder] = useState("DSC");
    useEffect(() => {
        async function fetchData() {
            try {
                const response = await axios.get(
                    `http://localhost:8080/admin/category?list_order=${sortOrder}`,
                    {
                        withCredentials: true,
                    }
                );

                if (response.data?.status && response.data?.data?.categories) {
                    setCategories(response.data.data.categories);
                }
            } catch (error) {
                console.error('Failed to load categories:', error);
            }
        }

        fetchData();
    }, [sortOrder]);

    const handleDelete = async (id: number) => {
        const confirmed = window.confirm(
            'Are you sure you want to delete this category? This will delete all the products under this categoy.'
        );

        if (!confirmed) return;

        try {
            await axios.delete(
                `http://localhost:8080/admin/category/${id}`,
                {
                    withCredentials: true,
                }
            );

            setCategories((prev) =>
                prev.filter((category) => category.id !== id)
            );
        } catch (error) {
            console.error(error);
        }
    };

    return (
        <div className="min-h-screen bg-gray-100">
            {/* Header */}
            <div className="flex justify-between items-center bg-green-700 text-white px-8 py-6">
                <div className="flex items-center gap-4">
                    {/* Back Button */}
                    <Link
                        href="/admin"
                        className="bg-green-800 hover:bg-green-900 text-white px-4 py-2 rounded text-sm font-medium transition flex items-center gap-1 border border-green-600"
                    >
                        ← Back
                    </Link>
                    <h1 className="text-3xl font-bold">
                        Categories Management
                    </h1>
                </div>
                 <div className="flex items-center gap-4">
                    <select
                        
                        onChange={(e) => setSortOrder(e.target.value)}
                        className="bg-white text-black border rounded px-3 py-2"
                    >
                        <option value="DSC">Newest First</option>
                        <option value="ASC">Oldest First</option>
                    </select>

                    <Link
                        href="/admin/categories/add"
                        className="bg-blue-600 hover:bg-blue-700 px-4 py-2 rounded"
                    >
                        + Add Category
                    </Link>
                </div>
            </div>

            <div className="max-w-7xl mx-auto mt-8 px-4">
                <div className="overflow-x-auto bg-white rounded-lg shadow">
                    <div className="max-w-7xl mx-auto mt-8 px-4">
                        {categories.length === 0 ? (
                            <div className="bg-white rounded-lg shadow p-10 text-center text-gray-500">
                                No categories found.
                            </div>
                        ) : (
                            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
                                {categories.map((category) => (
                                    <CategoryCard
                                        key={category.id}
                                        category={category}
                                        onDelete={handleDelete}
                                    />
                                ))}
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
}