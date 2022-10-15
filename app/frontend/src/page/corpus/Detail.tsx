import { FC, useState, useEffect } from "react"
import { useParams } from "react-router-dom"
import { DefaultService, Corpus, Translation } from "openapi/babel"
import { Breadcrumb, Select } from "antd"
import { Link } from "react-router-dom"
import { FormattedMessage } from "react-intl"

import routePath from "route"
import Layout from "Layout"
import PageSpin from "component/PageSpin"

const { Option } = Select

export interface Props {
    corpus: Corpus
}

const CorpusDetail: FC<Props> = (props) => {
    const { corpus } = props

    const [translations, setTranslations] = useState<Translation[]>([])
    useEffect(() => {
        DefaultService.listCorpusTranslations(corpus.id).then(
            ({ translations }) => {
                setTranslations(translations)
            }
        )
    }, [])

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

            <Select
                mode="multiple"
                allowClear={true}
                showArrow={true}
                style={{ width: "100%", marginTop: 16 }}
                placeholder={
                    <FormattedMessage id="select_corpus_translation_placeholder" />
                }
            >
                {translations.map((t) => (
                    <Option key={t.id}>{t.title}</Option>
                ))}
            </Select>
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
