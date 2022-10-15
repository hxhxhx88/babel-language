/** @jsxImportSource @emotion/react */

import { FC } from "react"
import { css } from "@emotion/react"
import { FormattedMessage } from "react-intl"

const cssCapitalizeFirst = css`
    ::first-letter {
        text-transform: capitalize;
    }
`

type TransformName = "capitalize" | "uppercase" | "lowercase" | "capitalize-first"

export const Text: FC<{
    transform: TransformName
    children: React.ReactNode
}> = ({ transform, children }) => {
    switch (transform) {
        case "capitalize-first":
            return <div css={cssCapitalizeFirst}>{children}</div>
        default:
            return <span style={{ textTransform: transform }}>{children}</span>
    }
}

export const I18nText: FC<{ transform: TransformName; id: string }> = ({ transform, id }) => {
    return (
        <Text transform={transform}>
            <FormattedMessage id={id} />
        </Text>
    )
}
