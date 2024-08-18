import "@mantine/core/styles.css";
import React from "react";
import ReactDOM from "react-dom/client";
import Layout from "../components/layout";
import ChatWindowsComponent from "../components/thread";
import Sidebar from "../components/thread_sidebar";

function App() {
  return (
    <Layout
      main={<ChatWindowsComponent />}
      nav={<Sidebar />}
      showNavBurger={true}
    />
  );
}

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
