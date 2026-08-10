'use client'
import axios from "axios";
import { useEffect,useState } from "react"
import { useParams,useRouter } from 'next/navigation'
import ProductForm from "@/components/product-form";
import { useCategories } from "@/hooks/useCategories"


export default function EditProduct(){
     const { categories } = useCategories();

    const router = useRouter();
    const params = useParams();
    const id = params.id; // Extracts 'id' directly from the URL route
    const [isLoading,setIsLoading]=useState(false);
    const [formData,setFormData]=useState({
        categoryName:'',
        productName:'',
        productDescription:'',
        productImageUrl:'',
        price:0,
        stock:0,
        size:'',
        popular:false,
        hasOffer:false,
        discountPercentage:0,
        discountAmount:0,
        totalDiscountedAmount:0
    })
    useEffect(()=>{
        async function fetchCategory(){
            try{
                const response= await axios.get(`http://localhost:8080/admin/product/${id}`,
                    {withCredentials:true},
                );

                const product = response.data.data;

                setFormData({
                    categoryName:product.category_name,
                    productName:product.product_name,
                    productDescription:product.product_description,
                    productImageUrl:product.product_image_url,
                    price:product.price,
                    stock:product.stock,
                    size:product.size,
                    popular:product.popular,
                    hasOffer:product.has_offer,
                    discountPercentage:product.offer_discount_percent,
                    discountAmount:product.discount_amount,
                    totalDiscountedAmount:product.total_discounted_amount
                });
            } catch(error){
                console.log(error)
            }
        }
        fetchCategory();
    },[id])
    
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
            setIsLoading(true)
        // 3. formData is already a JS object containing all current values
        console.log('Form Data from State:', formData);

        try {
            const response = await axios.put(`http://localhost:8080/admin/product/${id}`, 
                {
                    'category_name':formData.categoryName,
                    'product_name':formData.productName,
                    'product_description':formData.productDescription,
                    'product_image_url':formData.productImageUrl,
                    'price':formData.price,
                    'stock':formData.stock,
                    'popular':formData.popular,
                    'size':formData.size,
                    'has_offer':formData.hasOffer,
                    'offer_discount_percent':formData.discountPercentage,
                    'discount_amount':formData.discountAmount,
                    'total_discounted_amount':formData.totalDiscountedAmount
                },
                { withCredentials: true }
            );
            if (response.status==200){
                console.log('success')
                router.push('/admin/products');
            }
            console.log(response.data);
        } catch (error) {
            
            console.error(error);
        } finally{
            setIsLoading(false);
        }
    };
    
    return (
        <>
            <div className="font-bold text-2xl">
                Edit Category
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