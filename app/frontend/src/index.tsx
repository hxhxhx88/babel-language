import React from "react"
import ReactDOM from "react-dom/client"
import { IntlProvider } from "react-intl"
import App from "./App"

import "./index.css"
import "antd/dist/antd.min.css"

import(`locale/${process.env.REACT_APP_LANG}.json`).then((lang) => {
    const root = ReactDOM.createRoot(document.getElementById("root") as HTMLElement)
    root.render(
        <React.StrictMode>
            <IntlProvider locale="en-GB" messages={lang}>
                <App />
            </IntlProvider>
        </React.StrictMode>
    )
})
