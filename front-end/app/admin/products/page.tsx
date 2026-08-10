'use client'

import axios from 'axios';
import Link from 'next/link';
import { useEffect, useState } from 'react';

interface Product {
    id: number;
    created_at: string;
    category_id: number;
    category_name: string;
    product_name: string;
    product_description: string;
    product_image_url: string;
    price: number;
    stock: number;
    popular: boolean;
    size: string;
    has_offer: boolean;
    offer_discount_percent: number;
    discount_amount: number;
    total_discounted_amount: number;
}

export default function Products() {
    const [products, setProducts] = useState<Product[]>([]);
    const [sortOrder, setSortOrder] = useState("DSC");
    useEffect(() => {
        async function fetchData() {
            try {
                const response = await axios.get(
                    `http://localhost:8080/admin/product?list_order=${sortOrder}`,
                    {
                        withCredentials: true,
                    }
                );

                if (response.data?.status && response.data?.data?.products) {
                    setProducts(response.data.data.products);
                }
            } catch (error) {
                console.error('Failed to load products:', error);
            }
        }

        fetchData();
    }, [sortOrder]);

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
                prev.filter((product) => product.id !== id)
            );
        } catch (error) {
            console.error(error);
        }
    };

    return (
        <div className="min-h-screen bg-gray-100">
            {/* Header */}
            <div className="flex justify-between items-center bg-green-700 text-white px-8 py-6">
                <h1 className="text-3xl font-bold">
                    Products Management
                </h1>
                <select
                        
                    onChange={(e) => setSortOrder(e.target.value)}
                    className="bg-white text-black border rounded px-3 py-2"
                >
                    <option value="DSC">Newest First</option>
                    <option value="ASC">Oldest First</option>
                </select>
                <Link
                    href="/admin/products/add"
                    className="bg-blue-600 hover:bg-blue-700 px-5 py-2 rounded font-medium transition"
                >
                    + Add Product
                </Link>
            </div>

            {/* Table */}
            <div className="p-6">
                <div className="overflow-x-auto bg-white rounded-lg shadow border">

                    <table className="min-w-full text-sm">
                        <thead className="bg-gray-200 text-gray-700">
                            <tr>
                                <th className="px-3 py-3 text-center">S.No</th>
                                <th className="px-3 py-3">Created At</th>
                                <th className="px-3 py-3">Category</th>
                                <th className="px-3 py-3">Product</th>
                                <th className="px-3 py-3">Description</th>
                                <th className="px-3 py-3">Image URL</th>
                                <th className="px-3 py-3 text-right">Price</th>
                                <th className="px-3 py-3 text-center">Stock</th>
                                <th className="px-3 py-3 text-center">Size</th>
                                <th className="px-3 py-3 text-center">Popular</th>
                                <th className="px-3 py-3 text-center">Offer</th>
                                <th className="px-3 py-3 text-center">Discount %</th>
                                <th className="px-3 py-3 text-right">Discount</th>
                                <th className="px-3 py-3 text-right">Final Price</th>
                                <th className="px-3 py-3 text-center">Actions</th>
                            </tr>
                        </thead>

                        <tbody>
                            {products.length === 0 ? (
                                <tr>
                                    <td
                                        colSpan={15}
                                        className="text-center py-8 text-gray-500"
                                    >
                                        No products found.
                                    </td>
                                </tr>
                            ) : (
                                products.map((product, index) => (
                                    <tr
                                        key={product.id}
                                        className="border-t hover:bg-gray-50"
                                    >
                                        <td className="px-3 py-3 text-center">
                                            {index + 1}
                                        </td>

                                        <td className="px-3 py-3 whitespace-nowrap">
                                            {new Date(
                                                product.created_at
                                            ).toLocaleString('en-IN', {
                                                day: '2-digit',
                                                month: 'short',
                                                year: 'numeric',
                                                hour: '2-digit',
                                                minute: '2-digit',
                                            })}
                                        </td>

                                        <td className="px-3 py-3 whitespace-nowrap">
                                            {product.category_name}
                                        </td>

                                        <td className="px-3 py-3 whitespace-nowrap font-medium">
                                            {product.product_name}
                                        </td>

                                        <td
                                            className="px-3 py-3 max-w-xs truncate"
                                            title={product.product_description}
                                        >
                                            {product.product_description}
                                        </td>

                                        <td
                                            className="px-3 py-3 max-w-xs truncate text-blue-600"
                                            title={product.product_image_url}
                                        >
                                            {product.product_image_url || 'No image'}
                                        </td>

                                        <td className="px-3 py-3 text-right whitespace-nowrap">
                                            ₹{product.price}
                                        </td>

                                        <td className="px-3 py-3 text-center">
                                            {product.stock}
                                        </td>

                                        <td className="px-3 py-3 text-center">
                                            {product.size}
                                        </td>

                                        <td className="px-3 py-3 text-center">
                                            {product.popular ? 'Yes' : 'No'}
                                        </td>

                                        <td className="px-3 py-3 text-center">
                                            {product.has_offer ? 'Yes' : 'No'}
                                        </td>

                                        <td className="px-3 py-3 text-center">
                                            {product.has_offer
                                                ? `${product.offer_discount_percent}%`
                                                : '-'}
                                        </td>

                                        <td className="px-3 py-3 text-right">
                                            {product.has_offer
                                                ? `₹${product.discount_amount}`
                                                : '-'}
                                        </td>

                                        <td className="px-3 py-3 text-right font-semibold">
                                            {product.has_offer
                                                ? `₹${product.total_discounted_amount}`
                                                : `₹${product.price}`}
                                        </td>

                                        <td className="px-3 py-3">
                                            <div className="flex justify-center gap-2">
                                                <Link
                                                    href={`/admin/products/edit/${product.id}`}
                                                    className="bg-blue-500 hover:bg-blue-600 text-white px-3 py-1 rounded text-xs"
                                                >
                                                    Edit
                                                </Link>

                                                <button
                                                    onClick={() =>
                                                        handleDelete(product.id)
                                                    }
                                                    className="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded text-xs"
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