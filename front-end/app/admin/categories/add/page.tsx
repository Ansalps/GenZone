'use client';

import CategoryForm from '@/components/category-from';
import axios from 'axios';
import { useRouter } from 'next/navigation';
import { useState } from 'react';

export default function AddCategory() {
    const router = useRouter();

    const [isLoading, setIsLoading] = useState(false);

    const [formData, setFormData] = useState({
        categoryName: '',
        categoryDescription: '',
        categoryImage: null as File | null,
    });

    // Text fields
    const handleChange = (
        e: React.ChangeEvent<HTMLInputElement>
    ) => {
        const { name, value } = e.target;

        setFormData((prevData) => ({
            ...prevData,
            [name]: value,
        }));
    };

    // Image field
    const handleImageChange = (
        e: React.ChangeEvent<HTMLInputElement>
    ) => {
        const file = e.target.files?.[0] || null;

        setFormData((prevData) => ({
            ...prevData,
            categoryImage: file,
        }));
    };

    const handleSubmit = async (
        e: React.FormEvent<HTMLFormElement>
    ) => {
        e.preventDefault();

        if (!formData.categoryImage) {
            alert('Please select an image file.');
            return;
        }

        setIsLoading(true);

        const data = new FormData();

        data.append(
            'category_name',
            formData.categoryName
        );

        data.append(
            'description',
            formData.categoryDescription
        );

        data.append(
            'image',
            formData.categoryImage
        );

        try {
            const response = await axios.post(
                'http://localhost:8080/admin/category',
                data,
                {
                    withCredentials: true,
                }
            );

            if (response.status === 200) {
                console.log('Success:', response.data);

                router.push('/admin/categories');
            }
        } catch (error) {
            console.error(
                'Error uploading category:',
                error
            );
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="min-h-screen bg-slate-950 px-4 py-8 text-white sm:px-6 lg:px-8">
            <div className="mx-auto max-w-7xl">
                <header className="mb-6 flex items-center justify-between gap-4 rounded-3xl border border-white/10 bg-white/5 p-5 shadow-2xl shadow-slate-950/30 backdrop-blur-xl">
                    <div>
                        <h1 className="text-2xl font-bold">Add Category</h1>
                        <p className="text-sm text-slate-300">Create a new category for your store</p>
                    </div>
                </header>

                <main className="rounded-3xl border border-white/10 bg-slate-900/60 p-6 shadow-2xl shadow-slate-950/30">
                    <div className="flex justify-center">
                        <form
                            onSubmit={handleSubmit}
                            className="w-full max-w-lg"
                        >
                            <CategoryForm
                                formData={formData}
                                onChange={handleChange}
                                onImageChange={handleImageChange}
                                isLoading={isLoading}
                                isEdit={false}
                            />
                        </form>
                    </div>
                </main>
            </div>
        </div>
    );
}