import React, { useEffect } from "react";
import { url } from "../main";
import axios from "axios";
import { toast, ToastContainer } from "react-toastify";

const Cart = () => {
    const getCarts = async () => {
        const response = await axios.get(url+"/cart", {
            headers: { Authorization: sessionStorage.getItem("token") },
        });
        console.log(response.data)
    };

    useEffect(() => {
        getCarts();
    }, [])

   

    return <div>
        
    </div>;
};

export default Cart;
