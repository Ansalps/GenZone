'use client'

import axios from 'axios'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useState, type SubmitEvent } from 'react'
import { toast } from 'sonner'
import { z } from 'zod'

const loginSchema = z.object({
    email: z.email('Invalid email address'),
    password: z.string().min(6, 'Password must be at least 6 characters'),
})

type FormData = z.infer<typeof loginSchema>

export default function Login() {
    const router = useRouter()
    const [formData, setFormData] = useState<FormData>({
        email: '',
        password: '',
    })
    const [errors, setErrors] = useState<Partial<Record<keyof FormData, string>>>({})
    const [showPassword, setShowPassword] = useState(false)

    async function handleSubmit(e: SubmitEvent<HTMLFormElement>) {
        e.preventDefault()
        setErrors({})

        const validationResult = loginSchema.safeParse(formData)
        if (!validationResult.success) {
            const fieldErrors: Partial<Record<keyof FormData, string>> = {}

            validationResult.error.issues.forEach((issue) => {
                const path = issue.path[0] as keyof FormData
                if (path && !fieldErrors[path]) {
                    fieldErrors[path] = issue.message
                }
            })

            setErrors(fieldErrors)
            toast.error(validationResult.error.issues[0].message)
            return
        }

        try {
            const response = await axios.post(
                'http://localhost:8080/public/login',
                {
                    email: formData.email,
                    password: formData.password,
                },
                {
                    withCredentials: true,
                }
            )

            toast.success(response.data?.message || 'Login successful')
            router.push('/dashboard')
        } catch (error) {
            console.error('Login error:', error)

            if (axios.isAxiosError(error)) {
                toast.error(error.response?.data?.message || 'Login failed')
            } else {
                toast.error('Something went wrong')
            }
        }
    }

    function onChange(e: React.ChangeEvent<HTMLInputElement>) {
        const { name, value } = e.target
        setFormData((prev) => ({
            ...prev,
            [name]: value,
        }))
    }

    return (
        <main className="relative min-h-screen overflow-hidden bg-slate-950 text-white">
            <div className="absolute inset-0 -z-10 overflow-hidden">
                <div className="absolute -left-16 top-20 h-72 w-72 rounded-full bg-cyan-500/20 blur-3xl" />
                <div className="absolute right-0 top-40 h-80 w-80 rounded-full bg-violet-500/20 blur-3xl" />
                <div className="absolute bottom-10 left-1/3 h-64 w-64 rounded-full bg-emerald-500/10 blur-3xl" />
            </div>

            <header className="mx-auto flex max-w-7xl items-center justify-between px-6 py-5">
                <Link href="/" className="flex items-center gap-3">
                    <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-cyan-400 to-violet-500 font-black text-slate-950 shadow-lg shadow-cyan-500/30">
                        G
                    </div>
                    <div>
                        <p className="text-lg font-bold tracking-tight">GenZone</p>
                        <p className="text-[10px] uppercase tracking-[0.25em] text-slate-400">Store</p>
                    </div>
                </Link>

                <Link
                    href="/"
                    className="rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-cyan-400/40 hover:bg-cyan-500/10"
                >
                    Home
                </Link>
            </header>

            <section className="mx-auto flex min-h-[calc(100vh-88px)] max-w-6xl items-center justify-center px-6 py-10">
                <div className="grid w-full max-w-5xl overflow-hidden rounded-[2rem] border border-white/10 bg-slate-900/70 shadow-2xl shadow-slate-950/50 backdrop-blur-xl lg:grid-cols-[1.1fr_0.9fr]">
                    <div className="relative hidden flex-col justify-between bg-gradient-to-br from-slate-800 via-slate-900 to-slate-950 p-8 lg:flex">
                        <div>
                            <span className="inline-flex rounded-full border border-cyan-400/30 bg-cyan-500/10 px-3 py-1 text-[10px] font-semibold uppercase tracking-[0.25em] text-cyan-200">
                                Welcome back
                            </span>
                            <h1 className="mt-6 max-w-sm text-4xl font-black tracking-tight text-white">
                                Your style, your routine, your account.
                            </h1>
                        </div>

                        <div className="overflow-hidden rounded-[1.5rem] border border-white/10 bg-slate-950/60 p-4">
                            <img
                                src="https://images.unsplash.com/photo-1521572267360-ee0c2909d518?auto=format&fit=crop&w=900&q=80"
                                alt="Fashion collection"
                                className="h-60 w-full rounded-[1.1rem] object-cover"
                            />
                        </div>
                    </div>

                    <div className="p-6 sm:p-8 lg:p-10">
                        <div className="mb-6">
                            <p className="text-xs font-semibold uppercase tracking-[0.3em] text-cyan-300">Member login</p>
                            <h2 className="mt-3 text-3xl font-bold text-white">Sign in</h2>
                            <p className="mt-2 text-sm text-slate-400">Access your saved items, orders, and account details.</p>
                        </div>

                        <form onSubmit={handleSubmit} className="space-y-5">
                            <div className="space-y-2">
                                <label htmlFor="email" className="text-sm font-medium text-slate-200">
                                    Email address
                                </label>
                                <input
                                    type="email"
                                    id="email"
                                    name="email"
                                    value={formData.email}
                                    onChange={onChange}
                                    className="w-full rounded-xl border border-white/10 bg-slate-950/70 px-4 py-3 text-white outline-none transition focus:border-cyan-400/60"
                                    placeholder="you@example.com"
                                />
                                {errors.email && <span className="text-xs text-rose-400">{errors.email}</span>}
                            </div>

                            <div className="space-y-2">
                                <label htmlFor="password" className="text-sm font-medium text-slate-200">
                                    Password
                                </label>
                                <div className="flex items-center rounded-xl border border-white/10 bg-slate-950/70 px-3 transition focus-within:border-cyan-400/60">
                                    <input
                                        type={showPassword ? 'text' : 'password'}
                                        id="password"
                                        name="password"
                                        value={formData.password}
                                        onChange={onChange}
                                        className="w-full bg-transparent px-2 py-3 text-white outline-none placeholder:text-slate-500"
                                        placeholder="Enter your password"
                                    />
                                    <button
                                        type="button"
                                        className="px-2 text-xs font-medium text-cyan-300 transition hover:text-cyan-200"
                                        onClick={() => setShowPassword((prev) => !prev)}
                                    >
                                        {showPassword ? 'Hide' : 'Show'}
                                    </button>
                                </div>
                                {errors.password && <span className="text-xs text-rose-400">{errors.password}</span>}
                            </div>

                            <div className="flex items-center justify-between text-sm">
                                <label className="flex items-center gap-2 text-slate-300">
                                    <input type="checkbox" className="h-4 w-4 rounded border-white/10 bg-slate-950/70" />
                                    Remember me
                                </label>
                                <Link href="/signup" className="text-cyan-300 hover:text-cyan-200">
                                    Create account
                                </Link>
                            </div>

                            <button
                                type="submit"
                                className="w-full rounded-xl bg-gradient-to-r from-cyan-500 to-violet-500 px-4 py-3 text-sm font-semibold text-white shadow-lg shadow-cyan-500/20 transition hover:brightness-110"
                            >
                                Sign in
                            </button>
                        </form>
                    </div>
                </div>
            </section>
        </main>
    )
}

