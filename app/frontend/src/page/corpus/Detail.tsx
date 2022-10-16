import { FC, useState, useEffect, useMemo, useCallback } from "react"
import { useParams } from "react-router-dom"
import { DefaultService, Corpus, Translation, Block, BlockFilter } from "openapi/babel"
import { Divider, Switch, Space, Drawer, Alert, Breadcrumb, Button, PageHeader, DrawerProps, List, Select, Typography, InputNumber, Popover, Spin } from "antd"
import { RightOutlined, SettingOutlined, LeftOutlined } from "@ant-design/icons"
import { Link } from "react-router-dom"
import { Set as ImmutableSet, Map as ImmutableMap } from "immutable"

import { PAGE_SIZE } from "constant"
import routePath from "route"
import Layout from "Layout"
import PageSpin from "component/PageSpin"
import { I18nText } from "component/Text"

const { Text, Paragraph } = Typography
const { Option } = Select

export interface Props {
    corpus: Corpus
    translations: Translation[]
}

type Query = {
    page: number | undefined
    parents: Block[]
}

function filterFromQuery(q: Query): BlockFilter {
    const n = q.parents.length
    if (n === 0) {
        return {}
    }
    return {
        parent_block_id: q.parents[n - 1].id,
    }
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
    const [totalCount, setTotalCount] = useState<number>(0)
    const [query, setQuery] = useState<Query>({ parents: [], page: undefined })
    useEffect(() => {
        if (!reference || query.page !== undefined) return

        const n = query.parents.length
        setIsCountTranslationBlocks(true)
        DefaultService.countTranslationBlocks(reference, {
            filter: filterFromQuery(query),
        })
            .then(({ total_count }) => {
                setTotalCount(total_count)
                if (total_count > 0) {
                    setQuery({ ...query, page: 0 })
                }
            })
            .catch(console.error)
            .finally(() => setIsCountTranslationBlocks(false))
    }, [reference, query])

    const [referenceBlocks, setReferenceBlocks] = useState<Block[]>([])
    const [isListTranslationBlocks, setIsListTranslationBlocks] = useState<boolean>(false)
    useEffect(() => {
        if (!reference || query.page === undefined) return

        setIsListTranslationBlocks(true)
        DefaultService.searchTranslationBlocks(reference, {
            filter: filterFromQuery(query),
            pagination: { page: query.page, page_size: PAGE_SIZE },
        })
            .then(({ blocks }) => {
                setReferenceBlocks(blocks)
            })
            .catch(console.error)
            .finally(() => setIsListTranslationBlocks(false))
    }, [query])

    const [parallelBlocks, setParallelBlocks] = useState<ImmutableMap<Block["uuid"], ImmutableMap<Translation["id"], Block>>>(ImmutableMap())
    const [isTranslateBlock, setIsTranslateBlock] = useState<ImmutableSet<Block["id"]>>(ImmutableSet())
    const translateBlock = useCallback(
        (block: Block) => {
            setIsTranslateBlock(isTranslateBlock.add(block.id))
            DefaultService.translateBlock(block.id, {
                translation_ids: selected.filter((tid) => tid !== reference),
            })
                .then(({ blocks }) => {
                    const uuid = block.uuid
                    setParallelBlocks(
                        parallelBlocks.withMutations((p) => {
                            if (!p.has(uuid)) {
                                p.set(uuid, ImmutableMap())
                            }
                            const q = p.get(uuid)?.withMutations((q) => {
                                blocks.forEach((blk) => {
                                    const tid = blk.translation_id
                                    q.set(tid, blk)
                                })
                            })
                            if (q) {
                                p.set(uuid, q)
                            }
                        })
                    )
                })
                .catch(console.error)
                .finally(() => setIsTranslateBlock(isTranslateBlock.remove(block.id)))
        },
        [selected, parallelBlocks]
    )

    const [expandedBlocks, setExpandedBlocks] = useState<ImmutableSet<Block["id"]>>(ImmutableSet())

    return (
        <Layout>
            <Spin spinning={isListTranslationBlocks || isCountTranslationBlocks}>
                <PageHeader
                    ghost={false}
                    breadcrumb={
                        <Breadcrumb>
                            <Breadcrumb.Item>
                                <Link to={routePath(`/corpuses`)}>
                                    <I18nText id="corpus_list" transform="capitalize" />
                                </Link>
                            </Breadcrumb.Item>
                            <Breadcrumb.Item
                                onClick={() => {
                                    setQuery({ parents: [], page: undefined })
                                }}
                            >
                                {corpus.title}
                            </Breadcrumb.Item>
                            {query.parents.map((p, idx) => (
                                <Breadcrumb.Item
                                    key={p.id}
                                    onClick={() => {
                                        setQuery({ parents: [...query.parents.slice(0, idx + 1)], page: undefined })
                                    }}
                                >
                                    {p.content}
                                </Breadcrumb.Item>
                            ))}
                        </Breadcrumb>
                    }
                    extra={[<Button icon={<SettingOutlined />} key="setting" onClick={() => setIsSelectFormVisible(true)} />]}
                    title={corpus.title}
                >
                    {selected.length > 0 ? (
                        reference && (
                            <>
                                <I18nText id="selected_versions" transform="capitalize" />
                                <ul style={{ marginBottom: 0 }}>
                                    <li>
                                        <Text strong={true}>{translationLookup.get(reference)?.title}</Text>
                                    </li>
                                    {selected
                                        .filter((t) => t !== reference)
                                        .map((tid) => (
                                            <li key={tid}>
                                                <Text strong={tid === reference}>{translationLookup.get(tid)?.title}</Text>
                                            </li>
                                        ))}
                                </ul>
                            </>
                        )
                    ) : (
                        <Alert type="warning" showIcon message={<I18nText id="select_version_prompt" transform="capitalize-first" />} />
                    )}
                </PageHeader>
            </Spin>

            {totalCount > 0 && (
                <div style={{ marginTop: 16, marginBottom: 8, display: "flex", flexDirection: "row", justifyContent: "flex-end" }}>
                    <Pagination min={0} max={Math.ceil(totalCount / PAGE_SIZE) - 1} onConfirm={(page) => setQuery({ ...query, page })} />
                </div>
            )}

            {referenceBlocks.map((block) =>
                block.uuid.endsWith("/") ? (
                    <Button key={block.id} style={{ width: "100%", marginTop: 8 }} onClick={() => setQuery({ parents: [...query.parents, block], page: undefined })}>
                        {block.content}
                    </Button>
                ) : (
                    <Spin spinning={isTranslateBlock.has(block.id)} key={block.id}>
                        <Paragraph>
                            <pre
                                onClick={() => {
                                    if (expandedBlocks.has(block.id)) {
                                        setExpandedBlocks(expandedBlocks.remove(block.id))
                                    } else {
                                        setExpandedBlocks(expandedBlocks.add(block.id))
                                        if (!parallelBlocks.has(block.uuid)) {
                                            translateBlock(block)
                                        }
                                    }
                                }}
                            >
                                {block.content}

                                {expandedBlocks.has(block.id) &&
                                    selected
                                        .map((t) => parallelBlocks.get(block.uuid)?.get(t))
                                        .filter(Boolean)
                                        .map((b) => (
                                            <div key={b?.id}>
                                                <Divider style={{ margin: "8px 0" }} />
                                                {b?.content}
                                            </div>
                                        ))}
                            </pre>
                        </Paragraph>
                    </Spin>
                )
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
                    <Button type="primary" disabled={selected.size > 0 && reference === undefined} onClick={() => reference && onConfirm(selected.toArray().sort(), reference)}>
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
