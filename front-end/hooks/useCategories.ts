import { useEffect, useState } from "react";
import axios from "axios";

interface Category {
    id: number;
    category_name: string;
}

export function useCategories() {
    const [categories, setCategories] = useState<Category[]>([]);
    const [isCategoriesLoading, setIsCategoriesLoading] =
        useState(true);
    const [categoriesError, setCategoriesError] =
        useState<unknown>(null);

    useEffect(() => {
        const fetchCategories = async () => {
            try {
                setIsCategoriesLoading(true);
                setCategoriesError(null);

                const response = await axios.get(
                    "http://localhost:8080/admin/category",
                    {
                        withCredentials: true,
                    }
                );

                console.log("Category API response:", response.data);

                const fetchedCategories =
                    response.data?.data?.categories;

                if (Array.isArray(fetchedCategories)) {
                    setCategories(fetchedCategories);
                } else {
                    console.error(
                        "Unexpected category response:",
                        response.data
                    );

                    setCategories([]);
                }
            } catch (error) {
                console.error(
                    "Failed to fetch categories:",
                    error
                );

                setCategoriesError(error);
                setCategories([]);
            } finally {
                setIsCategoriesLoading(false);
            }
        };

        fetchCategories();
    }, []);

    return {
        categories,
        isCategoriesLoading,
        categoriesError,
    };
}