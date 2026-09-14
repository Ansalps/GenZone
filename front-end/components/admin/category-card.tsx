'use client';

import Link from 'next/link';
import { Category } from '@/types/category';



interface CategoryCardProps {
    category: Category;
    onRequestDelete: (id: number, name?: string) => void;
}

export default function CategoryCard({
    category,
    onRequestDelete,
}: CategoryCardProps) {
    return (
        <div className="rounded-2xl border border-white/10 bg-slate-900/80 shadow-lg overflow-hidden transition hover:shadow-2xl">
            {/* Image */}
            <div className="h-44 bg-gray-800">
                {category.category_image_url ? (
                    <img
                        src={category.category_image_url}
                        alt={category.category_name}
                        className="w-full h-full object-cover"
                    />
                ) : (
                    <div className="w-full h-full flex items-center justify-center text-slate-400">
                        No Image
                    </div>
                )}
            </div>

            {/* Content */}
            <div className="p-4">
                <h2 className="text-lg font-semibold text-white mb-1">
                    {category.category_name}
                </h2>

                <p
                    className="text-slate-300 text-sm line-clamp-3 mb-3"
                    title={category.category_description}
                >
                    {category.category_description}
                </p>

                <p className="text-xs text-slate-400 mb-4">
                    Created:{' '}
                    {new Date(category.category_created_at).toLocaleDateString('en-IN', {
                        day: '2-digit',
                        month: 'short',
                        year: 'numeric',
                    })}
                </p>

                {/* Actions */}
                <div className="flex gap-2">
                    <Link
                        href={`/admin/categories/edit/${category.id}`}
                        className="flex-1 text-center rounded-xl bg-gradient-to-r from-violet-500 to-purple-500 text-white text-sm px-3 py-2 transition hover:opacity-95"
                    >
                        Edit
                    </Link>

                    <button
                        onClick={() => onRequestDelete(category.id, category.category_name)}
                        className="flex-1 rounded-xl bg-red-600/90 text-white text-sm px-3 py-2 transition hover:opacity-95"
                    >
                        Delete
                    </button>
                </div>
            </div>
        </div>
    );
}