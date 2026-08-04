

    interface CategoryFormProps {
        formData: {
            categoryName: string;
            categoryDescription: string;
            categoryImageUrl:string;
           
        };
        onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
        isLoading:boolean;
    }

    export default function CategoryForm({ formData, onChange,isLoading }: CategoryFormProps){
            
        return (
            <>
                <div className="flex flex-col">
                    <label htmlFor="categoryName">Category name *</label>
                    <input 
                        type="text" 
                        id="categoryName" 
                        name="categoryName"
                        value={formData.categoryName}
                        onChange={onChange}
                        className="border-black border p-1" 
                        required/>
                </div>


                <div className="flex flex-col">
                    <label htmlFor="categoryDescription">Category description *</label>
                    <input 
                        type="text" 
                        id="categoryDescription"
                        name="categoryDescription"
                        value={formData.categoryDescription}
                        onChange={onChange} 
                        className="border-black border p-1" 
                        required/>
                </div>

                <div className="flex flex-col">
                    <label htmlFor="categoryImageUrl">Category image url *</label>
                    <input 
                        type="text" 
                        id="categoryImageUrl"
                        name="categoryImageUrl"
                        value={formData.categoryImageUrl}
                        onChange={onChange} 
                        className="border-black border p-1" 
                        required/>
                </div>

               

                <button disabled={isLoading} type="submit" className="bg-black text-white p-2 mt-2 rounded">
                        {isLoading ? "Saving..." : "Submit"}
                </button>
                
            </>
        )
    }