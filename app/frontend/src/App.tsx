import { BrowserRouter as Router, Routes, Route } from "react-router-dom"
import routePath from "route"

import CorpusList from "page/corpus/List"

function App (): JSX.Element {
    return (
        <Router>
            <Routes>
                <Route path={"/"} element={<CorpusList />} />
                <Route path={process.env.PUBLIC_URL} element={<CorpusList />} />
                <Route path={routePath()} element={<CorpusList />} />
                <Route path={routePath("/corpuses")} element={<CorpusList />} />
            </Routes>
        </Router>
    )
}

export default App
