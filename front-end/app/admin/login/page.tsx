'use client'
import axios from 'axios';
import { useState } from 'react';
export default function Login(){

    //single state object holding form values
    const [formData,setFormData]=useState({
        email:'',
        password:'',
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

        try {
            const response = await axios.post('http://localhost:8080/admin/login', 
                formData,
                { withCredentials: true }
            );
            if (response.status==200){
                console.log('success')
                window.location.href='/admin'
            }
            console.log(response.data);
        } catch (error) {
            
            console.error(error);
        }
    };
    return (
        /* 1. Added flex-col to stack items vertically */
        /* 2. Added h-screen so the box stretches from top to bottom of the screen */
        <div className="flex flex-col justify-center items-center h-screen gap-4">
            <form action="" onSubmit={handleSubmit}>
                    <div className="flex flex-col">
                    <label htmlFor="email">Email</label>
                    <input 
                        type="text" 
                        id="email"
                        name='email' 
                        value={formData.email}
                        onChange={handleChange}
                        className="border-black border p-1"
                        required
                    />
                    </div>

                    <div className="flex flex-col">
                        <label htmlFor="password">Password</label>
                        <input 
                            type="text" 
                            id="password" 
                            name='password'
                            value={formData.password}
                            onChange={handleChange}
                            className="border-black border p-1"
                            required
                        />
                    </div>
                    <button type="submit"
                        className='bg-blue-500 text-amber-50 cursor-pointer'
                    >
                        Submit
                    </button>
            </form>
        </div>
    )
}
