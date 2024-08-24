import { Center, Container } from "@mantine/core";
import "@mantine/core/styles.css";
import React from "react";
import ReactDOM from "react-dom/client";
import Layout from "../components/layout";

function Main() {
  return (
    <Container>
      <Center>
        <h1>Hello Vibrain</h1>
      </Center>
    </Container>
  );
}

function App() {
  return <Layout main={<Main />} />;
}

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
