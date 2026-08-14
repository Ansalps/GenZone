'use client';

import Link from 'next/link';
import { Category } from '@/types/category';



interface CategoryCardProps {
    category: Category;
    onDelete: (id: number) => void;
}

export default function CategoryCard({
    category,
    onDelete,
}: CategoryCardProps) {
    return (
        <div className="bg-white rounded-xl shadow-md overflow-hidden border border-gray-200 hover:shadow-lg transition">
            {/* Image */}
            <div className="h-48 bg-gray-100">
                {category.category_image_url ? (
                    <img
                        src={category.category_image_url}
                        alt={category.category_name}
                        className="w-full h-full object-cover"
                    />
                ) : (
                    <div className="w-full h-full flex items-center justify-center text-gray-400">
                        No Image
                    </div>
                )}
            </div>

            {/* Content */}
            <div className="p-5">
                <h2 className="text-xl font-semibold text-gray-800 mb-2">
                    {category.category_name}
                </h2>

                <p
                    className="text-gray-600 text-sm line-clamp-3 mb-4"
                    title={category.category_description}
                >
                    {category.category_description}
                </p>

                <p className="text-xs text-gray-400 mb-4">
                    Created:{' '}
                    {new Date(
                        category.category_created_at
                    ).toLocaleDateString('en-IN', {
                        day: '2-digit',
                        month: 'short',
                        year: 'numeric',
                    })}
                </p>

                {/* Actions */}
                <div className="flex gap-2">
                    <Link
                        href={`/admin/categories/edit/${category.id}`}
                        className="flex-1 text-center bg-blue-500 hover:bg-blue-600 text-white text-sm px-3 py-2 rounded transition"
                    >
                        Edit
                    </Link>

                    <button
                        onClick={() => onDelete(category.id)}
                        className="flex-1 bg-red-500 hover:bg-red-600 text-white text-sm px-3 py-2 rounded transition"
                    >
                        Delete
                    </button>
                </div>
            </div>
        </div>
    );
}