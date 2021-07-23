import React, {useState} from 'react';

import Navbar from 'react-bootstrap/Navbar';
import Container from 'react-bootstrap/Container';
import Col from "react-bootstrap/Col";
import Row from "react-bootstrap/Row"

import './App.css';

const ws = new WebSocket('ws://' + window.location.host + '/api/polling/subscribe', 'echo-protocol');

const App = () => {
    const [message, setMessage] = useState("hi");
    ws.onopen = function () {
        console.log('socket connection opened properly');
        ws.send("Hello World"); // send a message
        setMessage('message sent');
    };
    ws.onmessage = function (evt) {
        console.log("Message received = " + evt.data);
        setMessage(evt.data);
    };
    ws.onclose = function () {
        console.log("Connection closed...");
        setMessage('closed');
    };
    return (
        <>
            <Navbar bg="light" expand="lg">
                <Container>
                    <Navbar.Brand href="#home">Green Mountain Grill</Navbar.Brand>
                </Container>
            </Navbar>
            <Container>
                <Row>
                    <Col>Current Temp</Col>
                    <Col>Target Temp</Col>
                </Row>
                <p>
                    {message}
                </p>
            </Container>
        </>
    );
};

export default App;