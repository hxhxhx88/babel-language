import { FC, useState, useEffect } from "react"
import { useParams } from "react-router-dom"
import { DefaultService, Corpus } from "openapi/babel"
import { Breadcrumb } from "antd"
import { Link } from "react-router-dom"
import { FormattedMessage } from "react-intl"

import routePath from "route"
import Layout from "Layout"
import PageSpin from "component/PageSpin"

export interface Props {
    corpus: Corpus
}

const CorpusDetail: FC<Props> = (props) => {
    const { corpus } = props

    return (
        <Layout>
            <Breadcrumb>
                <Breadcrumb.Item>
                    <Link to={routePath(`/corpuses`)}>
                        <FormattedMessage id="corpus_list" />
                    </Link>
                </Breadcrumb.Item>
                <Breadcrumb.Item>{corpus.title}</Breadcrumb.Item>
            </Breadcrumb>
        </Layout>
    )
}

const Detail: FC = () => {
    const { corpusId = "" } = useParams()

    const [corpus, setCorpus] = useState<Corpus | undefined>(undefined)
    useEffect(() => {
        DefaultService.getCorpus(corpusId).then(({ corpus }) => {
            setCorpus(corpus)
        })
    }, [])

    return corpus ? <CorpusDetail corpus={corpus} /> : <PageSpin />
}

export default Detail
