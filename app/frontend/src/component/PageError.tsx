import React, { FC } from "react"
import { Alert } from "antd"
import { FormattedMessage } from "react-intl"

const PageError: FC<{ message: React.ReactNode }> = ({ message }) => {
    const m = (
        <span style={{ textTransform: "capitalize" }}>
            <FormattedMessage id="error" />
        </span>
    )
    return (
        <div style={{ width: "100vw", height: "100vh", display: "flex", justifyContent: "center", alignItems: "center", padding: 16 }}>
            <Alert message={m} style={{ width: "100%" }} description={message} type="error" showIcon />
        </div>
    )
}

export default PageError
