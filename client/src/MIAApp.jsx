import { AppRouter } from "./router/AppRouter"
import { AppTheme } from "./theme/AppTheme"

function MIAApp() {

    return (
        <>
            <AppTheme>
                <AppRouter />
            </AppTheme>
        </>
    )
}

export default MIAApp
