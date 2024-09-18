import React, { useState } from "react";
import {
    MapContainer,
    TileLayer,
    Marker,
    Popup,
    Polyline,
} from "react-leaflet";
import "leaflet/dist/leaflet.css";
import capitals from "../../brazilcapitals.json";

const MyMap = ({ paths }) => {
    const position = [capitals.Brasília.latitude, capitals.Brasília.longitude]; // Latitude e longitude para o centro do mapa
    const [selectedCity, setSelectedCity] = useState("");

    function selectCity(city) {
        setSelectedCity(city);
    }
    const capitalsArray = Object.entries(capitals).map(([capital, data]) => ({
        name: capital,
        ...data,
    }));

    return (
        <>
            <MapContainer
                minZoom={3.5}
                center={position}
                zoom={3.5}
                style={{ height: "100%", width: "100%" }}
            >
                <TileLayer
                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                    attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                />
                {capitalsArray.map((capital) => (
                    <Marker
                        key={capital.name}
                        position={[capital.latitude, capital.longitude]}
                        eventHandlers={{
                            click: () => selectCity(capital.name),
                        }}
                    >
                        <Popup key={capital.name}>
                            {capital.name}, {capital.state}
                        </Popup>
                    </Marker>
                ))}
                {paths && (
                    <>
                        {paths.map((line, i) => (
                            <Polyline
                                key={i}
                                positions={[
                                    [
                                        line.Path[0].Latitude,
                                        line.Path[0].Longitude,
                                    ],
                                    [
                                        line.Path[1].Latitude,
                                        line.Path[1].Longitude,
                                    ],
                                ]}
                                color="blue"
                            />
                        ))}
                    </>
                )}
            </MapContainer>
        </>
    );
};

export default MyMap;
