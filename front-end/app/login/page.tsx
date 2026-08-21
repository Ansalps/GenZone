'use client'
import axios from "axios";
import { useRouter } from "next/navigation";
import { useState,SubmitEvent } from "react"
import { toast } from "sonner";
import {z} from "zod"

const loginSchema=z.object({
    email:z.email("Invalid email address"),
    password:z.string().min(6,"password must be atleast 6 characters")
})
type  FormData= z.infer<typeof loginSchema>

export default function Login(){
    const router=useRouter()
    const [formData,setFormData]=useState<FormData>({
        email:'',
        password:''
    })
    const [errors,setErrors]=useState<Partial<Record<keyof FormData, string>>>({});
    const [showPassword, setShowPassword] = useState(false);

    async function handleSubmit(e:SubmitEvent<HTMLFormElement>){
        e.preventDefault()
        setErrors({})//clear previous errors
        const validationResult=loginSchema.safeParse(formData)
        if (!validationResult.success){
            const fieldErrors: Partial<Record<keyof FormData,string>> = {};

            validationResult.error.issues.forEach((issue) => {
                const path = issue.path[0] as keyof FormData;
                if (path && !fieldErrors[path]) {
                fieldErrors[path] = issue.message;
                }
            });

            setErrors(fieldErrors);

            toast.error(validationResult.error.issues[0].message);
            return; // Stop execution if validation fails
            
        }
        try{
            const response=await axios.post('http://localhost:8080/public/login',{
                'email':formData.email,
                'password':formData.password
            },
            {
                withCredentials: true
            })
            toast.success(
                response.data?.message || "Login successful"
            );
            router.push('/dashboard')
        } catch(error){
           console.error("Login error:", error);

            if (axios.isAxiosError(error)) {
            toast.error(
                    error.response?.data?.message ||
                    "Login failed"
                );
            } else {
                toast.error("Something went wrong");
            }
        }
        
    }
    function onChange(e:React.ChangeEvent<HTMLInputElement>){
        const {name,value}=e.target
        setFormData({...formData,
            [name]:value
        })
        
    }
    return (
        <div className="relative">
            <button
                type="button"
                onClick={() => router.push('/')}
                className="absolute top-5 left-5 border border-black px-3 py-2 rounded hover:bg-black hover:text-white"
            >
                Home
            </button>

            <form onSubmit={handleSubmit}>
                
                <div className="flex flex-col justify-center items-center h-screen gap-4">
                
                    <div className="flex flex-col">
                    <label htmlFor="email">Email</label>
                    <input 
                        type="text" 
                        id="email" 
                        name="email"
                        value={formData.email}
                        className="border-black border p-1"
                        onChange={onChange}
                    />
                </div>

                <div className="flex flex-col">
                    <label htmlFor="password">Password</label>
                    <div className="flex items-center border border-black">
                        <input 
                            type={showPassword ? "text" : "password"} 
                            id="password" 
                            className="flex-1 p-1 outline-none"
                            name="password"
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
                </div>

                    <div className="flex flex-col">
                        <button type="submit">Sumbit</button>
                    </div>

                </div>
            </form>
        </div>
        
    )
}
