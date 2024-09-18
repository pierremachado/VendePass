import axios from 'axios'
import React, { useEffect, useState } from 'react'
import { url } from '../main';
import src from "../assets/src.svg"
import arrow from "../assets/arrow-right.svg"
import "./ticket.css"

const Tickets = () => {
  const [tickets, setTickets] = useState([]);

  const getTickets = async () => {
    const response = await axios.get(url + "/tickets", {
        headers: { Authorization: sessionStorage.getItem("token") },
    });
    setTickets(response.data.Data.Tickets);
    console.log(response.data.Data.Tickets);
  };

  useEffect(() => {
    getTickets();
  }, []);

  return ( <div className="reservations">
    <div className="reservations-list">
        {tickets.map((ticket, i) => (
            <div key={i} className="ticket">
                <img src={src} alt="" />
                
                <h4>
                {ticket.Src.Name}
                </h4>
                 <img src={arrow} width={"24px"} />{" "}
                 <h4>
                {ticket.Dest.Name}
                </h4>
                
            </div>
        ))}
    </div>
    
</div>
    // <div className='tickets'>
    //     {tickets.map((ticket, i) => (
    //        <span className="ticket-row" key={i}>{ticket}</span> 
    //     ))}
    // </div>
  )
}

export default Tickets