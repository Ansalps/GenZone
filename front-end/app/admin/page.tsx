'use client'
import axios from "axios"
import Link from 'next/link';

export default function AdminDashboard(){
    async function hadleLogout() {
        try {
            // CRITICAL FIX: Added { withCredentials: true } as the 3rd argument
            const response = await axios.post(
                'http://localhost:8080/admin/logout', 
                {}, // Empty request body (2nd argument)
                { withCredentials: true } // Config object (3rd argument)
            );
            
            console.log(response.data);
            
            // Redirect cleanly back to the login page after the cookie is wiped
            window.location.href = '/admin/login';
            
        } catch (error) {
            console.error("Logout failed:", error);
        }
    }
    return (
        <div className="text-2xl 
      font-semibold flex justify-start gap-4 items-center">
            <Link 
                href="/admin/users" // Change to your actual route path
                className="bg-blue-500 text-amber-50 cursor-pointer p-2 rounded text-base"
            >
                users
            </Link>
            <Link 
                href="/admin/categories" // Change to your actual route path
                className="bg-blue-500 text-amber-50 cursor-pointer p-2 rounded text-base"
            >
                Categories
            </Link>
            <button className="bg-blue-500 text-amber-50 cursor-pointer">products</button>
            <button className="bg-blue-500 text-amber-50 cursor-pointer">orders</button>
            <button className="bg-blue-500 text-amber-50 cursor-pointer">coupons</button>
            <button className="bg-blue-500 text-amber-50 cursor-pointer">product offers</button>
            <button className="bg-blue-500 text-amber-50 cursor-pointer">sales report</button>
            <button className="bg-blue-500 text-amber-50 cursor-pointer">best selling</button>
            <button className="bg-blue-500 text-amber-50 cursor-pointer">invoice</button>
            <button className="bg-blue-500 text-amber-50 cursor-pointer" onClick={hadleLogout}>Log Out</button>
        </div>
    )
}