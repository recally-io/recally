import "@mantine/core/styles.css";
import React from "react";
import ReactDOM from "react-dom/client";
import Assistants from "../components/Assistants";
import Layout from "../components/Layout";

function App() {
  return <Layout main={<Assistants />} />;
}

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
