import "./shims"

import ReactDOM from "react-dom"
import { App } from "./App"
import React from "react"

const node = document.createElement("div")
node.id = "root"
document.body.appendChild(node)

const root = ReactDOM.createRoot(node)
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
)

console.debug("react version: %s", React.version)
console.debug("react-dom version: %s", ReactDOM.version)
