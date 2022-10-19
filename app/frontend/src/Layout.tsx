import { FC } from "react"
import { Typography, Layout } from "antd"
import { Link } from "react-router-dom"
import routePath from "route"

const { Title } = Typography
const { Content, Header } = Layout

export interface Prop {
    children?: React.ReactNode
}

const PageLayout: FC<Prop> = (props) => {
    return (
        <Layout style={{ minHeight: "100vh" }}>
            <Header style={{ display: "flex", flexDirection: "row", alignItems: "center", paddingRight: 16, paddingLeft: 16 }}>
                <Link to={routePath()}>
                    <Title style={{ color: "white", margin: 0, marginTop: 8, fontFamily: "Luminari, fantasy" }}>Babel</Title>
                </Link>
            </Header>
            <Content style={{ padding: 16 }}>{props.children}</Content>
        </Layout>
    )
}

export default PageLayout
