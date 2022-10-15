import { FC, useState, useEffect } from "react"
import { components } from "openapi/babel"

import { listCorpuses } from "api"

const List: FC = () => {
    const [corpuses, setCorpuses] = useState<readonly components["schemas"]["Corpus"][]>([])

    useEffect(() => {
        listCorpuses({}).then((resp) => {
            setCorpuses(resp.data.corpuses)
        }).catch(console.error)
    }, [])

    return (
        <div>
            {corpuses.map((c) => <p key={c.id}>{c.title}</p>)}
        </div>
    )
}

export default List
