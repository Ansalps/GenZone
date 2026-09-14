'use client'

import axios from 'axios';
import Link from 'next/link';
import { useEffect, useState } from 'react';
import { Product } from '@/types/product';
import ProductCard from '@/components/admin/product-card';
import ConfirmModal from '@/components/admin/confirm-modal';
import { toast } from 'sonner';
import { useCategories } from '@/hooks/useCategories';

export default function Products() {
    const [products, setProducts] = useState<Product[]>([]);
    const [sortOrder, setSortOrder] = useState("DSC");
    const [category, setCategory] = useState("");

    const {
        categories,
        isCategoriesLoading,
        categoriesError,
    } = useCategories();

    useEffect(() => {
        async function fetchData() {
            try {
                const response = await axios.get(
                    "http://localhost:8080/admin/product",
                    {
                        params: {
                            list_order: sortOrder,
                            category: category || undefined,
                        },
                        withCredentials: true,
                    }
                );

                if (
                    response.data?.status &&
                    response.data?.data?.products
                ) {
                    setProducts(response.data.data.products);
                }
            } catch (error) {
                console.error(
                    'Failed to load products:',
                    error
                );
            }
        }

        fetchData();
    }, [sortOrder, category]);

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
            await axios.delete(`http://localhost:8080/admin/product/${confirmId}`, { withCredentials: true });
            setProducts((prev) => prev.filter((p) => p.id !== confirmId));
            toast.success('Product deleted');
        } catch (error) {
            console.error('Failed to delete product:', error);
            toast.error('Failed to delete product');
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
                        <Link href="/admin" className="rounded-full border border-white/10 bg-slate-900/60 px-3 py-2 text-sm font-medium text-slate-200 hover:bg-slate-900 transition">← Back</Link>

                        <h1 className="text-2xl font-bold">Products Management</h1>
                    </div>

                    <div className="flex items-center gap-3">
                        <select value={category} onChange={(e) => setCategory(e.target.value)} disabled={isCategoriesLoading} className="rounded-xl border border-white/10 bg-slate-900/50 px-3 py-2 text-sm text-white outline-none">
                            <option value="">{isCategoriesLoading ? 'Loading categories...' : 'All Categories'}</option>
                            {categories.map((cat) => (<option key={cat.id} value={cat.category_name}>{cat.category_name}</option>))}
                        </select>

                        <select value={sortOrder} onChange={(e) => setSortOrder(e.target.value)} className="rounded-xl border border-white/10 bg-slate-900/50 px-3 py-2 text-sm text-white outline-none">
                            <option value="DSC">Newest First</option>
                            <option value="ASC">Oldest First</option>
                        </select>

                        <Link href="/admin/products/add" className="rounded-xl bg-gradient-to-r from-emerald-500 to-teal-500 px-4 py-2 text-sm font-semibold text-white shadow-lg shadow-emerald-500/20 hover:brightness-105 transition">+ Add Product</Link>
                    </div>
                </header>

                <main className="rounded-3xl border border-white/10 bg-slate-900/60 p-6 shadow-2xl shadow-slate-950/30">
                    {products.length === 0 ? (
                        <div className="rounded-xl border border-white/6 bg-slate-900/50 p-12 text-center text-slate-300">No products found.</div>
                    ) : (
                        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
                            {products.map((product) => (
                                <ProductCard key={product.id} product={product} onRequestDelete={requestDelete} />
                            ))}
                        </div>
                    )}

                    <ConfirmModal open={confirmOpen} title={`Delete product${confirmName ? `: ${confirmName}` : ''}`} description="This will delete the product permanently." confirmLabel="Delete" cancelLabel="Cancel" loading={isDeleting} onCancel={() => setConfirmOpen(false)} onConfirm={handleDelete} />
                </main>
            </div>

        </div>
    );
}