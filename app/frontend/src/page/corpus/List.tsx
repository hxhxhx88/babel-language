import { FC, useState, useEffect } from "react"
import { Button, Space } from "antd"
import { DefaultService, Corpus } from "openapi/babel"
import { Link } from "react-router-dom"

import PageSpin from "component/PageSpin"
import Layout from "Layout"
import routePath from "route"

const List: FC = () => {
    const [corpuses, setCorpuses] = useState<Corpus[]>([])

    const [isListCorpuses, setIsListCorpuses] = useState<boolean>(false)
    useEffect(() => {
        setIsListCorpuses(true)
        DefaultService.listCorpuses()
            .then(({ corpuses }) => setCorpuses(corpuses))
            .catch(console.error)
            .finally(() => setIsListCorpuses(false))
    }, [])

    return (
        <Layout>
            <Space direction="vertical" style={{ width: "100%" }}>
                {corpuses.map((c) => (
                    <Button key={c.id} style={{ width: "100%", textAlign: "left" }} size="large">
                        <Link to={routePath(`/corpus/${c.id}`)}>{c.title}</Link>
                    </Button>
                ))}
            </Space>
            {isListCorpuses && <PageSpin />}
        </Layout>
    )
}

export default List
