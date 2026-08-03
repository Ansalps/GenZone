'use client'
import CategoryForm from "@/components/category-from"
import axios from "axios"
import { useRouter } from "next/navigation"
import { useState } from "react"

export default function AddCategory(){
    const router=useRouter();
    const [isLoading,setIsLoading]=useState(false);
    //single state object holding form values
        const [formData,setFormData]=useState({
            categoryName:'',
            categoryDescription:'',
            categoryImageUrl:''
        })

        const handleChange=(e:React.ChangeEvent<HTMLInputElement>)=>{
            const {name,value}=e.target;
            setFormData((prevData)=>({
                ...prevData,
                [name]: value // Updates only the field being edited
            }))
        }

    const handleSubmit=async (e:React.FormEvent<HTMLFormElement>) =>{
        e.preventDefault();

        // 3. formData is already a JS object containing all current values
        console.log('Form Data from State:', formData);
        setIsLoading(true);
        try {
            const response = await axios.post('http://localhost:8080/admin/category', 
                {
                    'category_name':formData.categoryName,
                    'category_description':formData.categoryDescription,
                    'category_image_url':formData.categoryImageUrl
                },
                { withCredentials: true }
            );
            if (response.status==200){
                console.log('success')
               router.push("/admin/categories");
            }
            console.log(response.data);
        } catch (error) {
            
            console.error(error);
        } finally{
            setIsLoading(false)
        } 
    };
    return (
        <>
            <div className="font-bold text-2xl">
                Add Category
            </div>
            <div className="flex flex-col justify-center items-center h-screen gap-4">
                <form onSubmit={handleSubmit} className="flex flex-col gap-4 w-full max-w-xs">
                    <CategoryForm
                        formData={formData} 
                        onChange={handleChange}
                        isLoading={isLoading}
                    />
                </form>
            </div>
        </>
    )
}