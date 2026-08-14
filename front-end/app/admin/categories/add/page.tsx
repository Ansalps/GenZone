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
        <div className="p-6">
            <div className="font-bold text-2xl mb-6">
                Add Category
            </div>

            <div className="flex justify-center">
                <form
                    onSubmit={handleSubmit}
                    className="flex flex-col gap-4 w-full max-w-md"
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
        </div>
    );
}