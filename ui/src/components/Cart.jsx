import React, { useEffect, useState } from "react";
import { url } from "../main";
import axios from "axios";
import { toast, ToastContainer } from "react-toastify";
import "./cart.css";
import arrow from "../assets/arrow-right.svg"
import src from "../assets/src.svg"
const Cart = () => {
    const [reservations, setReservations] = useState([]);
    const getCarts = async () => {
        const response = await axios.get(url + "/cart", {
            headers: { Authorization: sessionStorage.getItem("token") },
        });
        setReservations(response.data.Data.Reservations);
    };

    const buyTicket = async (id) => {
        const response = await axios.post(url + "/buy", 
            {ReservationId : id}, {
            headers: { Authorization: sessionStorage.getItem("token") },
        });
        if (response.data.Error === "")
            setReservations(reservations.filter(reservation => reservation.Id !== id))

    }

    useEffect(() => {
        getCarts();
    }, []);

    return (
        <div className="reservations">
            <div className="reservations-list">
                {reservations.map((res, i) => (
                    <div key={i} className="reservation">
                        <img src={src} alt="" />
                        
                        <h4>
                        {res.Src.Name}
                        </h4>
                         <img src={arrow} width={"24px"} />{" "}
                         <h4>
                        {res.Dest.Name}
                        </h4>
                        <button className="buy" onClick={() => buyTicket(res.Id)}>Comprar</button>
                        <br />
                    </div>
                ))}
            </div>
            
        </div>
    );
};

export default Cart;
