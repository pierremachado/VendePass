import React from "react";
import logo from "../assets/logo.svg";
import profile from "../assets/profile.svg";
import logout from "../assets/logout.svg";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import { url } from "../main";

const Header = ({name}) => {
    const navigate = useNavigate();

    const handleLogout = async () => {
        const response = await axios.get(url + "/logout", {
            headers: { Authorization: sessionStorage.getItem("token") },
        });
        sessionStorage.setItem("token", "")
        console.log(response.data)
        navigate("/");
    };

    return (
        <header>
            <img src={logo} alt="logo" className="logo" width={"170px"} />
            <div className="line"></div>
            <div className="title">
                <h2>Hora de decolar.</h2>
                <h3>Qual é seu próximo destino?</h3>
            </div>
            <div className="line"></div>
            <div className="profile">
                <img src={profile} alt="" />
                <h4>{name}</h4>
                <img src={logout} alt="" className="logout-img" onClick={handleLogout} />
            </div>
        </header>
    );
};

export default Header;
