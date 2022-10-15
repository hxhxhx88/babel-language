import { FC, useState, useEffect } from "react"
import { DefaultService, Corpus } from "openapi/babel"

import Layout from "Layout"

const List: FC = () => {
    const [corpuses, setCorpuses] = useState<Corpus[]>([])

    useEffect(() => {
        DefaultService.listCorpuses().then((resp) => {
            setCorpuses(resp.corpuses)
        }).catch(console.error)
    }, [])

    return (
        <Layout>
            {corpuses.map((c) => <p key={c.id}>{c.title}</p>)}
        </Layout>
    )
}

export default List
