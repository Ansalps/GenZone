'use client'

import axios from 'axios';
import Link from 'next/link';
import { useEffect, useState } from 'react';

interface Category {
    id: number;
    category_created_at: string;
    category_name: string;
    category_description: string;
    category_image_url?: string;
}

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
            'Are you sure you want to delete this category?'
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
                    <table className="min-w-full table-auto">
                        <thead className="bg-gray-200">
                            <tr>
                                <th className="px-4 py-3 text-center">
                                    S.NO.
                                </th>

                                <th className="px-4 py-3 text-left">
                                    Created At
                                </th>

                                <th className="px-4 py-3 text-left">
                                    Category Name
                                </th>

                                <th className="px-4 py-3 text-left">
                                    Description
                                </th>

                                <th className="px-4 py-3 text-left">
                                    Image URL
                                </th>

                                <th className="px-4 py-3 text-center">
                                    Actions
                                </th>
                            </tr>
                        </thead>

                        <tbody>
                            {categories.length === 0 ? (
                                <tr>
                                    <td
                                        colSpan={6}
                                        className="text-center py-8 text-gray-500"
                                    >
                                        No categories found.
                                    </td>
                                </tr>
                            ) : (
                                categories.map((category, index) => (
                                    <tr
                                        key={category.id}
                                        className="border-t hover:bg-gray-50"
                                    >
                                        <td className="px-4 py-4 text-center font-semibold">
                                            {index + 1}
                                        </td>

                                        <td className="px-4 py-4 whitespace-nowrap">
                                            {new Date(
                                                category.category_created_at
                                            ).toLocaleString('en-IN', {
                                                day: '2-digit',
                                                month: 'short',
                                                year: 'numeric',
                                                hour: '2-digit',
                                                minute: '2-digit',
                                            })}
                                        </td>

                                        <td className="px-4 py-4 font-medium">
                                            {category.category_name}
                                        </td>

                                        <td
                                            className="px-4 py-4 max-w-sm truncate"
                                            title={category.category_description}
                                        >
                                            {category.category_description}
                                        </td>

                                        <td
                                            className="px-4 py-4 max-w-xs truncate text-blue-600"
                                            title={category.category_image_url}
                                        >
                                            {category.category_image_url ??
                                                'No Image'}
                                        </td>

                                        <td className="px-4 py-4">
                                            <div className="flex justify-center gap-2">
                                                <Link
                                                    href={`/admin/categories/edit/${category.id}`}
                                                    className="bg-blue-500 hover:bg-blue-600 text-white text-sm px-3 py-1 rounded"
                                                >
                                                    Edit
                                                </Link>

                                                <button
                                                    onClick={() =>
                                                        handleDelete(category.id)
                                                    }
                                                    className="bg-red-500 hover:bg-red-600 text-white text-sm px-3 py-1 rounded"
                                                >
                                                    Delete
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
}