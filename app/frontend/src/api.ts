import { Fetcher } from "openapi-typescript-fetch"
import { paths } from "openapi/babel"

const fetcher = Fetcher.for<paths>()
fetcher.configure({
    baseUrl: "/api",
})

export const createCorpus = fetcher.path("/corpuses").method("post").create()
export const createCorpusTranslation = fetcher.path("/corpus/{corpusId}/translations").method("post").create()
export const listCorpuses = fetcher.path("/corpuses").method("get").create()
export const getMetadata = fetcher.path("/metadata").method("get").create()
