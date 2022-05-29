import React, {useState} from 'react';
import {
    MDBBtn,
    MDBBtnGroup,
    MDBCard,
    MDBCardBody,
    MDBCardText,
    MDBCardTitle,
    MDBCol,
    MDBCollapse,
    MDBContainer,
    MDBRow,
    MDBTypography,
} from 'mdb-react-ui-kit';

import './App.css';

const ws = new WebSocket('ws://' + window.location.host + '/api/polling/subscribe', 'echo-protocol');

const App = () => {
    const [currentState, setCurrentState] = useState({
        ID: "",
        CreatedAt: "",
        UpdatedAt: "",
        DeletedAt: "",
        SessionUUID: "",
        CurrentTemperature: 0,
        TargetTemperature: 0,
        Probe1Temperature: 0,
        Probe1TargetTemperature: 0,
        Probe2Temperature: 0,
        Probe2TargetTemperature: 0,
        WarnCode: "",
        PowerState: "",
        FireState: ""
    })
    const [sessionHist, setSessionHist] = useState([]);

    ws.onopen = function () {
        console.log('socket connection opened properly');
    };
    ws.onmessage = function (evt) {
        console.log("Message received = " + evt.data);
        const parsed = JSON.parse(evt.data);
        setSessionHist(parsed);
        if (parsed.length > 0) {
            setCurrentState(parsed[parsed.length - 1])
        }
    };
    ws.onclose = function () {
        console.log("Connection closed...");
    };
    const [showShow, setShowShow] = useState(false);

    const toggleShow = () => setShowShow(!showShow);

    return (
        <>
            <header>
                <div className='p-5 text-center bg-light'>
                    <h1 className='mb-3'>Green Mountain Grill</h1>
                </div>
            </header>
            <MDBContainer className='text-center'>
                <MDBTypography tag='strong'>Current: </MDBTypography>
                {currentState.CurrentTemperature}
                <br/>
                <MDBTypography tag='strong'>Set Temp: </MDBTypography>
                {currentState.TargetTemperature}
                <br/>
                <MDBTypography tag='strong'>Probe 1 Current: </MDBTypography>
                {currentState.Probe1Temperature}
                <br/>
                <MDBTypography tag='strong'>Probe 1 Set Temp: </MDBTypography>
                {currentState.Probe1TargetTemperature}
                <br/>
                <MDBTypography tag='strong'>Probe 2 Current: </MDBTypography>
                {currentState.Probe2Temperature}
                <br/>
                <MDBTypography tag='strong'>Probe 2 Set Temp: </MDBTypography>
                {currentState.Probe2TargetTemperature}
            </MDBContainer>
        </>
    );
};

export default App;