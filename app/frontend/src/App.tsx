import React, { useState } from "react"
import { Button } from "antd"
import { Fetcher } from "openapi-typescript-fetch"

import { paths, ApiPaths } from "openapi/babel"
const fetcher = Fetcher.for<paths>()
fetcher.configure({
    baseUrl: "/api",
})

function App (): JSX.Element {
    const [text, setText] = useState("Fetching...")

    fetcher.path(ApiPaths.GetMetadata).method("get").create()({}).then((resp) => {
        setText(resp.data.commit_identifier ?? "")
    }).catch(console.error)

    return <Button type="primary">{text}</Button>
}

export default App
