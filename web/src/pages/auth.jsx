import "@mantine/core/styles.css";
import React from "react";
import ReactDOM from "react-dom/client";
import Layout from "../components/layout";
import { AuthenticationForm } from "../components/auth";

function App() {
  return <Layout main={<AuthenticationForm />} />;
}

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
