import React from "react";
import "./input.css";
import capitals from "../../brazilcapitals.json";
const Input = ({ placeholder, img, value, setValue }) => {
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
                {Object.keys(capitals).map((capital, i) => (
                    <option key={i} value={capital}>
                        {capital}
                    </option>
                ))}
            </select>
        </div>
    );
};

export default Input;
