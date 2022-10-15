const prefix = `${process.env.PUBLIC_URL}/_`

export default function routePath(endpoint = ""): string {
    return `${prefix}${endpoint}`
}

export function trimPrefix(path: string): string {
    if (path.startsWith(prefix)) {
        return path.slice(prefix.length)
    }
    if (path.startsWith(process.env.PUBLIC_URL)) {
        return path.slice(process.env.PUBLIC_URL.length)
    }
    return path
}
