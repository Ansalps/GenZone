'use client'

import axios from 'axios';
import Link from 'next/link';
import { useEffect, useState } from 'react';
import CategoryCard from '@/components/admin/category-card';
import ConfirmModal from '@/components/admin/confirm-modal';
import { toast } from 'sonner';
import { Category } from '@/types/category';



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

    const [confirmOpen, setConfirmOpen] = useState(false);
    const [confirmId, setConfirmId] = useState<number | null>(null);
    const [confirmName, setConfirmName] = useState<string | undefined>(undefined);
    const [isDeleting, setIsDeleting] = useState(false);

    const requestDelete = (id: number, name?: string) => {
        setConfirmId(id);
        setConfirmName(name);
        setConfirmOpen(true);
    };

    const handleDelete = async () => {
        if (!confirmId) return;

        setIsDeleting(true);

        try {
            await axios.delete(
                `http://localhost:8080/admin/category/${confirmId}`,
                { withCredentials: true }
            );

            setCategories((prev) => prev.filter((c) => c.id !== confirmId));
            toast.success('Category deleted');
        } catch (error) {
            console.error(error);
            toast.error('Failed to delete category');
        } finally {
            setIsDeleting(false);
            setConfirmOpen(false);
            setConfirmId(null);
            setConfirmName(undefined);
        }
    };

    return (
        <div className="min-h-screen bg-slate-950 px-4 py-8 text-white sm:px-6 lg:px-8">
            <div className="mx-auto max-w-7xl">
                <header className="mb-6 flex flex-col gap-4 rounded-3xl border border-white/10 bg-white/5 p-5 shadow-2xl shadow-slate-950/30 backdrop-blur-xl sm:flex-row sm:items-center sm:justify-between">
                    <div className="flex items-center gap-4">
                        <Link
                            href="/admin"
                            className="rounded-full border border-white/10 bg-slate-900/60 px-3 py-2 text-sm font-medium text-slate-200 hover:bg-slate-900 transition"
                        >
                            ← Back
                        </Link>

                        <h1 className="text-2xl font-bold">Categories Management</h1>
                    </div>

                    <div className="flex items-center gap-3">
                        <select
                            onChange={(e) => setSortOrder(e.target.value)}
                            className="rounded-xl border border-white/10 bg-slate-900/50 px-3 py-2 text-sm text-white outline-none"
                            value={sortOrder}
                        >
                            <option value="DSC">Newest First</option>
                            <option value="ASC">Oldest First</option>
                        </select>

                        <Link
                            href="/admin/categories/add"
                            className="rounded-xl bg-gradient-to-r from-cyan-500 to-violet-500 px-4 py-2 text-sm font-semibold text-white shadow-lg shadow-cyan-500/20 hover:brightness-105 transition"
                        >
                            + Add Category
                        </Link>
                    </div>
                </header>

                <main className="rounded-3xl border border-white/10 bg-slate-900/60 p-6 shadow-2xl shadow-slate-950/30">
                    {categories.length === 0 ? (
                        <div className="rounded-xl border border-white/6 bg-slate-900/50 p-12 text-center text-slate-300">
                            No categories found.
                        </div>
                    ) : (
                        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
                            {categories.map((category) => (
                                <CategoryCard
                                    key={category.id}
                                    category={category}
                                    onRequestDelete={requestDelete}
                                />
                            ))}
                        </div>
                    )}
                    <ConfirmModal
                        open={confirmOpen}
                        title={`Delete category${confirmName ? `: ${confirmName}` : ''}`}
                        description="This will delete the category and all products under it. This action is irreversible."
                        confirmLabel="Delete"
                        cancelLabel="Cancel"
                        loading={isDeleting}
                        onCancel={() => setConfirmOpen(false)}
                        onConfirm={handleDelete}
                    />
                </main>
            </div>
        </div>
    );
}