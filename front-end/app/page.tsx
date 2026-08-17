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

    const [isLoadingCategories, setIsLoadingCategories] =
        useState(true)

    const [isLoadingProducts, setIsLoadingProducts] =
        useState(true)

    useEffect(() => {

        const fetchCategories = async () => {

            try {

                const response = await axios.get(
                    'http://localhost:8080/public/category'
                )

                if (
                    response.data?.status &&
                    response.data?.data?.categories
                ) {
                    setCategories(
                        response.data.data.categories
                    )
                }

            } catch (error) {

                console.error(
                    'Failed to load categories:',
                    error
                )

            } finally {

                setIsLoadingCategories(false)

            }
        }

        const fetchProducts = async () => {

            try {

                const response = await axios.get(
                    'http://localhost:8080/public/product'
                )

                if (
                    response.data?.status &&
                    response.data?.data?.products
                ) {
                    setProducts(
                        response.data.data.products
                    )
                }

            } catch (error) {

                console.error(
                    'Failed to load products:',
                    error
                )

            } finally {

                setIsLoadingProducts(false)

            }
        }

        fetchCategories()
        fetchProducts()

    }, [])

    return (
        <div className="min-h-screen bg-gray-100">

            {/* ================= HEADER ================= */}

            <header className="bg-green-800 text-white">

                <div className="max-w-7xl mx-auto px-6 py-5">

                    <div className="flex items-center justify-between">

                        {/* Logo / Brand */}

                        <Link
                            href="/"
                            className="text-2xl font-bold"
                        >
                            GeZOne
                        </Link>

                        {/* Auth Buttons */}

                        <div className="flex items-center gap-3">

                            <Link
                                href="/login"
                                className="border border-white px-5 py-2 rounded-lg hover:bg-white hover:text-green-800 transition"
                            >
                                Log In
                            </Link>

                            <Link
                                href="/signup"
                                className="bg-white text-green-800 px-5 py-2 rounded-lg font-semibold hover:bg-gray-100 transition"
                            >
                                Sign Up
                            </Link>

                        </div>

                    </div>

                </div>

            </header>


            {/* ================= HERO ================= */}

            <section className="bg-green-700 text-white">

                <div className="max-w-7xl mx-auto px-6 py-16">

                    <div className="max-w-2xl">

                        <h1 className="text-4xl md:text-5xl font-bold">
                            Discover Something You’ll Love
                        </h1>

                        <p className="mt-4 text-green-100 text-lg">
                            Explore our latest products and
                            find everything you need in one place.
                        </p>

                        <Link
                            href="/signup"
                            className="inline-block mt-8 bg-white text-green-800 px-6 py-3 rounded-lg font-semibold hover:bg-gray-100 transition"
                        >
                            Get Started
                        </Link>

                    </div>

                </div>

            </section>


            {/* ================= CATEGORIES ================= */}

            <section className="max-w-7xl mx-auto px-6 py-12">

                <div className="flex items-center justify-between mb-6">

                    <div>

                        <h2 className="text-2xl md:text-3xl font-bold text-gray-800">
                            Shop by Category
                        </h2>

                        <p className="text-gray-500 mt-1">
                            Explore our product categories
                        </p>

                    </div>

                </div>


                {isLoadingCategories ? (

                    <div className="text-center py-10 text-gray-500">
                        Loading categories...
                    </div>

                ) : categories.length === 0 ? (

                    <div className="text-center py-10 text-gray-500">
                        No categories available.
                    </div>

                ) : (

                    <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">

                        {categories.map((category) => (

                            <CategoryCard
                                key={category.id}
                                category={category}
                            />

                        ))}

                    </div>

                )}

            </section>


            {/* ================= PRODUCTS ================= */}

            <section className="bg-gray-50">

                <div className="max-w-7xl mx-auto px-6 py-12">

                    <div className="mb-6">

                        <h2 className="text-2xl md:text-3xl font-bold text-gray-800">
                            Featured Products
                        </h2>

                        <p className="text-gray-500 mt-1">
                            Check out our latest products
                        </p>

                    </div>


                    {isLoadingProducts ? (

                        <div className="text-center py-10 text-gray-500">
                            Loading products...
                        </div>

                    ) : products.length === 0 ? (

                        <div className="text-center py-10 text-gray-500">
                            No products available.
                        </div>

                    ) : (

                        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">

                            {products.map((product) => (

                                <ProductCard
                                    key={product.id}
                                    product={product}
                                />

                            ))}

                        </div>

                    )}

                </div>

            </section>


            {/* ================= FOOTER ================= */}

            <footer className="bg-gray-900 text-gray-300">

                <div className="max-w-7xl mx-auto px-6 py-8">

                    <div className="flex flex-col md:flex-row justify-between gap-4">

                        <div>

                            <h3 className="text-xl font-bold text-white">
                                GeZOne
                            </h3>

                            <p className="text-sm mt-2">
                                Your one-stop online shopping destination.
                            </p>

                        </div>

                        <div className="flex gap-4">

                            <Link
                                href="/login"
                                className="hover:text-white"
                            >
                                Log In
                            </Link>

                            <Link
                                href="/signup"
                                className="hover:text-white"
                            >
                                Sign Up
                            </Link>

                        </div>

                    </div>

                </div>

            </footer>

        </div>
    )
}