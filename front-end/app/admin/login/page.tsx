'use client'
import axios from 'axios';
import { useState } from 'react';

const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

export default function Login() {
    const [formData, setFormData] = useState({
        email: '',
        password: '',
    });
    const [showPassword, setShowPassword] = useState(false);
    const [submitError, setSubmitError] = useState('');

    const [errors, setErrors] = useState<{ email?: string; password?: string }>({});

    const validateForm = () => {
        const nextErrors: { email?: string; password?: string } = {};

        if (!formData.email.trim()) {
            nextErrors.email = 'Email is required.';
        } else if (!emailRegex.test(formData.email)) {
            nextErrors.email = 'Enter a valid email address.';
        }

        if (!formData.password.trim()) {
            nextErrors.password = 'Password is required.';
        } else if (formData.password.length < 8) {
            nextErrors.password = 'Password must be at least 8 characters.';
        }

        setErrors(nextErrors);
        return Object.keys(nextErrors).length === 0;
    };

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;

        setFormData((prevData) => ({
            ...prevData,
            [name]: value,
        }));

        setErrors((prev) => ({
            ...prev,
            [name]: undefined,
        }));

        if (submitError) {
            setSubmitError('');
        }
    };

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setSubmitError('');

        if (!validateForm()) {
            return;
        }

        try {
            const response = await axios.post('http://localhost:8080/admin/login', formData, {
                withCredentials: true,
            });

            if (response.status === 200) {
                window.location.href = '/admin';
            }
        } catch (error: unknown) {
            const axiosError = error as { response?: { data?: { message?: string; error?: string } } };
            const backendMessage =
                axiosError.response?.data?.message ||
                axiosError.response?.data?.error ||
                'Unable to login. Please try again.';

            setSubmitError(backendMessage);
            console.error(error);
        }
    };

    return (
        <div className="flex min-h-screen items-center justify-center bg-slate-950 px-4 py-10">
            <div className="absolute inset-0 overflow-hidden">
                <div className="absolute -left-20 top-16 h-72 w-72 rounded-full bg-cyan-500/20 blur-3xl" />
                <div className="absolute bottom-10 right-10 h-80 w-80 rounded-full bg-violet-500/20 blur-3xl" />
            </div>

            <div className="relative w-full max-w-md rounded-3xl border border-white/10 bg-white/5 p-8 shadow-2xl shadow-cyan-950/40 backdrop-blur-xl">
                <div className="mb-8 text-center">
                    <div className="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-2xl bg-gradient-to-br from-cyan-400 to-violet-500 text-2xl font-bold text-white shadow-lg shadow-cyan-500/30">
                        A
                    </div>
                    <p className="text-xs font-semibold uppercase tracking-[0.3em] text-cyan-300">
                        Admin Portal
                    </p>
                    <h1 className="mt-3 text-3xl font-bold text-white">Welcome back</h1>
                    <p className="mt-2 text-sm text-slate-300">
                        Sign in to manage your store and operations.
                    </p>
                </div>

                <form onSubmit={handleSubmit} className="space-y-5">
                    {submitError && (
                        <div className="rounded-xl border border-rose-500/40 bg-rose-500/10 px-3 py-2 text-sm text-rose-300">
                            {submitError}
                        </div>
                    )}

                    <div className="space-y-2">
                        <label htmlFor="email" className="text-sm font-medium text-slate-200">
                            Email address
                        </label>
                        <input
                            type="email"
                            id="email"
                            name="email"
                            value={formData.email}
                            onChange={handleChange}
                            className={`w-full rounded-xl border bg-slate-900/80 px-4 py-3 text-white outline-none transition focus:ring-2 ${
                                errors.email
                                    ? 'border-rose-500 focus:border-rose-400 focus:ring-rose-400/30'
                                    : 'border-slate-700 focus:border-cyan-400 focus:ring-cyan-400/30'
                            }`}
                            placeholder="admin@example.com"
                        />
                        {errors.email && <p className="text-xs text-rose-400">{errors.email}</p>}
                    </div>

                    <div className="space-y-2">
                        <label htmlFor="password" className="text-sm font-medium text-slate-200">
                            Password
                        </label>
                        <div className="relative">
                            <input
                                type={showPassword ? 'text' : 'password'}
                                id="password"
                                name="password"
                                value={formData.password}
                                onChange={handleChange}
                                className={`w-full rounded-xl border bg-slate-900/80 px-4 py-3 pr-12 text-white outline-none transition focus:ring-2 ${
                                    errors.password
                                        ? 'border-rose-500 focus:border-rose-400 focus:ring-rose-400/30'
                                        : 'border-slate-700 focus:border-cyan-400 focus:ring-cyan-400/30'
                                }`}
                                placeholder="Enter your password"
                            />
                            <button
                                type="button"
                                onClick={() => setShowPassword((prev) => !prev)}
                                className="absolute inset-y-0 right-3 flex items-center text-xs font-medium text-cyan-300 transition hover:text-cyan-200"
                                aria-label={showPassword ? 'Hide password' : 'Show password'}
                            >
                                {showPassword ? 'Hide' : 'Show'}
                            </button>
                        </div>
                        {errors.password && <p className="text-xs text-rose-400">{errors.password}</p>}
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
    );
}
