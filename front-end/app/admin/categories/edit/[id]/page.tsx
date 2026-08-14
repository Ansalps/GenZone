'use client';

import axios from 'axios';
import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import CategoryForm from '@/components/category-from';

export default function EditCategory() {
    const router = useRouter();
    const params = useParams();
    const id = params?.id as string;

    const [isLoading, setIsLoading] = useState(false);
    const [isFetching, setIsFetching] = useState(true);

    const [formData, setFormData] = useState({
        categoryName: '',
        categoryDescription: '',
        categoryImage: null as File | null,
        categoryImageUrl: '',
    });

    const [imageFile, setImageFile] = useState<File | null>(null);

    useEffect(() => {
        if (!id) return;

        async function fetchCategory() {
            try {
                const response = await axios.get(
                    `http://localhost:8080/admin/category/${id}`,
                    {
                        withCredentials: true,
                    }
                );

                if (response.data?.data) {
                    const category = response.data.data;

                    setFormData({
                        categoryName: category.category_name || '',
                        categoryDescription: category.category_description || '',
                        categoryImage: null,
                        categoryImageUrl: category.category_image_url || '',
                    });
                }
            } catch (error) {
                console.error(
                    'Failed to fetch category details:',
                    error
                );
            } finally {
                setIsFetching(false);
            }
        }

        fetchCategory();
    }, [id]);

    const handleChange = (
        e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
    ) => {
        const { name, value } = e.target;

        setFormData((prevData) => ({
            ...prevData,
            [name]: value,
        }));
    };

    const handleImageChange = (
        e: React.ChangeEvent<HTMLInputElement>
    ) => {
        const file = e.target.files?.[0] || null;

        setImageFile(file);
    };

    const handleSubmit = async (
        e: React.FormEvent<HTMLFormElement>
    ) => {
        e.preventDefault();

        setIsLoading(true);

        try {
            const data = new FormData();

            data.append(
                'category_name',
                formData.categoryName
            );

            data.append(
                'description',
                formData.categoryDescription
            );

            // Only send image if user selected a new one
            if (imageFile) {
                data.append('image', imageFile);
            }

            const response = await axios.put(
                `http://localhost:8080/admin/category/${id}`,
                data,
                {
                    withCredentials: true,
                }
            );

            if (response.status === 200) {
                router.push('/admin/categories');
            }
        } catch (error) {
            console.error(
                'Failed to update category:',
                error
            );
        } finally {
            setIsLoading(false);
        }
    };

    if (isFetching) {
        return (
            <div className="flex justify-center items-center h-screen text-gray-500">
                Loading category details...
            </div>
        );
    }

    return (
        <div className="p-6 max-w-xl mx-auto">
            <h1 className="font-bold text-2xl mb-6">
                Edit Category
            </h1>

            <form
                onSubmit={handleSubmit}
                className="w-full"
            >
                <CategoryForm
                    formData={formData}
                    onChange={handleChange}
                    onImageChange={handleImageChange}
                    isLoading={isLoading}
                />
            </form>
        </div>
    );
}