import React, { useEffect, useState } from "react";
import { url } from "../main";
import axios from "axios";
import { toast, ToastContainer } from "react-toastify";
import "./cart.css";
import arrow from "../assets/arrow-right.svg"
import src from "../assets/src.svg"
import trash from "../assets/trash.svg"
const Cart = () => {
    const [reservations, setReservations] = useState([]);
    const getCarts = async () => {
        const response = await axios.get(url + "/cart", {
            headers: { Authorization: sessionStorage.getItem("token") },
        });
        setReservations(response.data.Data.Reservations);
    };

    const deleteReservation = async (id) => {
        try {
            const response = await axios.delete(`${url}/reservation`, {
                data: { ReservationId: id },
                headers: { Authorization: sessionStorage.getItem("token") }
            });
    
            if (response.data.Error === "") {
                setReservations(reservations.filter(reservation => reservation.Id !== id));
            }
        } catch (error) {
            console.error('Error deleting reservation:', error);
        }
    };

    const buyTicket = async (id) => {
        const response = await axios.post(url + "/ticket", 
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
                        <img src={trash} width={"24px"} onClick={() => deleteReservation(res.Id)} className="delete" />
                        <br />
                    </div>
                ))}
            </div>
            
        </div>
    );
};

export default Cart;
