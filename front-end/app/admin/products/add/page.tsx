'use client'
import ProductForm from "@/components/product-form"
import axios from "axios"
import { useRouter } from "next/navigation"
import { useState,useEffect } from "react"

interface Category {
  id: number;
  category_name: string;
}

export default function AddProduct(){
    const [categories, setCategories] = useState<Category[]>([]);

    useEffect(() => {
    const fetchCategories = async () => {
        try {
        const response = await axios.get(
            "http://localhost:8080/admin/category",{
                withCredentials:true
            }
        );
        
        setCategories(response.data.data.categories);
        } catch (error) {
        console.error(error);
        }
    };

    fetchCategories();
    }, []);

    const router=useRouter();
    const [isLoading,setIsLoading]=useState(false);
    //single state object holding form values
        const [formData,setFormData]=useState({
            categoryName:'',
            productName:'',
            productDescription:'',
            productImageUrl:'',
            price:0,
            stock:0,
            popular:false,
            hasOffer:false,
            size:''
        })

      const onChange = (
            e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
            ) => {
            const { name, value, type } = e.target;

            setFormData((prev) => ({
                ...prev,
                [name]:
                type === "checkbox"
                    ? (e.target as HTMLInputElement).checked
                    : type === "number"
                    ? Number(value)
                    : value,
            }));
        };

    const handleSubmit=async (e:React.FormEvent<HTMLFormElement>) =>{
        e.preventDefault();

        // 3. formData is already a JS object containing all current values
        console.log('Form Data from State:', formData);
        setIsLoading(true);
        try {
            console.log(`popular: ${formData.popular}`)
            const response = await axios.post('http://localhost:8080/admin/product', 
                {
                    'category_name':formData.categoryName,
                    'product_name':formData.productName,
                    'product_description':formData.productDescription,
                    'product_image_url':formData.productImageUrl,
                    'price':formData.price,
                    'stock':formData.stock,
                    'popular':formData.popular,
                    'size':formData.size,
                },
                { withCredentials: true }
            );
            if (response.status==200){
                console.log('success')
               router.push("/admin/products");
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
                    <ProductForm
                        formData={formData} 
                        onChange={onChange}
                        isLoading={isLoading}
                        categories={categories}
                    />
                </form>
            </div>
        </>
    )
}