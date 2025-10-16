import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import "./index.css";
import Loader from "./components/Loader";
import ItemsPage from "./components/ItemsPage";

const router = createBrowserRouter([
  { path: "/", element: <Loader /> },
  { path: "/items", element: <ItemsPage /> },
]);

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);
