import React from 'react'

const Tickets = ({tickets}) => {
  return (
    <div className='tickets'>
        {tickets.map((ticket, i) => (
           <span className="ticket-row" key={i}>{ticket}</span> 
        ))}
    </div>
  )
}

export default Tickets