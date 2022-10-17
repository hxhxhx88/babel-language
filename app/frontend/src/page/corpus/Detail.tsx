import { FC, useState, useEffect, useCallback } from "react"
import { useParams } from "react-router-dom"
import { DefaultService, Corpus, Translation, Block, BlockFilter } from "openapi/babel"
import { Input, Form, Tag, Divider, Switch, Space, Drawer, Breadcrumb, Button, DrawerProps, List, Select, Typography, InputNumber, Popover, Spin } from "antd"
import { SearchOutlined, DownOutlined, RightOutlined, SettingOutlined, LeftOutlined } from "@ant-design/icons"
import { Link } from "react-router-dom"
import { Set as ImmutableSet, Map as ImmutableMap } from "immutable"

import { PAGE_SIZE } from "constant"
import routePath from "route"
import Layout from "Layout"
import PageSpin from "component/PageSpin"
import { I18nText } from "component/Text"

const { Paragraph } = Typography
const { Option } = Select

export interface Props {
    corpus: Corpus
    translations: Translation[]
}

type Search = {
    content?: string
}

function isEmptySearch(s: Search): boolean {
    if (!!s.content) {
        return false
    }
    return true
}

type Query = {
    page: number | undefined
    parents: Block[]
    search: Search | undefined
}

function makeFilter(q: Query): BlockFilter {
    const f: BlockFilter = { ...q.search }
    const n = q.parents.length
    if (n === 0) {
        return f
    }
    return {
        ...f,
        parent_block_id: q.parents[n - 1].id,
    }
}

const CorpusDetail: FC<Props> = (props) => {
    const { corpus, translations } = props

    const [isSelectFormVisible, setIsSelectFormVisible] = useState<boolean | undefined>(true)
    const [selected, setSelected] = useState<Translation["id"][]>([])
    const [reference, setReference] = useState<Translation["id"] | undefined>(undefined)
    const [isCountTranslationBlocks, setIsCountTranslationBlocks] = useState<boolean>(false)
    const [totalCount, setTotalCount] = useState<number>(0)
    const [query, setQuery] = useState<Query>({ parents: [], page: undefined, search: undefined })
    const [isSearchFormVisible, setIsSearchFormVisible] = useState<boolean | undefined>(undefined)
    useEffect(() => {
        if (!reference || query.page !== undefined) return

        const n = query.parents.length
        setIsCountTranslationBlocks(true)
        DefaultService.countTranslationBlocks(reference, {
            filter: makeFilter(query),
        })
            .then(({ total_count }) => {
                setTotalCount(total_count)
                if (total_count > 0) {
                    setQuery({ ...query, page: 0 })
                } else {
                    setReferenceBlocks([])
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
            filter: makeFilter(query),
            pagination: { page: query.page, page_size: PAGE_SIZE },
        })
            .then(({ blocks }) => {
                setReferenceBlocks(blocks)
            })
            .catch(console.error)
            .finally(() => setIsListTranslationBlocks(false))
    }, [reference, query])

    const [parallelBlocks, setParallelBlocks] = useState<ImmutableMap<Block["uuid"], ImmutableMap<Translation["id"], Block>>>(ImmutableMap())
    const [isTranslateBlock, setIsTranslateBlock] = useState<ImmutableSet<Block["id"]>>(ImmutableSet())
    const translateBlock = useCallback(
        (block: Block) => {
            const tids = selected.filter((tid) => tid !== reference)
            if (tids.length === 0) {
                return
            }

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
            <Breadcrumb style={{ cursor: "pointer" }}>
                <Breadcrumb.Item>
                    <Link to={routePath(`/corpuses`)}>
                        <I18nText id="corpus_list" transform="capitalize" />
                    </Link>
                </Breadcrumb.Item>
                <Breadcrumb.Item
                    onClick={() => {
                        setQuery({ ...query, parents: [], page: undefined })
                    }}
                >
                    {corpus.title}
                </Breadcrumb.Item>
                {query.parents.map((p, idx) => (
                    <Breadcrumb.Item
                        key={p.id}
                        onClick={() => {
                            setQuery({ ...query, parents: [...query.parents.slice(0, idx + 1)], page: undefined })
                        }}
                    >
                        {p.content}
                    </Breadcrumb.Item>
                ))}
            </Breadcrumb>

            <div style={{ marginTop: 16, marginBottom: 16, display: "flex", flexDirection: "row", justifyContent: "space-between" }}>
                <Pagination min={0} max={Math.ceil(totalCount / PAGE_SIZE) - 1} onConfirm={(page) => setQuery({ ...query, page })} />
                <Space>
                    <Button icon={<SearchOutlined />} onClick={() => setIsSearchFormVisible(true)} type={query.search === undefined ? "default" : "primary"} />
                    <Button icon={<SettingOutlined />} onClick={() => setIsSelectFormVisible(true)} />
                </Space>
            </div>

            {referenceBlocks.map((block) =>
                block.uuid.endsWith("/") ? (
                    <Button key={block.id} style={{ width: "100%", marginBottom: 8 }} onClick={() => setQuery({ ...query, parents: [...query.parents, block], page: undefined })}>
                        {block.content}
                    </Button>
                ) : (
                    <div key={block.id} style={{ width: "100%", display: "flex", justifyContent: "center", alignItems: "flex-start" }}>
                        <Paragraph style={{ flexGrow: 1, marginRight: 8, marginBottom: 0 }}>
                            <pre style={{ marginTop: 0 }}>
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
                        <Spin spinning={isTranslateBlock.has(block.id)}>
                            <Button
                                style={{ width: 35, height: 35 }}
                                icon={expandedBlocks.has(block.id) ? <DownOutlined /> : <RightOutlined />}
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
                            />
                        </Spin>
                    </div>
                )
            )}

            {isSelectFormVisible !== undefined && (
                <SelectTranslationDrawer
                    width="50%"
                    title={<I18nText id="select_versions" transform="capitalize" />}
                    placement="right"
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

            {isSearchFormVisible !== undefined && (
                <SearchBlockDrawer
                    search={query.search}
                    width="50%"
                    title={<I18nText id="search" transform="capitalize" />}
                    placement="right"
                    closable={false}
                    open={isSearchFormVisible}
                    afterOpenChange={(open) => !open && setIsSearchFormVisible(undefined)}
                    onCancel={() => setIsSearchFormVisible(false)}
                    onConfirm={(search) => {
                        setQuery({ ...query, page: undefined, search })
                        setIsSearchFormVisible(false)
                    }}
                />
            )}

            {(isListTranslationBlocks || isCountTranslationBlocks) && <PageSpin />}
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
                    <Button onClick={onCancel} disabled={initialReference === undefined}>
                        <I18nText id="cancel" transform="capitalize-first" />
                    </Button>
                    <Button type="primary" disabled={selected.size == 0 || reference === undefined} onClick={() => reference && onConfirm(selected.toArray().sort(), reference)}>
                        <I18nText id="confirm" transform="capitalize-first" />
                    </Button>
                </Space>
            }
        >
            <div style={{ display: "flex", flexDirection: "column", height: "100%" }}>
                <Form layout="vertical">
                    <Form.Item label={<I18nText id="reference_version" transform="capitalize" />} style={{ marginBottom: 16 }}>
                        <Select disabled={selected.size == 0} style={{ width: "100%" }} onChange={setReference} value={reference}>
                            {translations
                                .filter((t) => selected.has(t.id))
                                .map((t) => (
                                    <Option key={t.id}>
                                        <Space>
                                            {t.title}
                                            <Tag color="blue">
                                                <I18nText id={`iso_639_3.${t.language_iso_639_3}`} />
                                            </Tag>
                                        </Space>
                                    </Option>
                                ))}
                        </Select>
                    </Form.Item>
                </Form>
                <List
                    style={{ flexGrow: 1, overflow: "scroll" }}
                    dataSource={translations}
                    renderItem={(t) => (
                        <List.Item key={t.id}>
                            <Space>
                                {t.title}
                                <Tag color="blue">
                                    <I18nText id={`iso_639_3.${t.language_iso_639_3}`} />
                                </Tag>
                            </Space>
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
    const disabled = max < min

    return (
        <Space>
            <Button
                icon={<LeftOutlined />}
                disabled={page == min || disabled}
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
                <Button style={{ width: 80 }} disabled={disabled}>
                    {disabled ? `0 / 0` : `${page + 1} / ${max + 1}`}
                </Button>
            </Popover>
            <Button
                icon={<RightOutlined />}
                disabled={page == max || disabled}
                onClick={() => {
                    setPage(page + 1)
                    onConfirm(page + 1)
                }}
            />
        </Space>
    )
}

const SearchBlockDrawer: FC<
    DrawerProps & {
        search?: Search
        onCancel: () => void
        onConfirm: (search: Search | undefined) => void
    }
> = (props) => {
    const { search: initialSearch = {}, onCancel, onConfirm, ...drawerProps } = props
    const [search, setSearch] = useState<Search>(initialSearch)

    return (
        <Drawer
            {...drawerProps}
            extra={
                <Space>
                    <Button onClick={onCancel}>
                        <I18nText id="cancel" transform="capitalize-first" />
                    </Button>
                    <Button onClick={() => onConfirm(undefined)}>
                        <I18nText id="reset" transform="capitalize-first" />
                    </Button>
                    <Button type="primary" onClick={() => onConfirm(isEmptySearch(search) ? undefined : search)}>
                        <I18nText id="confirm" transform="capitalize-first" />
                    </Button>
                </Space>
            }
        >
            <div style={{ display: "flex", flexDirection: "column", height: "100%" }}>
                <Form layout="vertical">
                    <Form.Item label={<I18nText id="content" transform="capitalize" />}>
                        <Input value={search.content} onChange={(e) => setSearch({ ...search, content: e.target.value })} />
                    </Form.Item>
                </Form>
            </div>
        </Drawer>
    )
}
