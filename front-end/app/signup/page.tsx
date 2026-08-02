"use client"
import axios from 'axios';
export default function SignUp(){
    const handleSubmit = async (e:React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        // Handle form submission logic here
        const response=await axios.get('http://localhost:8080')
        
        console.log(response)
    };
    return (
        /* 1. Added flex-col to stack items vertically */
        /* 2. Added h-screen so the box stretches from top to bottom of the screen */
        <div className="flex flex-col justify-center items-center h-screen gap-4">
            <form onSubmit={handleSubmit} className="flex flex-col gap-4 w-full max-w-xs">
                <div className="flex flex-col">
                    <label htmlFor="firstName">First name *</label>
                    <input type="text" id="firstName" className="border-black border p-1" required/>
                </div>

                <div className="flex flex-col">
                    <label htmlFor="lastName">Last name</label>
                    <input type="text" id="lastName" className="border-black border p-1" required/>
                </div>

                <div className="flex flex-col">
                    <label htmlFor="email">Email *</label>
                    <input type="email" id="email" className="border-black border p-1" required/>
                </div>

                <div className="flex flex-col">
                    <label htmlFor="password">Password *</label>
                    <input type="password" id="password" className="border-black border p-1" required/>
                </div>

                <div className="flex flex-col">
                    <label htmlFor="confirmPassword">Confirm password *</label>
                    <input type="password" id="confirmPassword" className="border-black border p-1" required/>
                </div>

                <div className="flex flex-col">
                    <label htmlFor="phone">Phone *</label>
                    <input type="tel" id="phone" className="border-black border p-1" required/>
                </div>

                <button type="submit" className="bg-black text-white p-2 mt-2 rounded">
                        Sign Up
                </button>

            </form>
        </div>
    )
}
