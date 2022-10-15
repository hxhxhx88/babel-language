import React, { FC } from "react"
import { Alert } from "antd"
import { I18nText } from "component/Text"

const PageError: FC<{ message: React.ReactNode }> = ({ message }) => {
    return (
        <div style={{ width: "100vw", height: "100vh", display: "flex", justifyContent: "center", alignItems: "center", padding: 16 }}>
            <Alert message={<I18nText id="error" transform="capitalize" />} style={{ width: "100%" }} description={message} type="error" showIcon />
        </div>
    )
}

export default PageError
