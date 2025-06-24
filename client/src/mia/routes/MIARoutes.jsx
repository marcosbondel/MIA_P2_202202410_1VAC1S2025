import { MIAPage, TerminalPage } from "../pages"
import { Route, Routes } from "react-router-dom";

export const MIARoutes = () => {
    return (
        <Routes>
            <Route path="terminal" element={<TerminalPage />} />
            
            <Route path="/" element={<MIAPage />} />
        </Routes>
    )
}
