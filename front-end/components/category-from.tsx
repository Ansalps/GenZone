interface CategoryFormProps {
    formData: {
        categoryName: string;
        categoryDescription: string;
        categoryImage: File | null;
    };
    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
    isLoading: boolean;
}

export default function CategoryForm({ formData, onChange, isLoading }: CategoryFormProps) {
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
                    required 
                />
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
                    required 
                />
            </div>

            {/* Replaced categoryImageUrl string input with categoryImage file input */}
            <div className="flex flex-col">
                <label htmlFor="categoryImage">Category image *</label>
                <input 
                    type="file" 
                    id="categoryImage"
                    name="categoryImage"
                    accept="image/*"
                    onChange={onChange} 
                    className="border-black border p-1 cursor-pointer" 
                    required 
                />
            </div>

            <button 
                disabled={isLoading} 
                type="submit" 
                className="bg-black text-white p-2 mt-2 rounded disabled:bg-gray-400 cursor-pointer"
            >
                {isLoading ? "Saving..." : "Submit"}
            </button>
        </>
    )
}