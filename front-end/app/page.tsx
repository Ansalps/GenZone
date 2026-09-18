'use client'

import axios from 'axios'
import Link from 'next/link'
import { useEffect, useState } from 'react'

import CategoryCard from '@/components/landing/category-card'
import ProductCard from '@/components/landing/product-card'

interface Category {
    id: number
    category_name: string
    category_description: string
    category_image_url: string
}

interface Product {
    id: number
    category_id: number
    category_name: string
    product_name: string
    product_description: string
    product_image_url: string
    price: number
    stock: number
    popular: boolean
    size: string
    discount_percentage?: number
}

export default function LandingPage() {
    const [categories, setCategories] = useState<Category[]>([])
    const [products, setProducts] = useState<Product[]>([])
    const [search, setSearch] = useState('')
    const [isLoadingCategories, setIsLoadingCategories] = useState(true)
    const [isLoadingProducts, setIsLoadingProducts] = useState(true)

    useEffect(() => {
        const fetchCategories = async () => {
            try {
                const response = await axios.get('http://localhost:8080/public/category')

                if (response.data?.status && response.data?.data?.categories) {
                    setCategories(response.data.data.categories)
                }
            } catch (error) {
                console.error('Failed to load categories:', error)
            } finally {
                setIsLoadingCategories(false)
            }
        }

        fetchCategories()
    }, [])

    useEffect(() => {
        const fetchProducts = async () => {
            try {
                setIsLoadingProducts(true)

                const response = await axios.get('http://localhost:8080/public/product', {
                    params: {
                        search: search || undefined,
                    },
                })

                if (response.data?.status && response.data?.data?.products) {
                    setProducts(response.data.data.products)
                } else {
                    setProducts([])
                }
            } catch (error) {
                console.error('Failed to load products:', error)
                setProducts([])
            } finally {
                setIsLoadingProducts(false)
            }
        }

        fetchProducts()
    }, [search])

    return (
        <main className="min-h-screen bg-slate-950 text-white">
            <div className="absolute inset-0 -z-10 overflow-hidden">
                <div className="absolute -left-10 top-16 h-72 w-72 rounded-full bg-cyan-500/20 blur-3xl" />
                <div className="absolute right-0 top-40 h-80 w-80 rounded-full bg-violet-500/20 blur-3xl" />
                <div className="absolute bottom-20 left-1/3 h-64 w-64 rounded-full bg-emerald-500/10 blur-3xl" />
            </div>

            <header className="sticky top-0 z-20 border-b border-white/10 bg-slate-950/80 backdrop-blur-xl">
                <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-4">
                    <Link href="/" className="flex items-center gap-3">
                        <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-cyan-400 to-violet-500 font-black text-slate-950 shadow-lg shadow-cyan-500/30">
                            G
                        </div>
                        <div>
                            <p className="text-lg font-bold tracking-tight">GenZone</p>
                            <p className="text-[10px] uppercase tracking-[0.25em] text-slate-400">Store</p>
                        </div>
                    </Link>

                    <div className="flex items-center gap-3">
                        <Link href="/login" className="rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-cyan-400/40 hover:bg-cyan-500/10">
                            Log in
                        </Link>
                        <Link href="/signup" className="rounded-xl bg-gradient-to-r from-cyan-500 to-violet-500 px-4 py-2 text-sm font-semibold text-white shadow-lg shadow-cyan-500/20 transition hover:brightness-110">
                            Sign up
                        </Link>
                    </div>
                </div>
            </header>

            <section className="mx-auto max-w-7xl px-6 py-14 md:py-20">
                <div className="grid items-center gap-10 lg:grid-cols-[1.2fr_0.8fr]">
                    <div>
                        <span className="inline-flex rounded-full border border-cyan-400/30 bg-cyan-500/10 px-3 py-1 text-xs font-semibold uppercase tracking-[0.25em] text-cyan-200">
                            New season arrivals
                        </span>
                        <h1 className="mt-6 max-w-xl text-4xl font-black tracking-tight text-white md:text-6xl">
                            Make your everyday style feel premium.
                        </h1>
                        <p className="mt-5 max-w-xl text-base text-slate-300 md:text-lg">
                            Discover curated essentials, trending looks, and everyday favorites built for comfort, value, and confidence.
                        </p>

                        <div className="mt-8 flex flex-wrap gap-4">
                            <Link href="/signup" className="rounded-xl bg-gradient-to-r from-cyan-500 to-violet-500 px-6 py-3 font-semibold text-white shadow-lg shadow-cyan-500/20 transition hover:brightness-110">
                                Shop now
                            </Link>
                            <Link href="/login" className="rounded-xl border border-white/10 bg-white/5 px-6 py-3 font-semibold text-slate-100 transition hover:border-cyan-400/40 hover:bg-cyan-500/10">
                                Member login
                            </Link>
                        </div>

                        <div className="mt-9 grid max-w-lg grid-cols-3 gap-4">
                            <div className="rounded-2xl border border-white/10 bg-white/5 p-4">
                                <p className="text-2xl font-bold text-white">2K+</p>
                                <p className="mt-1 text-xs text-slate-400">Happy shoppers</p>
                            </div>
                            <div className="rounded-2xl border border-white/10 bg-white/5 p-4">
                                <p className="text-2xl font-bold text-white">180+</p>
                                <p className="mt-1 text-xs text-slate-400">Top picks</p>
                            </div>
                            <div className="rounded-2xl border border-white/10 bg-white/5 p-4">
                                <p className="text-2xl font-bold text-white">4.9/5</p>
                                <p className="mt-1 text-xs text-slate-400">Customer rating</p>
                            </div>
                        </div>
                    </div>

                    <div className="relative">
                        <div className="absolute -inset-4 rounded-[2rem] bg-gradient-to-br from-cyan-500/20 to-violet-500/20 blur-2xl" />
                        <div className="relative overflow-hidden rounded-[2rem] border border-white/10 bg-slate-900/80 p-5 shadow-2xl shadow-slate-950/50 backdrop-blur">
                            <div className="rounded-[1.5rem] bg-gradient-to-br from-slate-800 to-slate-900 p-5">
                                <div className="mb-5 flex items-center justify-between">
                                    <div>
                                        <p className="text-xs uppercase tracking-[0.25em] text-cyan-300">Featured</p>
                                        <h2 className="mt-2 text-2xl font-bold text-white">Trending collection</h2>
                                    </div>
                                    <span className="rounded-full border border-emerald-400/30 bg-emerald-500/10 px-3 py-1 text-xs font-medium text-emerald-200">
                                        In stock
                                    </span>
                                </div>

                                <div className="grid gap-4">
                                    <div className="overflow-hidden rounded-2xl border border-white/10 bg-slate-950/60">
                                        <img
                                            src="https://images.unsplash.com/photo-1521572267360-ee0c2909d518?auto=format&fit=crop&w=900&q=80"
                                            alt="Featured product"
                                            className="h-64 w-full object-cover"
                                        />
                                    </div>

                                    <div className="grid gap-3 sm:grid-cols-2">
                                        <div className="rounded-2xl border border-white/10 bg-white/5 p-4">
                                            <p className="text-xs uppercase tracking-[0.2em] text-slate-400">Category</p>
                                            <h3 className="mt-2 text-lg font-semibold text-white">Menswear</h3>
                                        </div>
                                        <div className="rounded-2xl border border-white/10 bg-white/5 p-4">
                                            <p className="text-xs uppercase tracking-[0.2em] text-slate-400">Best price</p>
                                            <h3 className="mt-2 text-lg font-semibold text-white">₹1,499</h3>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </section>

            <section className="mx-auto max-w-7xl px-6 py-6">
                <div className="rounded-[2rem] border border-white/10 bg-white/5 p-4 shadow-2xl shadow-slate-950/30 backdrop-blur-xl">
                    <label htmlFor="product-search" className="mb-2 block text-xs font-semibold uppercase tracking-[0.25em] text-slate-400">
                        Search the catalog
                    </label>
                    <div className="flex items-center gap-3 rounded-2xl border border-white/10 bg-slate-900/80 px-4 py-3">
                        <span className="text-xl text-cyan-300">⌕</span>
                        <input
                            id="product-search"
                            type="text"
                            value={search}
                            onChange={(e) => setSearch(e.target.value)}
                            placeholder="Search products, categories, styles..."
                            className="w-full bg-transparent text-base text-white placeholder:text-slate-500 focus:outline-none"
                        />
                    </div>
                </div>
            </section>

            <section className="mx-auto max-w-7xl px-6 py-12">
                <div className="mb-6 flex items-end justify-between gap-4">
                    <div>
                        <p className="text-xs font-semibold uppercase tracking-[0.25em] text-cyan-300">Browse</p>
                        <h2 className="mt-2 text-3xl font-bold text-white">Shop by category</h2>
                    </div>
                </div>

                {isLoadingCategories ? (
                    <div className="rounded-2xl border border-white/10 bg-white/5 p-10 text-center text-slate-300">
                        Loading categories...
                    </div>
                ) : categories.length === 0 ? (
                    <div className="rounded-2xl border border-white/10 bg-white/5 p-10 text-center text-slate-300">
                        No categories available right now.
                    </div>
                ) : (
                    <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
                        {categories.map((category) => (
                            <CategoryCard key={category.id} category={category} />
                        ))}
                    </div>
                )}
            </section>

            <section className="bg-slate-900/60 py-12">
                <div className="mx-auto max-w-7xl px-6">
                    <div className="mb-6 flex items-end justify-between gap-4">
                        <div>
                            <p className="text-xs font-semibold uppercase tracking-[0.25em] text-cyan-300">Collection</p>
                            <h2 className="mt-2 text-3xl font-bold text-white">{search ? 'Search results' : 'Featured products'}</h2>
                        </div>
                        {search && <p className="text-sm text-slate-400">Showing results for “{search}”</p>}
                    </div>

                    {isLoadingProducts ? (
                        <div className="rounded-2xl border border-white/10 bg-white/5 p-10 text-center text-slate-300">
                            Searching products...
                        </div>
                    ) : products.length === 0 ? (
                        <div className="rounded-2xl border border-dashed border-white/15 bg-slate-950/60 p-12 text-center">
                            <div className="mb-3 text-5xl">🔍</div>
                            <h3 className="text-xl font-semibold text-white">No products found</h3>
                            <p className="mt-2 text-slate-400">
                                {search ? `We couldn't find any products matching “${search}”.` : 'There are no products available right now.'}
                            </p>
                        </div>
                    ) : (
                        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 xl:grid-cols-4">
                            {products.map((product) => (
                                <ProductCard key={product.id} product={product} />
                            ))}
                        </div>
                    )}
                </div>
            </section>

            <section className="mx-auto max-w-7xl px-6 py-12">
                <div className="rounded-[2rem] border border-cyan-400/20 bg-gradient-to-r from-cyan-500/10 via-violet-500/10 to-emerald-500/10 p-8 text-center shadow-2xl shadow-cyan-950/30 md:p-12">
                    <p className="text-xs font-semibold uppercase tracking-[0.3em] text-cyan-300">Ready to shop</p>
                    <h2 className="mt-4 text-3xl font-black text-white md:text-4xl">Start your next favorite collection today.</h2>
                    <p className="mx-auto mt-3 max-w-2xl text-slate-300">
                        Join GenZone to explore exclusive deals, new arrivals, and premium essentials curated for your routine.
                    </p>
                    <div className="mt-8 flex justify-center gap-4">
                        <Link href="/signup" className="rounded-xl bg-gradient-to-r from-cyan-500 to-violet-500 px-6 py-3 font-semibold text-white shadow-lg shadow-cyan-500/20 transition hover:brightness-110">
                            Create account
                        </Link>
                        <Link href="/login" className="rounded-xl border border-white/10 bg-white/5 px-6 py-3 font-semibold text-white transition hover:border-cyan-400/40 hover:bg-cyan-500/10">
                            Log in
                        </Link>
                    </div>
                </div>
            </section>

            <footer className="border-t border-white/10 bg-slate-950/80">
                <div className="mx-auto flex max-w-7xl flex-col gap-4 px-6 py-8 md:flex-row md:items-center md:justify-between">
                    <div>
                        <p className="text-xl font-bold text-white">GenZone</p>
                        <p className="mt-1 text-sm text-slate-400">Your one-stop destination for everyday essentials.</p>
                    </div>

                    <div className="flex items-center gap-5 text-sm text-slate-300">
                        <Link href="/login" className="hover:text-white">Log in</Link>
                        <Link href="/signup" className="hover:text-white">Sign up</Link>
                        <Link href="/admin/login" className="hover:text-white">Admin</Link>
                    </div>
                </div>
            </footer>
        </main>
    )
}
