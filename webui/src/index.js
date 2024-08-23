import React from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import App from "./App";

// Find the root element in the DOM
const rootElement = document.getElementById("root");

// Create a root and render the app
const root = createRoot(rootElement);
root.render(<App />);
