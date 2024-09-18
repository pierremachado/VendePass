import axios from 'axios'
import React, { useEffect, useState } from 'react'
import { url } from '../main';
import src from "../assets/src.svg"
import arrow from "../assets/arrow-right.svg"
import "./ticket.css"
import trash from "../assets/trash.svg"
const Tickets = () => {
  const [tickets, setTickets] = useState([]);

  const getTickets = async () => {
    const response = await axios.get(url + "/tickets", {
        headers: { Authorization: sessionStorage.getItem("token") },
    });
    setTickets(response.data.Data.Tickets);
    console.log(response.data.Data.Tickets);
  };

  const deleteTickets = async (id) => {
    console.log(id)
    try {
        const response = await axios.delete(`${url}/ticket`, {
            data: { TicketId: id },
            headers: { Authorization: sessionStorage.getItem("token") }
        });

        if (response.data.Error === "") {
            setTickets(tickets.filter(ticket => ticket.Id !== id));
        }
    } catch (error) {
        console.error('Error deleting ticket:', error);
    }
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
                <img src={trash} width={"24px"} onClick={() => deleteTickets(ticket.Id)} className="delete" />
            </div>
        ))}
    </div>
    
</div>
  )
}

export default Tickets