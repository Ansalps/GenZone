'use client'
import axios from 'axios';
import Link from 'next/link'; 
import { useEffect, useState } from 'react';

interface Product {
    id: number;
    category_id:number;
    category_name: string;
    product_name:string;
    category_description: string;
    category_image_url?: string; 
}

export default function Categories(){
    const [categories, setCategories] = useState<Product[]>([]);

    useEffect(()=>{
        async function fetchData() {
            try {
                const response = await axios.get('http://localhost:8080/admin/product', {
                    withCredentials: true
                });
                if (response.data?.status && response.data?.data?.products) {
                    setCategories(response.data.data.categories);
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

        setCategories(prev =>
            prev.filter(category => category.id !== id)
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
                    Categories Management
                </div>
                <Link 
                    href="/admin/categories/add" 
                    className="bg-blue-600 hover:bg-blue-700 text-sm font-medium text-amber-50 cursor-pointer py-2 px-4 rounded transition-colors"
                >
                    + Add Category
                </Link>
            </div>

            {/* Table Wrapper Block */}
            <div className="p-6 w-full max-w-6xl mx-auto flex flex-col gap-2">
                
                {/* 1. Adjusted Header Grid Columns to allocate 2 tracks for Actions */}
                <div className='grid grid-cols-12 bg-gray-200 p-3 rounded font-bold text-base text-gray-700 shadow-sm border border-gray-300'>
                    <div className='col-span-1 text-center'>S.NO.</div>
                    <div className='col-span-3'>Category Name</div>
                    <div className='col-span-4'>Category Description</div> {/* Changed from col-span-5 to 4 */}
                    <div className='col-span-2'>Category Image URL</div>   {/* Changed from col-span-3 to 2 */}
                    <div className='col-span-2 text-center'>Actions</div>   {/* 2 columns for buttons */}
                </div>

                {categories.length === 0 ? (
                    <div className="text-center text-gray-500 py-8 border border-dashed rounded mt-2">
                        No categories found. Click Add Category to create one.
                    </div>
                ) : (
                    categories.map((category, index) => (
                        /* 2. Matched data row layout with header grid allocation */
                        <div 
                            key={category.id} 
                            className='grid grid-cols-12 items-center bg-white p-3 rounded text-sm text-gray-600 border border-gray-200 hover:bg-gray-50 shadow-xs transition-colors'
                        >
                            <div className='col-span-1 text-center font-semibold text-gray-800'>
                                {index + 1}
                            </div>
                            
                            <div className='col-span-3 font-medium text-gray-900 capitalize'>
                                {category.category_name}
                            </div>
                            
                            <div className='col-span-4 pr-4 truncate' title={category.category_description}>
                                {category.category_description}
                            </div>
                            
                            <div className='col-span-2 truncate font-mono text-xs text-blue-500' title={category.category_image_url}>
                                {category.category_image_url || "No image"}
                            </div>

                            {/* 3. Actions Column with Next.js Link passing dynamic Category ID */}
                            <div className='col-span-2 flex justify-between items-center'>
                                <Link
                                    href={`/admin/categories/edit/${category.id}`}
                                    className="bg-blue-500 hover:bg-blue-600 text-white font-medium text-xs py-1.5 px-4 rounded transition-colors cursor-pointer shadow-xs"
                                >
                                    Edit
                                </Link>
                                <button
                                    onClick={() => handleDelete(category.id)}
                                    className="bg-red-500 hover:bg-red-600 text-white font-medium text-xs py-1.5 px-4 rounded transition-colors"
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
