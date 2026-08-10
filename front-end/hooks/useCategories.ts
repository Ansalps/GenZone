import { useState,useEffect } from "react";
import axios from "axios";

interface Category {
  id: number;
  category_name: string;
}

export function useCategories(){
    const [categories, setCategories] = useState<Category[]>([]);
    const [isCategoriesLoading, setIsCategoriesLoading] = useState<boolean>(true);
    const [categoriesError, setCategoriesError] = useState<unknown>(null);
    useEffect(() => {
    const fetchCategories = async () => {
        try {
        setIsCategoriesLoading(true);
        const response = await axios.get(
            "http://localhost:8080/admin/category",{
                withCredentials:true
            }
        );
        
        setCategories(response.data.data.categories);
        } catch (error) {
        console.error(error);
        setCategoriesError(error);
        }
    };

    fetchCategories();
    }, []);
    return { categories, isCategoriesLoading, categoriesError };
}