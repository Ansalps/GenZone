"use client"

import axios from 'axios'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { z } from 'zod'
import { toast } from 'sonner'

const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*])[A-Za-z\d!@#$%^&*]{8,}$/

const signupSchema = z
    .object({
        firstName: z.string().min(1, 'First name is required'),
        lastName: z.string(),
        email: z.string().email('Invalid email address'),
        password: z
            .string()
            .min(8, 'Password should have a minimum length of 8')
            .refine(
                (value) => passwordRegex.test(value),
                'Password should contain at least one uppercase letter, one lowercase letter, one digit, and one special character'
            ),
        confirmPassword: z.string(),
        phone: z.string().regex(/^\d{10}$/, 'Phone number must be 10 digits'),
    })
    .refine((data) => data.password === data.confirmPassword, {
        message: 'Passwords do not match',
        path: ['confirmPassword'],
    })

type FormData = z.infer<typeof signupSchema>

export default function SignUp() {
    const router = useRouter()

    const [formData, setFormData] = useState<FormData>({
        firstName: '',
        lastName: '',
        email: '',
        password: '',
        confirmPassword: '',
        phone: '',
    })

    const [errors, setErrors] = useState<Partial<Record<keyof FormData, string>>>({})
    const [showPassword, setShowPassword] = useState(false)
    const [showConfirmPassword, setShowConfirmPassword] = useState(false)

    function onChange(e: React.ChangeEvent<HTMLInputElement>) {
        const { name, value } = e.target
        setFormData((prev) => ({
            ...prev,
            [name]: value,
        }))
    }

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault()
        setErrors({})

        const validationResult = signupSchema.safeParse(formData)
        if (!validationResult.success) {
            const fieldErrors: Partial<Record<keyof FormData, string>> = {}

            validationResult.error.issues.forEach((issue) => {
                const path = issue.path[0] as keyof FormData
                if (path && !fieldErrors[path]) {
                    fieldErrors[path] = issue.message
                }
            })

            setErrors(fieldErrors)
            const firstError = validationResult.error.issues[0]?.message || 'Validation failed'
            toast.error(firstError)
            return
        }

        try {
            const response = await axios.post(
                'http://localhost:8080/public/signup',
                {
                    first_name: formData.firstName,
                    last_name: formData.lastName,
                    email: formData.email,
                    password: formData.password,
                    confirm_password: formData.confirmPassword,
                    phone: formData.phone,
                },
                {
                    withCredentials: true,
                }
            )

            toast.success(response.data?.message || 'Sign up successful')
            router.push(`/verify-otp?email=${encodeURIComponent(formData.email)}`)
        } catch (error) {
            toast.error('Something went wrong')
        }
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
                <div className="grid w-full max-w-5xl overflow-hidden rounded-[2rem] border border-white/10 bg-slate-900/70 shadow-2xl shadow-slate-950/50 backdrop-blur-xl lg:grid-cols-[0.9fr_1.1fr]">
                    <div className="relative hidden flex-col justify-between bg-gradient-to-br from-slate-800 via-slate-900 to-slate-950 p-8 lg:flex">
                        <div>
                            <span className="inline-flex rounded-full border border-violet-400/30 bg-violet-500/10 px-3 py-1 text-[10px] font-semibold uppercase tracking-[0.25em] text-violet-200">
                                New here
                            </span>
                            <h1 className="mt-6 max-w-sm text-4xl font-black tracking-tight text-white">
                                Create your account and discover the best deals.
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
                            <p className="text-xs font-semibold uppercase tracking-[0.3em] text-cyan-300">Create account</p>
                            <h2 className="mt-3 text-3xl font-bold text-white">Sign up</h2>
                            <p className="mt-2 text-sm text-slate-400">Join GenZone and unlock fresh arrivals, profile perks, and faster checkout.</p>
                        </div>

                        <form onSubmit={handleSubmit} className="space-y-5">
                            <div className="grid gap-5 sm:grid-cols-2">
                                <div className="space-y-2">
                                    <label htmlFor="firstName" className="text-sm font-medium text-slate-200">
                                        First name
                                    </label>
                                    <input
                                        type="text"
                                        id="firstName"
                                        name="firstName"
                                        value={formData.firstName}
                                        onChange={onChange}
                                        className="w-full rounded-xl border border-white/10 bg-slate-950/70 px-4 py-3 text-white outline-none transition focus:border-cyan-400/60"
                                        placeholder="John"
                                    />
                                    {errors.firstName && <span className="text-xs text-rose-400">{errors.firstName}</span>}
                                </div>

                                <div className="space-y-2">
                                    <label htmlFor="lastName" className="text-sm font-medium text-slate-200">
                                        Last name
                                    </label>
                                    <input
                                        type="text"
                                        id="lastName"
                                        name="lastName"
                                        value={formData.lastName}
                                        onChange={onChange}
                                        className="w-full rounded-xl border border-white/10 bg-slate-950/70 px-4 py-3 text-white outline-none transition focus:border-cyan-400/60"
                                        placeholder="Doe"
                                    />
                                    {errors.lastName && <span className="text-xs text-rose-400">{errors.lastName}</span>}
                                </div>
                            </div>

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
                                        placeholder="Create a strong password"
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

                            <div className="space-y-2">
                                <label htmlFor="confirmPassword" className="text-sm font-medium text-slate-200">
                                    Confirm password
                                </label>
                                <div className="flex items-center rounded-xl border border-white/10 bg-slate-950/70 px-3 transition focus-within:border-cyan-400/60">
                                    <input
                                        type={showConfirmPassword ? 'text' : 'password'}
                                        id="confirmPassword"
                                        name="confirmPassword"
                                        value={formData.confirmPassword}
                                        onChange={onChange}
                                        className="w-full bg-transparent px-2 py-3 text-white outline-none placeholder:text-slate-500"
                                        placeholder="Repeat password"
                                    />
                                    <button
                                        type="button"
                                        className="px-2 text-xs font-medium text-cyan-300 transition hover:text-cyan-200"
                                        onClick={() => setShowConfirmPassword((prev) => !prev)}
                                    >
                                        {showConfirmPassword ? 'Hide' : 'Show'}
                                    </button>
                                </div>
                                {errors.confirmPassword && <span className="text-xs text-rose-400">{errors.confirmPassword}</span>}
                            </div>

                            <div className="space-y-2">
                                <label htmlFor="phone" className="text-sm font-medium text-slate-200">
                                    Phone number
                                </label>
                                <input
                                    type="tel"
                                    id="phone"
                                    name="phone"
                                    value={formData.phone}
                                    onChange={onChange}
                                    className="w-full rounded-xl border border-white/10 bg-slate-950/70 px-4 py-3 text-white outline-none transition focus:border-cyan-400/60"
                                    placeholder="9876543210"
                                />
                                {errors.phone && <span className="text-xs text-rose-400">{errors.phone}</span>}
                            </div>

                            <button
                                type="submit"
                                className="w-full rounded-xl bg-gradient-to-r from-cyan-500 to-violet-500 px-4 py-3 text-sm font-semibold text-white shadow-lg shadow-cyan-500/20 transition hover:brightness-110"
                            >
                                Create account
                            </button>

                            <p className="text-center text-sm text-slate-400">
                                Already have an account?{' '}
                                <Link href="/login" className="font-medium text-cyan-300 hover:text-cyan-200">
                                    Sign in
                                </Link>
                            </p>
                        </form>
                    </div>
                </div>
            </section>
        </main>
    )
}
