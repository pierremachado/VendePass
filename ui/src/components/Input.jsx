import React, { useEffect } from "react";
import "./input.css";
import capitals from "../../brazilcapitals.json";
const Input = ({ placeholder, img, value, setValue }) => {

    let capitals_names = Object.keys(capitals)
    capitals_names.sort()
    
    useEffect(() => setValue(""), [])

    return (
        <div className="input">
            <img src={img} alt="" height={"20px"} />
            <select
                className="route-input"
                value={value}
                onChange={(e) => setValue(e.target.value)}
            >
                <option value="" disabled>
                    {placeholder}
                </option>
                {capitals_names.map((capital, i) => (
                    <option key={i} value={capital}>
                        {capital}
                    </option>
                ))}
            </select>
        </div>
    );
};

export default Input;
