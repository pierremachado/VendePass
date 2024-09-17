import axios from "axios";
import React, { useEffect, useState } from "react";
import { url } from "../main";

const Wishlist = ({ paths }) => {
    const [flights, setFlights] = useState([]);

    useEffect(() => {

        const findFlights = async () => {
            if (!paths || !paths.length) return;
        
            const flightIds = paths.map((route) => route.FlightId);
        
            try {
                const response = await axios.post(
                    `${url}/flights`,
                    { FlightIds: flightIds }, 
                    {
                        headers: {
                            Authorization: `${sessionStorage.getItem("token")}` 
                        }
                    }
                );
        
                console.log(response.data.Data.Flights);
                setFlights(response.data.Data.Flights);
            } catch (error) {
                console.error('Error fetching flights:', error);
            }
        };

        findFlights();
    }, [paths]);

    return (
        <div className="wishlist">
            {flights && flights.map((flight, i) => (
                <div className="wish" key={i}>
                    {flight.Seats}
                    {flight.Src}
                    {flight.Dest}
                </div>
            ))}
        </div>
    );
};

export default Wishlist;
