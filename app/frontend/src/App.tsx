import React, { useState } from "react"
import { Button } from "antd"
import { Fetcher } from "openapi-typescript-fetch"

import { paths, ApiPaths } from "openapi/babel"
const fetcher = Fetcher.for<paths>()
fetcher.configure({
    baseUrl: "/api",
})

function createCorpus() {
    fetcher.path(ApiPaths.CreateCorpus).method("post").create()({
        corpus: {
            title: "hello",
            original_language_iso_639_3: "eng",
            translations: [{
                language_iso_639_3: "eng",
                blocks: [{
                    content: "hello",
                    rank: 1,
                    uuid: "sent.1"
                }, {
                    content: "world",
                    rank: 1,
                    uuid: "sent.2"
                }]
            }, {
                language_iso_639_3: "zho",
                blocks: [{
                    content: "你好",
                    rank: 1,
                    uuid: "sent.1"
                }, {
                    content: "世界",
                    rank: 1,
                    uuid: "sent.2"
                }]
            }]
        }
    }).then((resp) => {
        alert(resp.data.id)
    }).catch(console.error)
}

function App (): JSX.Element {
    const [text, setText] = useState("Fetching...")

    fetcher.path(ApiPaths.GetMetadata).method("get").create()({}).then((resp) => {
        setText(resp.data.commit_identifier ?? "")
    }).catch(console.error)

    return <Button type="primary" onClick={createCorpus}>{text}</Button>
}

export default App
