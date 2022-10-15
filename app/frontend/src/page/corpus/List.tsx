import { FC, useState, useEffect } from "react"
import { Button } from "antd"
import { DefaultService, Corpus } from "openapi/babel"
import { Link } from "react-router-dom"

import Layout from "Layout"
import routePath from "route"

const List: FC = () => {
    const [corpuses, setCorpuses] = useState<Corpus[]>([])

    useEffect(() => {
        DefaultService.listCorpuses()
            .then(({ corpuses }) => {
                setCorpuses(corpuses)
            })
            .catch(console.error)
    }, [])

    return (
        <Layout>
            {corpuses.map((c) => (
                <Button key={c.id} style={{ width: "100%" }}>
                    <Link to={routePath(`/corpus/${c.id}`)}>{c.title}</Link>
                </Button>
            ))}
        </Layout>
    )
}

export default List
