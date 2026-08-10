'use client'
import axios from "axios";
import { useEffect,useState } from "react"
import { useParams,useRouter } from 'next/navigation'
import CategoryForm from "@/components/category-from";


export default function EditCategory(){
    const router = useRouter();
    const params = useParams();
    const id = params.id; // Extracts 'id' directly from the URL route
    const [isLoading,setIsLoading]=useState(false);
    const [formData,setFormData]=useState({
        categoryName:'',
        categoryDescription:'',
        categoryImageUrl:''
    })
    useEffect(()=>{
        async function fetchCategory(){
            try{
                const response= await axios.get(`http://localhost:8080/admin/category/${id}`,
                    {withCredentials:true},
                );

                const category = response.data.data;

                setFormData({
                categoryName: category.category_name,
                categoryDescription: category.category_description,
                categoryImageUrl: category.category_image_url,
                });
            } catch(error){
                console.log(error)
            }
        }
        fetchCategory();
    },[id])
    
    const handleChange=(e:React.ChangeEvent<HTMLInputElement>)=>{
            const {name,value}=e.target;
            setFormData((prevData)=>({
                ...prevData,
                [name]: value // Updates only the field being edited
            }))
        }
        const handleSubmit=async (e:React.FormEvent<HTMLFormElement>) =>{
        e.preventDefault();
            setIsLoading(true);
        // 3. formData is already a JS object containing all current values
        console.log('Form Data from State:', formData);

        try {
            const response = await axios.put(`http://localhost:8080/admin/category/${id}`, 
                {
                    'category_name':formData.categoryName,
                    'category_description':formData.categoryDescription,
                    'category_image_url':formData.categoryImageUrl
                },
                { withCredentials: true }
            );
            if (response.status==200){
                console.log('success')
                router.push('/admin/categories');
            }
            console.log(response.data);
        } catch (error) {
            
            console.error(error);
        }finally{
            setIsLoading(false)
        } 
    };
    
    return (
        <>
            <div className="font-bold text-2xl">
                Edit Category
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