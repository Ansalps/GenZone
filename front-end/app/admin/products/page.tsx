'use client'
import axios from 'axios';
import Link from 'next/link'; 
import { useEffect, useState } from 'react';



interface Product {
    id: number;
    category_id:number;
    category_name: string;
    product_name:string;
    product_description:string;
    product_image_url:string;
    price:number;
    stock:number;
    popular:boolean;
    size:string;
    has_offer:boolean;
    offer_discount_percent:number
}

export default function Categories(){
    

    const [products, setProducts] = useState<Product[]>([]);

    useEffect(()=>{
        async function fetchData() {
            try {
                const response = await axios.get('http://localhost:8080/admin/product', {
                    withCredentials: true
                });
                if (response.data?.status && response.data?.data?.products) {
                    setProducts(response.data.data.products);
                }
            } catch (error) {
                console.error("Failed to load categories:", error);
            }
        }
        fetchData();
    },[])

   const handleDelete = async (id: number) => {
    const confirmed = window.confirm(
        "Are you sure you want to delete this category?"
    );

    if (!confirmed) return;

    try {
        await axios.delete(
            `http://localhost:8080/admin/product/${id}`,
            {
                withCredentials: true,
            }
        );

        setProducts(prev =>
            prev.filter(product => product.id !== id)
        );
    } catch (error) {
        console.error(error);
    }
};

    return (
        <div className="w-screen min-h-screen bg-gray-50 flex flex-col">
            {/* Header Layout */}
            <div className="flex justify-between items-center w-full h-24 text-2xl bg-green-700 text-amber-50 p-6 border-b">
                <div className='font-bold'>
                    Products Management
                </div>
                <Link 
                    href="/admin/products/add" 
                    className="bg-blue-600 hover:bg-blue-700 text-sm font-medium text-amber-50 cursor-pointer py-2 px-4 rounded transition-colors"
                >
                    + Add Product
                </Link>
            </div>

            {/* Table Wrapper Block */}
            <div className="p-6 w-full max-w-6xl mx-auto flex flex-col gap-2">
                
                {/* 1. Adjusted Header Grid Columns to allocate 2 tracks for Actions */}
               <div className="grid grid-cols-20 bg-gray-200 p-3 rounded font-bold text-sm text-gray-700 border">
                    <div className="col-span-1 text-center">S.No</div>
                    <div className="col-span-2">Category</div>
                    <div className="col-span-2">Product</div>
                    <div className="col-span-4">Description</div>
                    <div className="col-span-2">Image URL</div>
                    <div className="col-span-1 text-center">Price</div>
                    <div className="col-span-1 text-center">Stock</div>
                    <div className="col-span-1 text-center">Popular</div>
                    <div className="col-span-1 text-center">Offer</div>
                    <div className="col-span-2 text-center">Discount %</div>
                    <div className="col-span-3 text-center">Actions</div>
                </div>

                {products.length === 0 ? (
                    <div className="text-center text-gray-500 py-8 border border-dashed rounded mt-2">
                        No Products found. Click Add Product to create one.
                    </div>
                ) : (
                    products.map((product, index) => (
                        /* 2. Matched data row layout with header grid allocation */
                        <div
                            key={product.id}
                            className="grid grid-cols-20 items-center bg-white p-3 rounded border hover:bg-gray-50 text-sm"
                        >
                            <div className="col-span-1 text-center">
                                {index + 1}
                            </div>

                            <div className="col-span-2 truncate">
                                {product.category_name}
                            </div>

                            <div className="col-span-2 truncate">
                                {product.product_name}
                            </div>

                            <div
                                className="col-span-4 truncate"
                                title={product.product_description}
                            >
                                {product.product_description}
                            </div>

                            <div
                                className="col-span-2 truncate text-blue-600"
                                title={product.product_image_url}
                            >
                                {product.product_image_url || "No image"}
                            </div>

                            <div className="col-span-1 text-center">
                                ₹{product.price}
                            </div>

                            <div className="col-span-1 text-center">
                                {product.stock}
                            </div>

                            <div className="col-span-1 text-center">
                                {product.popular ? "Yes" : "No"}
                            </div>

                            <div className="col-span-1 text-center">
                                {product.has_offer ? "Yes" : "No"}
                            </div>

                            <div className="col-span-2 text-center">
                                {product.offer_discount_percent}%
                            </div>

                            <div className="col-span-3 flex justify-center gap-2">
                                <Link
                                    href={`/admin/product/edit/${product.id}`}
                                    className="bg-blue-500 hover:bg-blue-600 text-white px-3 py-1 rounded text-xs"
                                >
                                    Edit
                                </Link>

                                <button
                                    onClick={() => handleDelete(product.id)}
                                    className="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded text-xs"
                                >
                                    Delete
                                </button>
                            </div>
                        </div>
                    ))
                )}
            </div>
        </div>
    )
}
