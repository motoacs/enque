import React from "react";
import { createRoot } from "react-dom/client";
import "./lib/i18n";
import "./style.css";
import App from "./App";
import { ErrorBoundary } from "./components/ErrorBoundary";
import { registerEventListeners } from "./lib/events";

const container = document.getElementById("root");
const root = createRoot(container!);

// Register Wails event listeners
registerEventListeners();

root.render(
  <React.StrictMode>
    <ErrorBoundary>
      <App />
    </ErrorBoundary>
  </React.StrictMode>
);
