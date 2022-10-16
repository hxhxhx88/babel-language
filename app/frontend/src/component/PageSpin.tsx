import { FC } from "react"
import { Spin } from "antd"

const PageSpin: FC = () => {
    return (
        <div style={{ width: "100vw", height: "100vh", display: "flex", justifyContent: "center", alignItems: "center", position: "fixed", left: 0, top: 0 }}>
            <Spin size="large" />
        </div>
    )
}

export default PageSpin
