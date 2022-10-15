import { FC } from "react"
import { Spin } from "antd"

const Detail: FC = () => {
    return (
        <div
            style={{
                width: "100vw",
                height: "100vh",
                display: "flex",
                justifyContent: "center",
                alignItems: "center",
            }}
        >
            <Spin />
        </div>
    )
}

export default Detail
