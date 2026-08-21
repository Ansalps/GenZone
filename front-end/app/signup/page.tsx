"use client"
import axios from 'axios';
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { email, z } from "zod";
import { toast } from "sonner";

const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*])[A-Za-z\d!@#$%^&*]{8,}$/;

const signupSchema = z.object({
    firstName: z.string().min(1, "First name is required"),
    lastName: z.string(),
    email: z.string().email("Invalid email address"),
    password: z
        .string()
        .min(8, "Password shouls have a minimum length of 8")
        .refine(
            (value) => passwordRegex.test(value),
            "Password should contain at least one uppercase letter, one lowercase letter, one digit, and one special character"
        ),
    confirmPassword: z.string(),
    phone: z.string().regex(/^\d{10}$/, "Phone number must be 10 digits"),
}).refine((data) => data.password === data.confirmPassword, {
    message: "Passwords do not match",
    path: ["confirmPassword"],
});

type FormData = z.infer<typeof signupSchema>;

export default function SignUp() {
    const router = useRouter();

    const [formData, setFormData] = useState<FormData>({
        firstName: '',
        lastName: '',
        email: '',
        password: '',
        confirmPassword: '',
        phone: ''
    });

    const [errors, setErrors] = useState<Partial<Record<keyof FormData, string>>>({});
    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);

    function onChange(e: React.ChangeEvent<HTMLInputElement>) {
        const { name, value } = e.target;
        setFormData((prev) => ({
            ...prev,
            [name]: value
        }));
    }

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setErrors({});

        const validationResult = signupSchema.safeParse(formData);
        if (!validationResult.success) {
            const fieldErrors: Partial<Record<keyof FormData, string>> = {};

            validationResult.error.issues.forEach((issue) => {
                const path = issue.path[0] as keyof FormData;
                if (path && !fieldErrors[path]) {
                    fieldErrors[path] = issue.message;
                }
            });

            setErrors(fieldErrors);
            const firstError = validationResult.error.issues[0]?.message || "Validation failed";
            toast.error(firstError);
            return;
        }

        try {
            const response = await axios.post('http://localhost:8080/public/signup', {
                first_name: formData.firstName,
                last_name: formData.lastName,
                email: formData.email,
                password: formData.password,
                confirm_password: formData.confirmPassword,
                phone: formData.phone
            }, {
                withCredentials: true
            });

            toast.success(response.data?.message || "Sign up successful");
            router.push(`/verify-otp?email=${encodeURIComponent(formData.email)}`);
        } catch (error) {
            toast.error("Something went wrong");
        }
    };

    return (
        <div className="relative flex flex-col justify-center items-center min-h-screen gap-4">
            <button
                type="button"
                onClick={() => router.push('/')}
                className="absolute top-5 left-5 border border-black px-3 py-2 rounded hover:bg-black hover:text-white"
            >
                Home
            </button>

            <form onSubmit={handleSubmit} className="flex flex-col gap-4 w-full max-w-xs">
                
                <div className="flex flex-col">
                    <label htmlFor="firstName">First name *</label>
                    <input 
                        type="text" 
                        id="firstName" 
                        name='firstName'
                        className="border-black border p-1" 
                        value={formData.firstName}
                        onChange={onChange}
                    />
                    {errors.firstName && <span className="text-red-500 text-sm">{errors.firstName}</span>}
                </div>

                <div className="flex flex-col">
                    <label htmlFor="lastName">Last name</label>
                    <input 
                        type="text" 
                        id="lastName" 
                        name='lastName'
                        className="border-black border p-1" 
                        value={formData.lastName}
                        onChange={onChange}
                    />
                    {errors.lastName && <span className="text-red-500 text-sm">{errors.lastName}</span>}
                </div>

                <div className="flex flex-col">
                    <label htmlFor="email">Email *</label>
                    <input 
                        type="email" 
                        id="email" 
                        name="email"
                        className="border-black border p-1" 
                        value={formData.email}
                        onChange={onChange}
                    />
                    {errors.email && <span className="text-red-500 text-sm">{errors.email}</span>}
                </div>

                <div className="flex flex-col">
                    <label htmlFor="password">Password *</label>
                    <div className="flex items-center border border-black">
                        <input 
                            type={showPassword ? "text" : "password"} 
                            id="password" 
                            name="password"
                            className="flex-1 p-1 outline-none" 
                            value={formData.password}
                            onChange={onChange}
                        />
                        <button
                            type="button"
                            className="px-2 text-sm text-gray-700"
                            onClick={() => setShowPassword((prev) => !prev)}
                        >
                            {showPassword ? "Hide" : "Show"}
                        </button>
                    </div>
                    {errors.password && <span className="text-red-500 text-sm">{errors.password}</span>}
                </div>

                <div className="flex flex-col">
                    <label htmlFor="confirmPassword">Confirm password *</label>
                    <div className="flex items-center border border-black">
                        <input 
                            type={showConfirmPassword ? "text" : "password"} 
                            id="confirmPassword" 
                            name="confirmPassword"
                            className="flex-1 p-1 outline-none" 
                            value={formData.confirmPassword}
                            onChange={onChange}
                        />
                        <button
                            type="button"
                            className="px-2 text-sm text-gray-700"
                            onClick={() => setShowConfirmPassword((prev) => !prev)}
                        >
                            {showConfirmPassword ? "Hide" : "Show"}
                        </button>
                    </div>
                    {errors.confirmPassword && <span className="text-red-500 text-sm">{errors.confirmPassword}</span>}
                </div>

                <div className="flex flex-col">
                    <label htmlFor="phone">Phone *</label>
                    <input 
                        type="tel" 
                        id="phone" 
                        name="phone"
                        className="border-black border p-1" 
                        value={formData.phone}
                        onChange={onChange}
                    />
                    {errors.phone && <span className="text-red-500 text-sm">{errors.phone}</span>}
                </div>

                <button type="submit" className="bg-black text-white p-2 mt-2 rounded">
                    Sign Up
                </button>
            </form>
        </div>
    );
}