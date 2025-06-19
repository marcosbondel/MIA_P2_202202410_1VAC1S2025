import { AppContext } from "./AppContext"

export const AppProvider = ({ children }) => {
    return (
        <AppContext.Provider value={{}}>
            {children}
        </AppContext.Provider>
    )
}
