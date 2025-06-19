import { AppProvider } from "./context/AppProvider"
import { AppRouter } from "./router/AppRouter"
import { AppTheme } from "./theme/AppTheme"

function MIAApp() {

    return (
        <AppProvider>
            <AppTheme>
                <AppRouter />
            </AppTheme>
        </AppProvider>
    )
}

export default MIAApp
