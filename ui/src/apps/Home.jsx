import React, { useEffect, useState } from "react";
import MyMap from "../components/MyMap";
import "./home.css";

import Header from "../components/Header";
import world from "../assets/world.svg";
import ticket from "../assets/ticket.svg";
import cart from "../assets/cart.svg";
import src from "../assets/src.svg";
import dest from "../assets/dest.svg";

import Input from "../components/Input";
import axios from "axios";
import { url } from "../main";
import { useNavigate } from "react-router-dom";
import Tickets from "../components/Tickets";
import { ToastContainer, toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import Wishlist from "../components/Wishlist";
import Cart from "../components/Cart";

const Home = () => {
    const navigate = useNavigate();
    const notify = () =>
        toast.error("Não foram encontradas rotas para o caminho desejado.", {
            position: "top-center",
            autoClose: 5000,
            theme: "light",
        });

    const [selected, setSelected] = useState(3);
    const [tickets, setTickets] = useState([]);
    const [user, setUser] = useState({
        Name: "",
        Username: "",
    });
    const [path, setPath] = useState([]);

    const [source, setSource] = useState("");
    const [destination, setDestination] = useState("");

    const fetchUserData = async () => {
        try {
            const response = await axios.get(url + "/user", {
                headers: { Authorization: sessionStorage.getItem("token") },
            });

            if (response.data.Error !== "") {
                navigate("/");
            } else setUser(response.data.Data.user);
        } catch (error) {
            navigate("/");
        }
    };
    useEffect(() => {
        if (source && destination) {
            if (source === destination) {
                notify();
                setPath([]);
                return;
            }
            const getRoute = async () => {
                try {
                    const response = await axios.get(url + "/route", {
                        params: {
                            src: source,
                            dest: destination,
                        },
                        headers: {
                            Authorization: sessionStorage.getItem("token"),
                        },
                    });
                    if (response.data.Error !== "") {
                        notify();
                        setPath([]);
                        return;
                    }
                    setPath(response.data.Data.path);
                } catch (error) {
                    console.error("Error fetching the route:", error);
                }
            };
            getRoute();
        }
    }, [source, destination]);

    useEffect(() => {
        fetchUserData();
    }, []);

    const changeScreen = () => {
        switch (selected) {
            case 0:
                return <MyMap paths={path} />;
            case 1:
                return <Tickets tickets={tickets} />;
            case 2:
                return <Cart />
                
            case 3:
                return <div className="home-img" />;
        }
    };

    return (
        <div className="home">
            <Header name={user.Name} />
            <ToastContainer/>
            <main>
                <div className="grid">
                    <div className="row">
                        {selected == 0 && (
                            <>
                                <Input
                                    placeholder={"Origem"}
                                    img={src}
                                    value={source}
                                    setValue={setSource}
                                />
                                <Input
                                    placeholder={"Destino"}
                                    img={dest}
                                    value={destination}
                                    setValue={setDestination}
                                />
                            </>
                        )}
                    </div>
                    <div className="row-2">
                        <div className="col left-col">
                            <div
                                className={`option ${
                                    selected == 0 && "selected"
                                }`}
                                onClick={() => setSelected(0)}
                            >
                                <img src={world} alt="" width={"22px"} />
                                <h4>Visualizar Rotas</h4>
                            </div>
                            <div
                                className={`option ${
                                    selected == 1 && "selected"
                                }`}
                                onClick={() => setSelected(1)}
                            >
                                <img src={ticket} alt="" width={"22px"} />
                                <h4>Minhas passagens</h4>
                            </div>
                            <div
                                className={`option ${
                                    selected == 2 && "selected"
                                }`}
                                onClick={() => setSelected(2)}
                            >
                                <img src={cart} alt="" width={"22px"} />
                                <h4>Meu carrinho</h4>
                            </div>
                        </div>
                        <div className="col main-col">{changeScreen()}</div>
                        <div className="col right-col">
                            <Wishlist paths={path} setPaths={setPath}/>
                        </div>
                    </div>
                </div>
            </main>
            <footer></footer>
        </div>
    );
};

export default Home;
