import axios from "axios";
import React, { useEffect, useState } from "react";
import { url } from "../main";
import "./wishlist.css";
import arrow from "../assets/arrow-right.svg";
import cart from "../assets/cart.svg";
const Wishlist = ({ paths }) => {
    const [flights, setFlights] = useState([]);
    const [status, setStatus] = useState("");
    const handleReservation = async () => {
        try {
            const response = await axios.post(
                url + "/reservation",
                {
                    data: paths.map((route) => route.FlightId),
                },
                {
                    headers: {
                        Authorization: sessionStorage.getItem("token"),
                    },
                }
            );

            if (response.data.Error != "") {
                return;
            }
        } catch (error) {}
    };

    useEffect(() => {
        const findFlights = async () => {
            if (!paths) return;

            const flightIds = paths.map((route) => route.FlightId);
            try {
                const response = await axios.post(
                    `${url}/flights`,
                    { FlightIds: flightIds },
                    {
                        headers: {
                            Authorization: `${sessionStorage.getItem("token")}`,
                        },
                    }
                );

                setFlights(response.data.Data.Flights);
            } catch (error) {
                console.error("Error fetching flights:", error);
            }
        };
        findFlights();
    }, [paths]);

    return (
        <div className="wishlist">
            <h2>Passagens</h2>
            <hr />
            <div className="wishes">
                {flights &&
                    flights.map((flight, i) => (
                        <div
                            className={`wish full ${flight.Seats <= 0 && "full"}`}
                            key={i}
                        >
                            <h4>
                                {flight.Src} <img src={arrow} width={"24px"} />{" "}
                                {flight.Dest}
                            </h4>
                            <p>Vagas: {flight.Seats}</p>
                        </div>
                    ))}
            </div>
            <button className="add-to-cart" onClick={handleReservation}>
                <img src={cart} height={"20px"} alt="" />
                Adicionar ao Carrinho
            </button>
        </div>
    );
};

export default Wishlist;
