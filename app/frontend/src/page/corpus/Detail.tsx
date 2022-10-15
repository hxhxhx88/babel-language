import { FC, useState, useEffect, useMemo, HTMLAttributes } from "react"
import { useParams } from "react-router-dom"
import { DefaultService, Corpus, Translation } from "openapi/babel"
import { Switch, Space, Drawer, Alert, Breadcrumb, Button, PageHeader, DrawerProps, List, Select, Typography, InputNumber, Popover, Spin } from "antd"
import { SettingOutlined, LeftOutlined, RightOutlined } from "@ant-design/icons"
import { Link } from "react-router-dom"
import { Set as ImmutableSet, Map as ImmutableMap } from "immutable"

import { PAGE_SIZE } from "constant"
import routePath from "route"
import Layout from "Layout"
import PageSpin from "component/PageSpin"
import { I18nText } from "component/Text"

const { Option } = Select

export interface Props {
    corpus: Corpus
    translations: Translation[]
}

const CorpusDetail: FC<Props> = (props) => {
    const { corpus, translations } = props

    const translationLookup = useMemo(() => {
        return ImmutableMap(translations.map((t) => [t.id, t]))
    }, [translations])

    const [isSelectFormVisible, setIsSelectFormVisible] = useState<boolean | undefined>(undefined)
    const [selected, setSelected] = useState<Translation["id"][]>([])
    const [reference, setReference] = useState<Translation["id"] | undefined>(undefined)
    const [isCountTranslationBlocks, setIsCountTranslationBlocks] = useState<boolean>(false)
    const [totalCount, setTotalCount] = useState<number | undefined>(undefined)
    const [page, setPage] = useState<number>(0)
    useEffect(() => {
        if (!reference) return

        setIsCountTranslationBlocks(true)
        DefaultService.countTranslationBlocks(reference)
            .then(({ total_count }) => {
                setPage(0)
                setTotalCount(total_count)
            })
            .catch(console.error)
            .finally(() => setIsCountTranslationBlocks(false))
    }, [reference])

    return (
        <Layout>
            <Spin spinning={isCountTranslationBlocks}>
                <PageHeader
                    ghost={false}
                    breadcrumb={
                        <Breadcrumb>
                            <Breadcrumb.Item>
                                <Link to={routePath(`/corpuses`)}>
                                    <I18nText id="corpus_list" transform="capitalize" />
                                </Link>
                            </Breadcrumb.Item>
                            <Breadcrumb.Item>{corpus.title}</Breadcrumb.Item>
                        </Breadcrumb>
                    }
                    extra={[<Button icon={<SettingOutlined />} key="setting" onClick={() => setIsSelectFormVisible(true)} />]}
                    title={corpus.title}
                >
                    {selected.length > 0 ? (
                        <>
                            <I18nText id="selected_versions" transform="capitalize" />
                            <ul style={{ marginBottom: 0 }}>
                                {selected.map((tid) => (
                                    <li key={tid}>
                                        <Typography.Text strong={tid === reference}>{translationLookup.get(tid)?.title}</Typography.Text>
                                    </li>
                                ))}
                            </ul>
                        </>
                    ) : (
                        <Alert type="warning" showIcon message={<I18nText id="select_version_prompt" transform="capitalize-first" />} />
                    )}
                </PageHeader>
            </Spin>

            {(totalCount ?? 0) > 0 && (
                <div style={{ marginTop: 16, display: "flex", flexDirection: "row", justifyContent: "flex-end" }}>
                    <Pagination min={0} max={Math.ceil((totalCount ?? 0) / PAGE_SIZE) - 1} onConfirm={setPage} />
                </div>
            )}

            {isSelectFormVisible !== undefined && (
                <SelectTranslationDrawer
                    title={<I18nText id="select_versions" transform="capitalize" />}
                    placement="bottom"
                    closable={false}
                    open={isSelectFormVisible}
                    translations={translations}
                    reference={reference}
                    selected={selected}
                    afterOpenChange={(open) => !open && setIsSelectFormVisible(undefined)}
                    onCancel={() => setIsSelectFormVisible(false)}
                    onConfirm={(ids, reference) => {
                        setSelected(ids)
                        setReference(reference)
                        setIsSelectFormVisible(false)
                    }}
                />
            )}
        </Layout>
    )
}

const Detail: FC = () => {
    const { corpusId = "" } = useParams()

    const [corpus, setCorpus] = useState<Corpus | undefined>(undefined)
    const [translations, setTranslations] = useState<Translation[] | undefined>(undefined)
    useEffect(() => {
        DefaultService.getCorpus(corpusId)
            .then(({ corpus }) => setCorpus(corpus))
            .catch(console.error)
        DefaultService.listCorpusTranslations(corpusId)
            .then(({ translations }) => setTranslations(translations))
            .catch(console.error)
    }, [])

    return corpus && translations ? <CorpusDetail corpus={corpus} translations={translations} /> : <PageSpin />
}

export default Detail

const SelectTranslationDrawer: FC<
    DrawerProps & {
        translations: Translation[]
        reference?: Translation["id"]
        selected?: Translation["id"][]
        onCancel: () => void
        onConfirm: (ids: Translation["id"][], reference: Translation["id"]) => void
    }
> = (props) => {
    const { translations, selected: initialSelected = [], reference: initialReference, onCancel, onConfirm, ...drawerProps } = props
    const [selected, setSelected] = useState<ImmutableSet<Translation["id"]>>(ImmutableSet(initialSelected))
    const [reference, setReference] = useState<Translation["id"] | undefined>(initialReference)

    return (
        <Drawer
            {...drawerProps}
            extra={
                <Space>
                    <Button onClick={onCancel}>
                        <I18nText id="cancel" transform="capitalize-first" />
                    </Button>
                    <Button type="primary" disabled={selected.size > 0 && reference === undefined} onClick={() => reference && onConfirm(selected.toArray(), reference)}>
                        <I18nText id="confirm" transform="capitalize-first" />
                    </Button>
                </Space>
            }
        >
            <div style={{ display: "flex", flexDirection: "column", height: "100%" }}>
                <List
                    style={{ flexGrow: 1, overflow: "scroll" }}
                    dataSource={translations}
                    renderItem={(t) => (
                        <List.Item key={t.id}>
                            <div>{t.title}</div>
                            <Switch
                                checked={selected.has(t.id)}
                                onChange={(checked) => {
                                    setSelected(checked ? selected.add(t.id) : selected.delete(t.id))
                                    if (checked && !reference) {
                                        setReference(t.id)
                                    }
                                }}
                            />
                        </List.Item>
                    )}
                />
                <Select
                    disabled={selected.size == 0}
                    style={{ width: "100%" }}
                    placeholder={<I18nText id="reference_version" transform="capitalize" />}
                    onChange={setReference}
                    value={reference}
                >
                    {translations
                        .filter((t) => selected.has(t.id))
                        .map((t) => (
                            <Option key={t.id}>{t.title}</Option>
                        ))}
                </Select>
            </div>
        </Drawer>
    )
}

const Pagination: FC<{
    min: number
    max: number
    curr?: number
    onConfirm: (page: number) => void
}> = ({ min, max, curr = 0, onConfirm }) => {
    const [page, setPage] = useState<number>(curr)
    const [popoverOpen, setPopoverOpen] = useState<boolean>(false)

    return (
        <Space>
            <Button
                icon={<LeftOutlined />}
                disabled={page == min}
                onClick={() => {
                    setPage(page - 1)
                    onConfirm(page - 1)
                }}
            />
            <Popover
                trigger="click"
                placement="bottomRight"
                open={popoverOpen}
                onOpenChange={(open) => {
                    setPopoverOpen(open)
                    !open && onConfirm(page)
                }}
                content={
                    <Space>
                        <InputNumber min={min + 1} max={max + 1} step={1} precision={0} onChange={(v) => setPage((v ?? 1) - 1)} value={page + 1} />
                        <Button
                            type="primary"
                            onClick={() => {
                                setPopoverOpen(false)
                                onConfirm(page)
                            }}
                        >
                            <I18nText id="confirm" transform="capitalize" />
                        </Button>
                    </Space>
                }
            >
                <Button style={{ width: 80 }}>{`${page + 1} / ${max + 1}`}</Button>
            </Popover>
            <Button
                icon={<RightOutlined />}
                disabled={page == max}
                onClick={() => {
                    setPage(page + 1)
                    onConfirm(page + 1)
                }}
            />
        </Space>
    )
}
