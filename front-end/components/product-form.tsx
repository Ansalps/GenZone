
    interface Category {
        id: number;
        category_name: string;
    }
    interface ProductFormProps {
        formData: {
            categoryName: string;
            productName: string;
            productDescription: string;
            productImageUrl:string;
            price:number;
            stock:number;
            popular:boolean;
            hasOffer:boolean;
            size:string;
        };
        onChange: (e: React.ChangeEvent<HTMLInputElement| HTMLSelectElement >) => void;
        isLoading:boolean;
        categories:Category[]
    }

    export default function ProductForm({ formData, onChange,isLoading,categories }: ProductFormProps){
        
        return (
            <>
                <div className="flex flex-col">
                    <label htmlFor="categoryName">Category name *</label>

                    <select
                        id="categoryName"
                        name="categoryName"
                        value={formData.categoryName}
                        onChange={onChange}
                        className="bg-white text-black border rounded px-3 py-2"
                    >
                        <option value="">Select Category</option>

                        {categories.map((category) => (
                            <option
                                key={category.id}
                                value={category.category_name}
                            >
                                {category.category_name}
                            </option>
                        ))}
                    </select>
                </div>

                <div className="flex flex-col">
                    <label htmlFor="productName">Product name *</label>
                    <input 
                        type="text" 
                        id="productName" 
                        name="productName"
                        value={formData.productName}
                        onChange={onChange}
                        className="border-black border p-1" 
                        required/>
                </div>

                <div className="flex flex-col">
                    <label htmlFor="productDescription">Product description *</label>
                    <input 
                        type="text" 
                        id="productDescription"
                        name="productDescription"
                        value={formData.productDescription}
                        onChange={onChange} 
                        className="border-black border p-1" 
                        required/>
                </div>

                <div className="flex flex-col">
                    <label htmlFor="productImageUrl">Product image url *</label>
                    <input 
                        type="text" 
                        id="productImageUrl"
                        name="productImageUrl"
                        value={formData.productImageUrl}
                        onChange={onChange} 
                        className="border-black border p-1" 
                        required/>
                </div>

                <div className="flex flex-col">
                    <label htmlFor="price">Price *</label>
                    <input 
                        type="number" 
                        id="price"
                        name="price"
                        value={formData.price}
                        onChange={onChange} 
                        className="border-black border p-1" 
                        required/>
                </div>

                <div className="flex flex-col">
                    <label htmlFor="stock">Stock *</label>
                    <input 
                        type="number" 
                        id="stock"
                        name="stock"
                        value={formData.stock}
                        onChange={onChange} 
                        className="border-black border p-1" 
                        required/>
                </div>

                <div className="flex items-center gap-2">
                    <input
                        type="checkbox"
                        id="popular"
                        name="popular"
                        checked={formData.popular}
                        onChange={onChange}
                    />
                    <label htmlFor="popular" className="cursor-pointer">
                        Popular
                    </label>
                </div>

                <div className="flex items-center gap-2">
                    <input
                        type="checkbox"
                        id="hasOffer"
                        name="hasOffer"
                        checked={formData.hasOffer}
                        onChange={onChange}
                    />
                    <label htmlFor="popular" className="cursor-pointer">
                        Has Offer
                    </label>
                </div>

                <div className="flex flex-col">
                    <label htmlFor="size">Size *</label>

                    <select
                        id="size"
                        name="size"
                        value={formData.size}
                        onChange={onChange}
                        className="bg-white text-black border rounded px-3 py-2"
                    >
                        <option value="">Select size</option>
                        <option value="Small">Small</option>
                        <option value="Medium">Medium</option>
                        <option value="Large">Large</option>
                    </select>
                </div>

                <button disabled={isLoading} type="submit" className="bg-black text-white p-2 mt-2 rounded">
                        {isLoading ? "Saving..." : "Submit"}
                </button>
                
            </>
        )
    }