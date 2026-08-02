'use client'
import axios from "axios";
import { useParams } from "next/navigation";
import { useEffect } from "react"

export default function DeleteCategory(){
    const params = useParams();
    const id = params.id; // Extracts 'id' directly from the URL route
        
    useEffect(()=>{
        async function fetchCategory(){
            try{
                const response= await axios.delete(`http://localhost:8080/admin/category/${id}`,
                    {withCredentials:true},
                );
                window.location.href='/admin/categories'
                console.log(response)
            } catch(error){
                console.log(error)
            }
        }
        fetchCategory();
    },[id])
    return (
        <>
        </>
    )
}