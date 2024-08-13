import "@mantine/core/styles.css";
import React from "react";
import ReactDOM from "react-dom/client";
import Layout from "../components/layout";
import Sidebar from "../components/thread_sidebar";
import ChatWindowsComponent from "../components/thread";

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
