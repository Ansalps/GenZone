'use client'
import CategoryForm from "@/components/category-from"
import axios from "axios"
import { useRouter } from "next/navigation"
import { useState } from "react"

export default function AddCategory() {
    const router = useRouter();
    const [isLoading, setIsLoading] = useState(false);

    // 1. Update state to support string fields and File | null for image
    const [formData, setFormData] = useState<{
        categoryName: string;
        categoryDescription: string;
        categoryImage: File | null;
    }>({
        categoryName: '',
        categoryDescription: '',
        categoryImage: null,
    });

    // 2. Handle both text input and file input events
    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value, files, type } = e.target;

        if (type === 'file' && files && files[0]) {
            setFormData((prevData) => ({
                ...prevData,
                [name]: files[0], // Store raw file object
            }));
        } else {
            setFormData((prevData) => ({
                ...prevData,
                [name]: value,
            }));
        }
    };

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        if (!formData.categoryImage) {
            alert("Please select an image file.");
            return;
        }

        setIsLoading(true);

        // 3. Construct FormData payload for multipart submission
        const data = new FormData();
        data.append('category_name', formData.categoryName);
        data.append('description', formData.categoryDescription);
        data.append('image', formData.categoryImage); // Matches c.FormFile("image") in Go

        try {
            const response = await axios.post(
                'http://localhost:8080/admin/category',
                data,
                {
                    withCredentials: true,
                    headers: {
                        'Content-Type': 'multipart/form-data',
                    },
                }
            );

            if (response.status === 200) {
                console.log('Success:', response.data);
                router.push("/admin/categories");
            }
        } catch (error) {
            console.error('Error uploading category:', error);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <>
            <div className="font-bold text-2xl">
                Add Category
            </div>
            <div className="flex flex-col justify-center items-center h-screen gap-4">
                <form onSubmit={handleSubmit} className="flex flex-col gap-4 w-full max-w-xs">
                    <CategoryForm
                        formData={formData} 
                        onChange={handleChange}
                        isLoading={isLoading}
                    />
                </form>
            </div>
        </>
    );
}