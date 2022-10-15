/** @jsxImportSource @emotion/react */

import { FC } from "react"
import { css } from "@emotion/react"

const Text: FC<{
    transform: "capitalize" | "uppercase" | "lowercase" | "capitalize-first"
    children: React.ReactNode
}> = ({ transform, children }) => {
    switch (transform) {
        case "capitalize-first":
            return (
                <div
                    css={css`
                        ::first-letter {
                            text-transform: capitalize;
                        }
                    `}
                >
                    {children}
                </div>
            )
        default:
            return <span style={{ textTransform: transform }}>{children}</span>
    }
}

export default Text
