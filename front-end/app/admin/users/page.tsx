'use client'

import axios from "axios";
import { useEffect } from "react";

export default function Users(){
    useEffect(()=>{
         // 1. Declare the async function inside the effect
        const fetchData = async () => {
            try {
            const response = await axios.get('http://localhost:8080/admin/listusers', {
                withCredentials: true // Remeber to include this for cookies!
            });
            console.log(response.data);
            } catch (error) {
            console.error('Error fetching data:', error);
            }
        };

        // 2. Call the function immediately
        fetchData();
    },[])
    return (
        <>
            <div className="w-screen h-25 bg-green-800 text-amber-50">
                Users
            </div>
            <div>
                

            </div>
        </>
    )
}